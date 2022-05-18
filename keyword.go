package geoqlparser

const (
	ILLEGAL Token = iota
	UNUSED
	EOF

	TRIGGER  // trigger
	WHEN     // when
	VARS     // vars
	REPEAT   // repeat
	RESET    // reset
	AFTER    // after
	INTERVAL // interval
	EVERY    // every
	TIMES    // times

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

	GEQ // >=
	LEQ // <=
	NEQ // !=
	GTR // >
	LSS // <

	LAND // &&
	AND  // AND
	LOR  // ||
	OR   // OR
	NOT  // NOT

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
	"every":    EVERY,
	"speed":    SPEED,
	"=":        ASSIGN,
	";":        SEMICOLON,
	"(":        LPAREN,
	")":        RPAREN,
	",":        COMMA,
	">=":       GEQ,
	"<=":       LEQ,
	"!=":       NEQ,
	">":        GTR,
	"<":        LSS,
	"&&":       LAND,
	"||":       LOR,
	"or":       OR,
	"and":      AND,
	"not":      NOT,
	"[":        LBRACK,
	"]":        RBRACK,
}

var keywordStrings = map[Token]string{}

func init() {
	for str, id := range keywords {
		if id == UNUSED {
			continue
		}
		keywordStrings[id] = str
	}
}

func KeywordString(id Token) string {
	str, ok := keywordStrings[id]
	if !ok {
		switch id {
		case UNUSED:
			return "UNUSED"
		case FLOAT:
			return "FLOATVAL"
		case INT:
			return "INTVAL"
		case STRING:
			return "STRINGVAL"
		}
		return ""
	}
	return str
}
