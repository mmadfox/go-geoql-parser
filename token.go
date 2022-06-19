package geoqlparser

import (
	"errors"
	"io"
	"strings"

	"github.com/mmadfox/go-geoql-parser/scanner"
)

type (
	Pos   int
	Token int
)

type Tokenizer struct {
	s   scanner.Scanner
	hop int
	tok rune
	lit string
	err error
}

func NewTokenizer(r io.Reader) *Tokenizer {
	s := &Tokenizer{s: scanner.Scanner{}}
	s.s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings
	s.s.Init(r)
	s.s.Error = func(_ *scanner.Scanner, msg string) {
		s.err = errors.New(msg)
	}
	return s
}

func (t *Tokenizer) next() (rune, string) {
	if t.hop != 0 {
		t.hop = 0
	} else {
		t.tok, t.lit = t.s.Scan(), t.s.TokenText()
	}
	return t.tok, t.lit
}

func (t *Tokenizer) Reset() {
	t.hop = 1
}

func (t *Tokenizer) Unwind() {
	t.hop = 0
}

func (t *Tokenizer) ErrorCount() int {
	return t.s.ErrorCount
}

func (t *Tokenizer) Err() error {
	return t.err
}

func (t *Tokenizer) Offset() Pos {
	return Pos(t.s.Offset)
}

func (t *Tokenizer) TokenText() string {
	return t.s.TokenText()
}

func (t *Tokenizer) Scan() (tok Token, lit string) {
	r, s := t.next()
	lit = strings.ToLower(s)
	switch r {
	case scanner.EOF:
		tok = EOF
	case scanner.Int:
		tok = INT
	case scanner.Float:
		tok = FLOAT
	case scanner.String:
		tok = STRING
	case '.':
		nr, _ := t.next()
		switch nr {
		case '.':
			tok = RANGE
			lit = KeywordString(tok)
		default:
			tok = PERIOD
			t.Reset()
		}
	case '=':
		nr, _ := t.next()
		switch nr {
		case '=':
			tok = LEQL
			lit = KeywordString(tok)
		default:
			tok = ASSIGN
			t.Reset()
		}
	case '@':
		tok = IDENT
	case ';':
		tok = SEMICOLON
	case ',':
		tok = COMMA
	case '(':
		tok = LPAREN
	case ')':
		tok = RPAREN
	case ']':
		tok = RBRACK
	case '[':
		tok = LBRACK
	case '{':
		tok = LBRACE
	case '}':
		tok = RBRACE
	case '/':
		tok = QUO
	case '+':
		tok = ADD
	case '-':
		tok = SUB
	case '*':
		tok = MUL
	case '%':
		tok = REM
	case '!':
		nr, _ := t.next()
		switch nr {
		case '=':
			tok = LNEQ
			lit = KeywordString(tok)
		default:
			tok = ILLEGAL
			t.Reset()
		}
	case '>':
		nr, _ := t.next()
		switch nr {
		case '=':
			tok = GEQ
			lit = KeywordString(tok)
		default:
			tok = GTR
			t.Reset()
		}
	case '<':
		nr, _ := t.next()
		switch nr {
		case '=':
			tok = LEQ
			lit = KeywordString(tok)
		default:
			tok = LSS
			t.Reset()
		}
	case ':':
		tok = COLON
	case scanner.Ident:
		var (
			kwd   Token
			found bool
		)
		switch lit {
		case "point":
			found = true
			kwd = GEOMETRY_POINT
		case "multipoint":
			found = true
			kwd = GEOMETRY_MULTIPOINT
		case "line":
			found = true
			kwd = GEOMETRY_LINE
		case "multiline":
			found = true
			kwd = GEOMETRY_MULTILINE
		case "polygon":
			found = true
			kwd = GEOMETRY_POLYGON
		case "multipolygon":
			found = true
			kwd = GEOMETRY_MULTIPOLYGON
		case "circle":
			found = true
			kwd = GEOMETRY_CIRCLE
		case "collection":
			found = true
			kwd = GEOMETRY_COLLECTION
		case "true", "up", "down", "false":
			found = true
			kwd = BOOLEAN
		case "not":
			found = true
			r, s = t.next()
			lit = strings.ToLower(s)
			notlit := "not " + lit
			notkwd, exists := keywords[notlit]
			if exists {
				found = true
				kwd = notkwd
				lit = notlit
			}
		default:
			kwd, found = keywords[lit]
		}
		if !found {
			tok = SELECTOR
			return
		}
		tok = kwd
	}
	return
}

func (op Token) Precedence() (n int) {
	switch op {
	case OR:
		n = 1
	case AND:
		n = 2
	case LSS, LEQ, GTR, GEQ, EQL, LEQL, INTERSECTS, NOT_INTERSECTS,
		IN, NOT_IN, NEARBY, NOT_NEARBY, NOT_EQ:
		n = 3
	case ADD, SUB:
		n = 4
	case MUL, QUO, REM:
		n = 5
	}
	return
}
