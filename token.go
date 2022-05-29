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
	return fmt.Sprintf("line=%d,column=%d,offset=%d",
		pos.Line, pos.Column, pos.Offset)
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

func (t *Tokenizer) ErrorCount() int {
	return t.scanner.ErrorCount
}

func (t *Tokenizer) Offset() Pos {
	return Pos(t.scanner.Offset)
}

func (t *Tokenizer) hasNextToken(tok Token) (ok bool) {
	r := t.scanner.Peek()
	switch r {
	case '{':
		ok = KeywordString(tok) == "{"
	case '}':
		ok = KeywordString(tok) == "}"
	case ':':
		ok = KeywordString(tok) == ":"
	}
	return
}

func (t *Tokenizer) scanNMEA() (tok Token, lit string, found bool) {
	scanSemi := func() bool {
		r, _ := t.next()
		return r == ':'
	}
	if !scanSemi() {
		return
	}
	ptyp, prefix := t.next()
	if ptyp != scanner.Ident {
		return
	}
	if !scanSemi() {
		return
	}
	styp, ident := t.next()
	if styp != scanner.Ident {
		return
	}
	k := "nmea:" + prefix + ":" + ident
	tt, ok := keywords[k]
	if ok {
		tok = tt
		lit = k
		found = true
	}
	return
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
	case '/':
		tok = QUO
	case '*':
		tok = MUL
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
	case ':':
		tok = COLON
	case scanner.Ident:
		var (
			kwd   Token
			found bool
		)
		switch lit {
		case "nmea":
			kwd, lit, found = t.scanNMEA()
		case "not":
			found = true
			r, s = t.next()
			lit = strings.ToLower(s)
			switch lit {
			default:
				found = false
			case "between":
				kwd = NOTBETWEEN
				lit = "not between"
			}
		default:
			kwd, found = keywords[lit]
		}
		if !found {
			tok = UNUSED
			return
		}
		tok = kwd
	}
	return
}

func (op Token) Precedence() (n int) {
	switch op {
	case LOR, OR:
		n = 1
	case LAND, AND:
		n = 2
	case NEQ, LSS, LEQ, GTR, GEQ:
		n = 3
		//default:
		//	if isSelector(Op) {
		//		n = 3
		//	}
	}
	return
}
