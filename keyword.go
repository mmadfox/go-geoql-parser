package geoqlparser

import "strings"

const (
	ILLEGAL Token = iota
	UNUSED
	EOF

	keywordsBegin
	TRIGGER    // trigger
	WHEN       // when
	VARS       // vars
	REPEAT     // repeat
	RESET      // reset
	AFTER      // after
	INTERVAL   // interval
	TIMES      // times
	AND        // and
	OR         // or
	NOTBETWEEN // not between
	BETWEEN    // between
	keywordsEnd

	INT    // 1
	FLOAT  // 1.1
	STRING // "1"

	ASSIGN    // =
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	COMMA     // ,
	RBRACK    // ]
	LBRACK    // [
	RBRACE    // }
	LBRACE    // {
	QUO       // /
	MUL       // *

	GEQ  // >=
	LEQ  // <=
	NEQ  // !=
	GTR  // >
	LSS  // <
	LAND // &&
	LOR  // ||

	selectorBegin
	TRACKER // tracker
	OBJECT  // object
	SPEED   // speed
	selectorEnd
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
	"between":  BETWEEN,

	"not between": NOTBETWEEN,

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
	"[":   LBRACK,
	"]":   RBRACK,
	"{":   LBRACE,
	"}":   RBRACE,
	"/":   QUO,
	"*":   MUL,
	":":   COLON,

	"tracker": TRACKER,
	"object":  OBJECT,
	"speed":   SPEED,
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

func isSelector(tok Token) bool {
	return tok >= selectorBegin && tok <= selectorEnd
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
