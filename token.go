package geoqlparser

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type (
	Pos   int
	Token int
)

type Tokenizer struct {
	scanner scanner.Scanner
	hop     int
	tok     rune
	lit     string
}

func NewTokenizer(r io.Reader) *Tokenizer {
	s := &Tokenizer{scanner: scanner.Scanner{}}
	s.scanner.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings
	s.scanner.Init(r)
	return s
}

func (t *Tokenizer) errorPos() string {
	pos := t.scanner.Pos()
	return fmt.Sprintf("line:%d,column:%d,offset:%d", pos.Line, pos.Column, pos.Offset)
}

func (t *Tokenizer) next() (rune, string) {
	if t.hop != 0 {
		t.hop = 0
	} else {
		t.tok, t.lit = t.scanner.Scan(), t.scanner.TokenText()
	}
	return t.tok, t.lit
}

func (t *Tokenizer) Reset() {
	t.hop = 1
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
	case '=':
		tok = ASSIGN
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
	case '&':
		nr, _ := t.next()
		switch nr {
		case '&':
			tok = LAND
			lit = KeywordString(tok)
		default:
			tok = ILLEGAL
			t.Reset()
		}
	case '|':
		nr, _ := t.next()
		switch nr {
		case '|':
			tok = LOR
			lit = KeywordString(tok)
		default:
			tok = ILLEGAL
			t.Reset()
		}
	case '!':
		nr, _ := t.next()
		switch nr {
		case '=':
			tok = NEQ
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
	case scanner.Ident:
		keyword, found := keywords[lit]
		if !found {
			tok = UNUSED
			return
		}
		tok = keyword
	}
	return
}
