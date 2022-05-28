package geoqlparser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Parse(gql string) (Statement, error) {
	t := NewTokenizer(strings.NewReader(gql))
	s := newParser(t)
	return s.parse0()
}

type parser struct {
	t   *Tokenizer
	pos Pos
	tok Token
	lit string
	err error
}

func (s *parser) next() {
	s.tok, s.lit = s.t.Scan()
	s.pos = s.t.Offset()
}

func (s *parser) parse0() (stmt Statement, err error) {
	s.next()
	switch s.tok {
	case TRIGGER:
		return s.parseTriggerStmt()
	default:
		err = s.syntaxError(nil)
	}
	return
}

func (s *parser) parseTriggerStmt() (stmt *Trigger, err error) {
	s.next()
	if !s.except(WHEN, VARS) {
		return nil, s.syntaxError(nil)
	}
	stmt = new(Trigger)
	if s.except(VARS) {
		if err = s.parseTriggerStmtVars(stmt); err != nil {
			return nil, err
		}
		s.next()
		if !s.except(WHEN) {
			return nil, s.syntaxError(nil)
		}
	}
	if s.except(WHEN) {
		if err = s.parseTriggerStmtWhen(stmt); err != nil {
			return nil, s.syntaxError(nil)
		}
	}
	for i := 0; i < 3; i++ {
		s.next()
		if s.except(EOF) {
			break
		}
		switch s.tok {
		case REPEAT:
			if err = s.parseTriggerStmtRepeat(stmt); err != nil {
				return nil, err
			}
		case RESET:
			if err = s.parseTriggerStmtReset(stmt); err != nil {
				return nil, err
			}
		}
	}
	stmt.Pos = s.t.Offset()
	return stmt, nil
}

func (s *parser) except(in ...Token) bool {
	for i := 0; i < len(in); i++ {
		if s.tok == in[i] {
			return true
		}
	}
	return false
}

func (s *parser) parseTriggerStmtReset(stmt *Trigger) error {
	s.next()
	if !s.except(AFTER) {
		return s.syntaxError(nil)
	}
	buf := strings.Builder{}
	for {
		s.next()
		if s.except(EOF, SEMICOLON) {
			break
		}
		buf.WriteString(s.lit)
	}
	dur, err := time.ParseDuration(buf.String())
	if err != nil {
		return s.syntaxError(err)
	}
	if dur.Seconds() > 0 {
		stmt.Reset = &ResetExpr{Dur: DurVal{V: dur}, Pos: s.t.Offset()}
	}
	return nil
}

func (s *parser) parseTriggerStmtRepeat(stmt *Trigger) error {
	s.next()
	if s.except(EOF, SEMICOLON) {
		return nil
	}
	switch s.tok {
	default:
		return s.syntaxError(nil)
	case INT:
	case UNUSED:
		if s.lit != "once" {
			return s.syntaxError(nil)
		}
		stmt.Repeat = &RepeatExpr{V: 1, Pos: s.t.Offset()}
		return nil
	}
	nv, err := toIntVal(s.lit)
	if err != nil {
		return err
	}
	var short bool
	s.next()
	switch s.tok {
	default:
		return s.syntaxError(nil)
	case TIMES:
	case QUO:
		short = true
	}
	stmt.Repeat = &RepeatExpr{V: nv.V}
	if short {
		dur, err := s.parseDurVal()
		if err != nil {
			return err
		}
		stmt.Repeat.Interval = dur.V
		stmt.Repeat.Pos = s.t.Offset()
		return nil
	}
	s.next()
	if s.except(EOF, SEMICOLON) {
		return nil
	}
	if !s.except(INTERVAL) {
		return s.syntaxError(nil)
	}
	dur, err := s.parseDurVal()
	if err != nil {
		return err
	}
	stmt.Repeat.Interval = dur.V
	return nil
}

func (s *parser) parseDurVal() (v DurVal, err error) {
	buf := strings.Builder{}
	for {
		s.next()
		if s.except(EOF, SEMICOLON) {
			break
		}
		buf.WriteString(s.lit)
	}
	dur, err := time.ParseDuration(buf.String())
	if err != nil {
		return v, s.syntaxError(err)
	}
	v.V = dur
	return
}

func (s *parser) parseTriggerStmtVars(stmt *Trigger) error {
	for {
		vnTok, vnLit := s.t.Scan()
		if vnTok == EOF {
			break
		}
		if vnTok == WHEN {
			s.t.Reset()
			break
		}
		if vnTok != UNUSED {
			return s.errorFromTok(vnTok)
		}
		assignTok, _ := s.t.Scan()
		if assignTok != ASSIGN {
			return s.errorFromTok(assignTok)
		}
		valTok, valLit := s.t.Scan()
		switch valTok {
		case INT:
			n, err := strconv.Atoi(valLit)
			if err != nil {
				return s.errorFromTok(valTok)
			}
			stmt.initVars()
			stmt.Vars[vnLit] = IntVal{V: n}
		case STRING:
			stmt.initVars()
			stmt.Vars[vnLit] = StrVal{V: valLit}
		case FLOAT:
			n, err := strconv.ParseFloat(valLit, 64)
			if err != nil {
				return s.errorFromTok(valTok)
			}
			stmt.initVars()
			stmt.Vars[vnLit] = FloatVal{V: n}
		case LBRACE:
			s.t.Reset()
			list, err := s.parseList()
			if err != nil {
				return err
			}
			if list != nil {
				stmt.initVars()
				stmt.Vars[vnLit] = list
			}
		case LBRACK:
			s.t.Reset()
			array, err := s.parseArray()
			if err != nil {
				return err
			}
			if array != nil {
				stmt.initVars()
				stmt.Vars[vnLit] = array
			}
		default:
			return s.errorFromTok(valTok)
		}
	}
	return nil
}

func (s *parser) parseList() (interface{}, error) {
	var (
		index    int
		typ      Token
		intVal   map[int]struct{}
		strVal   map[string]struct{}
		floatVal map[float64]struct{}
	)

	for {
		s.next()
		if s.except(LBRACE, COMMA) {
			continue
		}
		if s.except(RBRACE, EOF) {
			break
		}
		switch s.tok {
		default:
			return nil, s.syntaxError(nil)
		case INT:
			if index == 0 {
				typ = INT
			}
			if index > 0 && typ != INT {
				return nil, s.syntaxError(nil)
			}
			val, err := toIntVal(s.lit)
			if err != nil {
				return nil, s.syntaxError(nil)
			}
			if intVal == nil {
				intVal = make(map[int]struct{})
			}
			intVal[val.V] = struct{}{}
		case STRING:
			if index == 0 {
				typ = STRING
			}
			if index > 0 && typ != STRING {
				return nil, s.syntaxError(nil)
			}
			val, err := toStringVal(s.lit)
			if err != nil {
				return nil, s.syntaxError(nil)
			}
			if strVal == nil {
				strVal = make(map[string]struct{})
			}
			strVal[val.V] = struct{}{}
		case FLOAT:
			if index == 0 {
				typ = FLOAT
			}
			if index > 0 && typ != FLOAT {
				return nil, s.syntaxError(nil)
			}
			val, err := toFloatVal(s.lit)
			if err != nil {
				return nil, s.syntaxError(nil)
			}
			if floatVal == nil {
				floatVal = make(map[float64]struct{})
			}
			floatVal[val.V] = struct{}{}
		}
		index++
	}
	switch typ {
	case INT:
		if intVal == nil {
			return nil, nil
		}
		return ListIntVal{V: intVal}, nil
	case FLOAT:
		if floatVal == nil {
			return nil, nil
		}
		return ListFloatVal{V: floatVal}, nil
	case STRING:
		if strVal == nil {
			return nil, nil
		}
		return ListStringVal{V: strVal}, nil
	default:
		return nil, nil
	}
}

func (s *parser) parseArray() (interface{}, error) {
	var (
		index    int
		typ      Token
		intVal   []int
		strVal   []string
		floatVal []float64
	)

	for {
		tok, lit := s.t.Scan()
		if tok == LBRACK || tok == COMMA {
			continue
		}
		if tok == RBRACK || tok == EOF {
			break
		}
		switch tok {
		default:
			return nil, fmt.Errorf("syntax error at position %s near '%s'",
				s.t.errorPos(), lit)
		case INT:
			if index == 0 {
				typ = INT
			}
			if index > 0 && typ != INT {
				return nil, fmt.Errorf("syntax error at position %s near '%s', got INTVAL, expected %s",
					s.t.errorPos(), lit, type2str(typ))
			}
			val, err := toIntVal(lit)
			if err != nil {
				return nil, s.errorFromTok(tok)
			}
			if intVal == nil {
				intVal = make([]int, 0)
			}
			intVal = append(intVal, val.V)
		case STRING:
			if index == 0 {
				typ = STRING
			}
			if index > 0 && typ != STRING {
				return nil, fmt.Errorf("syntax error at position %s near '%s', got STRINGVAL, expected %s",
					s.t.errorPos(), lit, type2str(typ))
			}
			val, err := toStringVal(lit)
			if err != nil {
				return nil, s.errorFromTok(tok)
			}
			if strVal == nil {
				strVal = make([]string, 0)
			}
			strVal = append(strVal, val.V)
		case FLOAT:
			if index == 0 {
				typ = FLOAT
			}
			if index > 0 && typ != FLOAT {
				return nil, fmt.Errorf("syntax error at position %s near '%s', got FLOATVAL, expected %s",
					s.t.errorPos(), lit, type2str(typ))
			}
			val, err := toFloatVal(lit)
			if err != nil {
				return nil, s.errorFromTok(tok)
			}
			if floatVal == nil {
				floatVal = make([]float64, 0)
			}
			floatVal = append(floatVal, val.V)
		}
		index++
	}
	switch typ {
	case INT:
		if intVal == nil {
			return nil, nil
		}
		return ArrayIntVal{V: intVal}, nil
	case FLOAT:
		if floatVal == nil {
			return nil, nil
		}
		return ArrayFloatVal{V: floatVal}, nil
	case STRING:
		if strVal == nil {
			return nil, nil
		}
		return ArrayStringVal{V: strVal}, nil
	default:
		return nil, nil
	}
}

func (s *parser) parseTriggerStmtWhen(stmt *Trigger) error {
	if !s.except(WHEN) {
		return s.syntaxError(nil)
	}
	expr, err := s.parseBinaryExpr(1)
	if err != nil {
		return err
	}
	stmt.When = expr
	return nil
}

func (s *parser) parseBinaryExpr(oprec0 int) (Expr, error) {
	s.next()
	left, err := s.parseUnaryExpr()
	if err != nil {
		// short form: [trigger when * ...].
		// mostly for tests.
		if err == errSkipExpr && oprec0 == 1 {
			return left, nil
		}
		return nil, err
	}
	s.next()
	if s.except(EOF) || s.except(RPAREN) {
		return left, nil
	}
	if !isOperator(s.tok) {
		return nil, s.syntaxError(nil)
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

		left = &BinaryExpr{Left: left, Right: right, Op: op, Pos: pos}
	}
}

func (s *parser) parseUnaryExpr() (Expr, error) {
	if isSelector(s.tok) {
		s.t.Reset()
		if s.tok == TRACKER {
			return s.parseTrackerSelectorExpr()
		} else {
			return s.parseBaseSelectorExpr()
		}
	}
	switch s.tok {
	case LPAREN:
		lpos := s.t.Offset()
		expr, err := s.parseBinaryExpr(s.tok.Precedence())
		if err != nil {
			return nil, err
		}
		rpos := s.t.Offset()
		return &ParenExpr{Expr: expr, LeftPos: lpos, RightPos: rpos}, nil
	case MUL:
		return &WildcardLit{Pos: s.t.Offset()}, errSkipExpr
	case INT:
		val, err := toIntVal(s.lit)
		if err != nil {
			return nil, s.errorFromLit(s.lit)
		}
		expr := BasicLit{Pos: s.t.Offset(), Kind: INT, V: val}
		return &expr, nil
	case STRING:
		val, err := toStringVal(s.lit)
		if err != nil {
			return nil, s.errorFromLit(s.lit)
		}
		expr := BasicLit{Pos: s.t.Offset(), Kind: STRING, V: val}
		return &expr, nil
	case FLOAT:
		val, err := toFloatVal(s.lit)
		if err != nil {
			return nil, s.errorFromLit(s.lit)
		}
		expr := BasicLit{Pos: s.t.Offset(), Kind: FLOAT, V: val}
		return &expr, nil
	}
	return nil, s.errorFromTok(s.tok)
}

func (s *parser) parseTrackerSelectorExpr() (Expr, error) {
	s.next()
	if !isSelector(s.tok) {
		return nil, s.syntaxError(nil)
	}
	expr := &TrackerSelectorExpr{Ident: s.tok, Radius: DefaultRadiusVal}
	// short form: tracker
	if !s.t.hasNextToken(LBRACE) && !s.t.hasNextToken(COLON) {
		expr.Wildcard = true
		expr.LeftPos = s.t.Offset()
		return expr, nil
	}

	parseRadius := func(expr *TrackerSelectorExpr) error {
		for i := 0; i < 2; i++ {
			s.next()
			if i == 0 && s.tok == COLON {
				continue
			}
			if !s.except(UNUSED) {
				return s.syntaxError(nil)
			}
			radius, err := toRadiusVal(s.lit)
			if err != nil {
				return err
			}
			expr.Radius = radius
		}
		return nil
	}

	// short form with radius: tracker:1km
	if s.t.hasNextToken(COLON) {
		expr.Wildcard = true
		// with radius
		if err := parseRadius(expr); err != nil {
			return nil, err
		}
		expr.LeftPos = s.t.Offset()
		return expr, nil
	}

	if s.t.hasNextToken(LBRACE) {
		// with args: {@var, *, "uuid"}
		var li int
		for {
			s.next()
			if s.except(LBRACE) {
				if li == 0 {
					li++
					continue
				} else {
					return nil, s.errorFromLit(s.lit)
				}
			}
			if s.except(RBRACE, EOF) {
				break
			}
			if s.except(COMMA) {
				continue
			}
			// - with vars
			if s.tok == ILLEGAL && s.lit == "@" {
				if expr.Vars == nil {
					expr.Vars = make(map[string]struct{})
				}
				s.next()
				if !s.except(UNUSED) {
					return nil, s.errorFromLit(s.lit)
				}
				expr.Vars[s.lit] = struct{}{}
			}
			// - with identifier
			if s.except(STRING) {
				if expr.Args == nil {
					expr.Args = make(map[string]struct{})
				}
				expr.Args[trim(s.lit)] = struct{}{}
			}
			// - with wildcard
			if s.except(MUL) {
				expr.Wildcard = true
			}
		}
	}

	// with radius
	if s.t.hasNextToken(COLON) {
		expr.Wildcard = true
		// with radius
		if err := parseRadius(expr); err != nil {
			return nil, err
		}
		return expr, nil
	}
	expr.LeftPos = s.t.Offset()
	return expr, nil
}

func (s *parser) parseBaseSelectorExpr() (Expr, error) {
	tok, _ := s.t.Scan()
	if !isSelector(tok) {
		return nil, s.errorFromTok(tok)
	}
	expr := BaseSelectorExpr{Ident: tok, Qualifier: Any}
	tok, _ = s.t.Scan()
	// short form: speed,object,etc...
	if tok != LBRACE {
		s.t.Reset()
		expr.Wildcard = true
		return &expr, nil
	}
	// with args: {@var, *, "uuid"}
	for {
		tok, lit := s.t.Scan()
		if tok == RBRACE || tok == EOF {
			break
		}
		if tok == COMMA {
			continue
		}
		// - with vars
		if tok == ILLEGAL && lit == "@" {
			if expr.Vars == nil {
				expr.Vars = make(map[string]struct{})
			}
			tok, lit = s.t.Scan()
			if tok != UNUSED {
				return nil, s.errorFromLit(lit)
			}
			expr.Vars[lit] = struct{}{}
		}
		// - with identifier
		if tok == STRING {
			if expr.Args == nil {
				expr.Args = make(map[string]struct{})
			}
			expr.Args[trim(lit)] = struct{}{}
		}
		// - with wildcard
		if tok == MUL {
			expr.Wildcard = true
		}
	}
	// qualifier
	tok, _ = s.t.Scan()
	if tok != COLON {
		s.t.Reset()
	} else {
		_, lit := s.t.Scan()
		switch strings.ToLower(lit) {
		default:
			return nil, s.errorFromTok(tok)
		case "any":
			expr.Qualifier = Any
		case "all":
			expr.Qualifier = All
		}
	}
	return &expr, nil
}

func (s *parser) errorFromLit(lit string) error {
	return newError(s.t, "near"+" '"+lit+"'")
}

func (s *parser) syntaxError(withCtx error) error {
	s.err = withCtx
	return s.errorFromTok(s.tok)
}

func (s *parser) errorFromTok(tok Token) error {
	var lit string
	if tok == UNUSED {
		lit = s.t.lit
	} else {
		lit = KeywordString(tok)
	}
	msg := "near" + " '" + lit + "'"
	if s.err != nil {
		msg += " with error:" + s.err.Error()
	}
	return newError(s.t, msg)
}

func newParser(t *Tokenizer) *parser {
	return &parser{t: t}
}
