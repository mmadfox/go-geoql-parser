package geoqlparser

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var errNegativeValue = errors.New("value cannot be negative")

func Parse(gql string) (Statement, error) {
	r := strings.NewReader(gql)
	t := NewTokenizer(r)
	s := newParser(t, r)
	return s.parse0()
}

type parser struct {
	r    *strings.Reader
	t    *Tokenizer
	tok  Token
	lit  string
	err  error
	sign Token
	lpos Pos
	rpos Pos
}

func (s *parser) parseTriggerStmt() (stmt *Trigger, err error) {
	if !s.except(WHEN, SET) {
		return nil, s.error()
	}
	stmt = new(Trigger)
	stmt.lpos = s.t.Offset()
	if s.except(SET) {
		if err = s.parseSet(stmt); err != nil {
			return nil, err
		}
		s.next()
		if !s.except(WHEN) {
			return nil, s.error()
		}
	}
	if s.except(WHEN) {
		if err = s.parseWhen(stmt); err != nil {
			return nil, err
		}
	}
	if s.except(REPEAT) {
		if err = s.parseRepeat(stmt); err != nil {
			return nil, err
		}
	}
	if s.except(RESET) {
		if err = s.parseReset(stmt); err != nil {
			return nil, err
		}
	}
	stmt.rpos = s.t.Offset()
	return stmt, nil
}

func (s *parser) parseWhen(stmt *Trigger) error {
	if !s.except(WHEN) {
		return s.error()
	}
	expr, err := s.parseBinaryExpr(1)
	if err != nil {
		return err
	}
	stmt.When = expr
	return nil
}

func (s *parser) parseRepeat(stmt *Trigger) (err error) {
	s.next()
	if s.except(EOF, RESET) {
		return
	}

	if !s.except(INT) {
		return s.error()
	}

	repeatCount, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	if _, ok := repeatCount.(*IntTyp); !ok {
		return s.error()
	}

	stmt.RepeatCount = repeatCount
	if s.except(EOF, RESET) {
		return
	}

	if !s.except(SELECTOR) {
		return s.error()
	}
	s.next()
	dur, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	if _, ok := dur.(*DurationTyp); !ok {
		return s.error()
	}
	stmt.RepeatInterval = dur
	return
}

func (s *parser) parseSet(stmt *Trigger) error {
	s.next()
	for {
		if s.except(WHEN) {
			s.t.Reset()
			break
		}
		if s.except(EOF) {
			break
		}
		if !s.except(SELECTOR) {
			return s.error()
		}
		ident := Ident{Val: s.t.TokenText(), lpos: s.t.Offset()}
		ident.rpos = s.t.Offset() + Pos(len(ident.Val)-1)
		s.t.Unwind()
		s.next()
		if !s.except(ASSIGN) {
			return s.error()
		}
		tokPos := s.t.Offset() - 1
		expr, err := s.parseUnaryExpr()
		if err != nil {
			return err
		}
		if s.except(SEMICOLON) {
			s.next()
		}
		switch typ := expr.(type) {
		case *Ref:
			return s.error()
		case *ArrayTyp:
			if typ.Kind == IDENT {
				return s.error()
			}
		}
		stmt.initVars()
		va := &Assign{
			Left:   &ident,
			TokPos: tokPos,
			Right:  expr,
		}
		if er := stmt.SetVar(va); er != nil {
			s.err = er
			return s.error()
		}
	}
	return nil
}

func (s *parser) parseReset(stmt *Trigger) (err error) {
	s.next()
	s.next()
	var dur Expr
	switch s.tok {
	case INT:
		dur, err = s.parseIntTypes()
	case FLOAT:
		dur, err = s.parseFloatTypes()
	}
	if err != nil {
		return err
	}
	_, ok := dur.(*DurationTyp)
	if !ok {
		return s.error()
	}
	stmt.ResetAfter = dur
	return
}

func (s *parser) next() {
	s.tok, s.lit = s.t.Scan()
}

func (s *parser) parse0() (stmt Statement, err error) {
	s.next()
	switch s.tok {
	case TRIGGER, WHEN:
		if s.except(TRIGGER) {
			s.next()
		}
		return s.parseTriggerStmt()
	default:
		err = s.error()
	}
	return
}

func (s *parser) except(in ...Token) bool {
	for i := 0; i < len(in); i++ {
		if s.tok == in[i] {
			return true
		}
	}
	return false
}

func (s *parser) parseBinaryExpr(oprec0 int) (Expr, error) {
	left, err := s.parseUnaryExpr()
	if err != nil {
		return nil, err
	}
	if s.except(EOF) || s.except(RPAREN) {
		return left, nil
	}
	for {
		if s.except(RPAREN) {
			return left, nil
		}
		op, oprec, pos := s.tok, s.tok.Precedence(), s.t.Offset()
		if oprec < oprec0 {
			return left, nil
		}

		right, err := s.parseBinaryExpr(oprec + 1)
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Right: right, Op: op, OpPos: pos}
	}
}

func (s *parser) parseUnaryExpr() (expr Expr, err error) {
	if s.t.Err() != nil {
		s.err = s.t.Err()
		return nil, s.error()
	}

	s.next()

	if s.except(SUB, ADD) {
		switch s.tok {
		case SUB:
			s.sign = SUB
		case ADD:
			s.sign = ADD
		}
		s.next()
	}

	switch s.tok {
	default:
		s.err = fmt.Errorf("illegal expression")
		err = s.error()
	case DATE:
		expr, err = s.parseDateExpr()
	case TIME:
		expr, err = s.parseTimeExpr()
	case WEEKDAY:
		expr, err = s.parseWeekdayExpr()
	case MONTH:
		expr, err = s.parseMonthExpr()
	case SELECTOR:
		expr, err = s.parseSelectorExpr()
	case MUL:
		expr, err = s.parseWildcardLit()
	case IDENT:
		expr, err = s.parseVarExpr()
	case LBRACK:
		expr, err = s.parseArrayExpr()
	case LPAREN:
		expr, err = s.parseParenExpr()
	case FLOAT:
		expr, err = s.parseFloatTypes()
	case INT:
		expr, err = s.parseIntTypes()
	case STRING:
		expr, err = s.parseStringLit()
	case GEOMETRY_POINT, GEOMETRY_LINE, GEOMETRY_POLYGON:
		expr, err = s.parseGeometryExpr()
	case GEOMETRY_MULTIPOINT, GEOMETRY_MULTILINE, GEOMETRY_MULTIPOLYGON:
		expr, err = s.parseGeometryMultiObject()
	case GEOMETRY_COLLECTION:
		expr, err = s.parseGeometryCollectionExpr()
	case BOOLEAN:
		expr, err = s.parseBooleanLit()
	}
	if err == nil {
		switch s.tok {
		case RANGE:
			expr, err = s.parseRangeExpr(expr)
		}
	}
	s.resetSign()
	return
}

func (s *parser) parseWildcardLit() (expr Expr, err error) {
	return &WildcardTyp{lpos: s.t.Offset()}, nil
}

func (s *parser) resetSign() {
	s.sign = ILLEGAL
}

func (s *parser) isSignMinus() bool {
	return s.sign == SUB
}

func (s *parser) isSignPlus() bool {
	return s.sign == ADD
}

func (s *parser) parseParenExpr() (expr Expr, err error) {
	lp := s.t.Offset()
	expr, err = s.parseBinaryExpr(s.tok.Precedence())
	if err != nil {
		return nil, err
	}
	rp := s.t.Offset()
	s.next()
	return &ParenExpr{Expr: expr, lpos: lp, rpos: rp}, nil
}

func (s *parser) error() error {
	err := Error{
		Offset: s.t.s.Offset,
		Err:    s.err,
		Lit:    s.t.lit,
	}
	_, er := s.r.Seek(0, io.SeekStart)
	if er == nil {
		buf := make([]byte, s.t.s.Offset)
		_, _ = s.r.Read(buf)
		err.Msg = string(buf)
	}
	return &err
}

func newParser(t *Tokenizer, r *strings.Reader) *parser {
	return &parser{t: t, r: r}
}

type Error struct {
	Offset int
	Err    error
	Msg    string
	Lit    string
}

func (e *Error) Error() string {
	var ctx string
	if e.Err != nil {
		ctx = "error: " + e.Err.Error()
	}
	return fmt.Sprintf("syntax error at position offset=%d, near=%s\n```\n%s ...^\n```\n%s",
		e.Offset, e.Lit, strings.TrimSpace(e.Msg), ctx)
}

func isGeometryToken(tok Token) (ok bool) {
	switch tok {
	case GEOMETRY_POINT, GEOMETRY_MULTIPOINT,
		GEOMETRY_LINE, GEOMETRY_MULTILINE,
		GEOMETRY_POLYGON, GEOMETRY_MULTIPOLYGON:
		ok = true
	}
	return
}

func trim(lit string) string {
	lit = strings.TrimLeft(lit, `"`)
	lit = strings.TrimRight(lit, `"`)
	return lit
}
