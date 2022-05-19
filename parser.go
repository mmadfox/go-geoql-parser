package geoqlparser

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(gql string) (Statement, error) {
	t := NewTokenizer(strings.NewReader(gql))
	s := newParser(t)
	return s.parse0()
}

type parser struct {
	t           *Tokenizer
	triggerStmt *Trigger
}

func (s *parser) parse0() (stmt Statement, err error) {
	tok, _ := s.t.Scan()
	switch tok {
	case TRIGGER:
		return s.parseTriggerStmt()
	default:
		err = s.toError(tok)
	}
	return
}

func (s *parser) parseTriggerStmt() (stmt *Trigger, err error) {
	tok, _ := s.t.Scan()
	if tok != WHEN && tok != VARS {
		return nil, s.toError(tok)
	}

	stmt = new(Trigger)
	if tok == VARS {
		if err = s.parseTriggerStmtVars(stmt); err != nil {
			return nil, err
		}
		tok, _ = s.t.Scan()
		if tok != WHEN {
			return nil, s.toError(tok)
		}
	}

	if tok == WHEN {
		if err = s.parseTriggerStmtWhen(stmt); err != nil {
			return nil, s.toError(tok)
		}
	}

	for i := 0; i < 3; i++ {
		tok, _ = s.t.Scan()
		if tok == EOF {
			break
		}
		switch tok {
		case REPEAT:
		case RESET:
		}
	}
	return stmt, nil
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
			return s.toError(vnTok)
		}
		assignTok, _ := s.t.Scan()
		if assignTok != ASSIGN {
			return s.toError(assignTok)
		}
		valTok, valLit := s.t.Scan()
		switch valTok {
		case INT:
			n, err := strconv.Atoi(valLit)
			if err != nil {
				return s.toError(valTok)
			}
			stmt.initVars()
			stmt.Vars[vnLit] = IntVal{V: n}
		case STRING:
			stmt.initVars()
			stmt.Vars[vnLit] = StrVal{V: valLit}
		case FLOAT:
			n, err := strconv.ParseFloat(valLit, 64)
			if err != nil {
				return s.toError(valTok)
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
			return s.toError(valTok)
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
				return nil, s.toError(tok)
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
				return nil, s.toError(tok)
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
				return nil, s.toError(tok)
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
				return nil, s.toError(tok)
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
				return nil, s.toError(tok)
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
				return nil, s.toError(tok)
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
	return nil
}

func (s *parser) toError(tok Token) error {
	var lit string
	if tok == UNUSED {
		lit = s.t.lit
	} else {
		lit = KeywordString(tok)
	}
	return fmt.Errorf("syntax error at position %s near '%s'",
		s.t.errorPos(), lit)
}

func newParser(t *Tokenizer) *parser {
	return &parser{t: t}
}
