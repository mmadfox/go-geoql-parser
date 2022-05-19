package geoqlparser

import "strings"

const (
	ILLEGAL Token = iota
	UNUSED
	EOF

	keywordsBegin
	TRIGGER  // trigger
	WHEN     // when
	VARS     // vars
	REPEAT   // repeat
	RESET    // reset
	AFTER    // after
	INTERVAL // interval
	TIMES    // times
	AND      // and
	OR       // or
	NOT      // not
	keywordsEnd

	INT    // 1
	FLOAT  // 1.1
	STRING // "1"

	ASSIGN    // =
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	COMMA     // ,
	RBRACK    // ]
	LBRACK    // [
	RBRACE    // }
	LBRACE    // {
	QUO       // /

	GEQ  // >=
	LEQ  // <=
	NEQ  // !=
	GTR  // >
	LSS  // <
	LAND // &&
	LOR  // ||

	SPEED // speed selector
)

var keywords = map[string]Token{
	"trigger":  TRIGGER,
	"vars":     VARS,
	"when":     WHEN,
	"repeat":   REPEAT,
	"reset":    RESET,
	"after":    AFTER,
	"interval": INTERVAL,
	"times":    TIMES,

	"=":   ASSIGN,
	";":   SEMICOLON,
	"(":   LPAREN,
	")":   RPAREN,
	",":   COMMA,
	">=":  GEQ,
	"<=":  LEQ,
	"!=":  NEQ,
	">":   GTR,
	"<":   LSS,
	"&&":  LAND,
	"||":  LOR,
	"or":  OR,
	"and": AND,
	"not": NOT,
	"[":   LBRACK,
	"]":   RBRACK,
	"{":   LBRACE,
	"}":   RBRACE,
	"/":   QUO,
}

var keywordStrings = map[Token]string{}

func init() {
	for str, id := range keywords {
		keywordStrings[id] = str
	}
}

func KeywordString(id Token) string {
	str, ok := keywordStrings[id]
	if !ok {
		return type2str(id)
	}
	if id >= keywordsBegin && id <= keywordsEnd {
		str = strings.ToUpper(str)
	}
	return str
}

func type2str(id Token) (str string) {
	switch id {
	case UNUSED:
		str = "UNUSED"
	case FLOAT:
		str = "FLOATVAL"
	case INT:
		str = "INTVAL"
	case STRING:
		str = "STRINGVAL"
	}
	return
}
