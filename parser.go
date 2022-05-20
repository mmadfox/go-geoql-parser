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
	t *Tokenizer
}

func (s *parser) parse0() (stmt Statement, err error) {
	tok, _ := s.t.Scan()
	switch tok {
	case TRIGGER:
		return s.parseTriggerStmt()
	default:
		err = s.errorFromTok(tok)
	}
	return
}

func (s *parser) parseTriggerStmt() (stmt *Trigger, err error) {
	tok, _ := s.t.Scan()
	if tok != WHEN && tok != VARS {
		return nil, s.errorFromTok(tok)
	}

	stmt = new(Trigger)
	stmt.Reset = DefaultResetVal
	stmt.Repeat = DefaultRepeatVal
	if tok == VARS {
		if err = s.parseTriggerStmtVars(stmt); err != nil {
			return nil, err
		}
		tok, _ = s.t.Scan()
		if tok != WHEN {
			return nil, s.errorFromTok(tok)
		}
	}

	if tok == WHEN {
		if err = s.parseTriggerStmtWhen(stmt); err != nil {
			return nil, s.errorFromTok(tok)
		}
	}

	for i := 0; i < 3; i++ {
		tok, _ = s.t.Scan()
		if tok == EOF {
			break
		}
		switch tok {
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
	return stmt, nil
}

func (s *parser) parseTriggerStmtReset(stmt *Trigger) error {
	tok, lit := s.t.Scan()
	if tok != AFTER {
		return fmt.Errorf("syntax error at position %s near 'RESET', expected 'RESET AFTER'",
			s.t.errorPos())
	}
	buf := strings.Builder{}
	for {
		tok, lit = s.t.Scan()
		if tok == EOF || tok == SEMICOLON {
			break
		}
		buf.WriteString(lit)
	}
	dur, err := time.ParseDuration(buf.String())
	if err != nil {
		return fmt.Errorf("syntax error at position %s near 'RESET AFTER' %s",
			s.t.errorPos(), err)
	}
	if dur.Seconds() > 0 {
		stmt.Reset = DurVal{V: dur}
	}
	return nil
}

func (s *parser) parseTriggerStmtRepeat(stmt *Trigger) error {
	tok, n := s.t.Scan()
	if tok == EOF || tok == SEMICOLON {
		return nil
	}

	switch tok {
	default:
		return fmt.Errorf("syntax error at position %s near 'REPEAT'",
			s.t.errorPos())
	case INT:
	case UNUSED:
		if n != "once" {
			return fmt.Errorf("syntax error at position %s near 'REPEAT', got 'REPEAT %s', expected 'REPEAT once'",
				s.t.errorPos(), n)
		}
		stmt.Repeat = Repeat{V: 1}
		return nil
	}

	var short bool
	tok, v := s.t.Scan()
	switch tok {
	default:
		return fmt.Errorf("syntax error at position %s near 'REPEAT %s %s'",
			s.t.errorPos(), n, v)
	case TIMES:
	case QUO:
		short = true
	}

	nv, err := toIntVal(n)
	if err != nil {
		return err
	}
	stmt.Repeat.V = nv.V

	if short {
		dur, err := s.parseDurVal()
		if err != nil {
			return err
		}
		stmt.Repeat.Interval = dur.V
		return nil
	}

	tok, _ = s.t.Scan()
	if tok == EOF || tok == SEMICOLON {
		return nil
	}
	if tok != INTERVAL {
		return fmt.Errorf("syntax error at position %s near 'REPEAT %s %s'",
			s.t.errorPos(), n, v)
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
		tok, lit := s.t.Scan()
		if tok == EOF || tok == SEMICOLON {
			break
		}
		buf.WriteString(lit)
	}
	dur, err := time.ParseDuration(buf.String())
	if err != nil {
		return v, fmt.Errorf("syntax error at position %s near '%s'",
			s.t.errorPos(), s.t.lit)
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
		tok, lit := s.t.Scan()
		if tok == LBRACE || tok == COMMA {
			continue
		}
		if tok == RBRACE || tok == EOF {
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
				intVal = make(map[int]struct{})
			}
			intVal[val.V] = struct{}{}
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
				strVal = make(map[string]struct{})
			}
			strVal[val.V] = struct{}{}
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
	expr, err := s.parseWhen()
	if err != nil {
		return err
	}
	stmt.When = expr
	return nil
}

func (s *parser) parseWhen() (Expr, error) {
	expr, err := s.parseExpr()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (s *parser) parseExpr() (Expr, error) {
	tok, _ := s.t.Scan()
	if isSelector(tok) {
		s.t.Reset()
		if tok == TRACKER {
			return s.parseTrackerSelector()
		} else {
			return s.parseBaseSelector()
		}
	}
	if tok == MUL {
		return &WildcardLit{}, nil
	}
	return nil, s.errorFromTok(tok)
}

func (s *parser) parseTrackerSelector() (Expr, error) {
	tok, _ := s.t.Scan()
	if !isSelector(tok) {
		return nil, s.errorFromTok(tok)
	}
	expr := TrackerSelectorLit{Ident: tok, Radius: DefaultRadiusVal}
	tok, _ = s.t.Scan()
	// short form: tracker
	if tok != LBRACE {
		s.t.Reset()
		expr.Wildcard = true
		tok, lit := s.t.Scan()
		if tok != COLON {
			s.t.Reset()
			return &expr, nil
		}
		// with radius
		tok, lit = s.t.Scan()
		if tok != UNUSED {
			return nil, s.errorFromLit(lit)
		}
		radius, err := toRadiusVal(lit)
		if err != nil {
			return nil, err
		}
		expr.Radius = radius
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
	// radius
	tok, _ = s.t.Scan()
	if tok != COLON {
		s.t.Reset()
	} else {
		tok, lit := s.t.Scan()
		if tok != UNUSED {
			return nil, s.errorFromLit(lit)
		}
		radius, err := toRadiusVal(lit)
		if err != nil {
			return nil, err
		}
		expr.Radius = radius
	}
	return &expr, nil
}

func (s *parser) parseBaseSelector() (Expr, error) {
	tok, _ := s.t.Scan()
	if !isSelector(tok) {
		return nil, s.errorFromTok(tok)
	}
	expr := BaseSelectorLit{Ident: tok, Qualifier: Any}
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

func (s *parser) errorFromTok(tok Token) error {
	var lit string
	if tok == UNUSED {
		lit = s.t.lit
	} else {
		lit = KeywordString(tok)
	}
	return newError(s.t, "near"+" '"+lit+"'")
}

func newParser(t *Tokenizer) *parser {
	return &parser{t: t}
}
