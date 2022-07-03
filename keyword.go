package geoqlparser

const (
	ILLEGAL Token = iota
	EOF

	TRIGGER // trigger
	WHEN    // when
	SET     // set
	REPEAT  // repeat
	RESET   // reset

	INT         // 1
	FLOAT       // 1.1
	STRING      // "1"
	SPEED       // 10kmp
	TIME        // 11:11, 11:11:11
	DATE        // 2030-10-02
	WEEKDAY     // Mon
	MONTH       // Jan
	DURATION    // 1h, 20s, 7h3m45s, 7h3m, 3m
	TEMPERATURE // -30C, +30C, -40F
	PRESSURE    // 2.2bar, 2.2psi
	DISTANCE    // 1km, 2m
	PERCENT     // 11%
	IDENT       // @ident
	RANGE       // low .. high
	BOOLEAN     // true | false
	SELECTOR    // index, count, speed, etc

	GEOMETRY_POINT        // point
	GEOMETRY_LINE         // line
	GEOMETRY_POLYGON      // polygon
	GEOMETRY_MULTILINE    // multiline
	GEOMETRY_MULTIPOINT   // multipoint
	GEOMETRY_MULTIPOLYGON // multipolygon
	GEOMETRY_COLLECTION   // collection

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
	PERIOD    // .

	QUO // /
	MUL // *
	SUB // -
	ADD // +
	REM // mod

	GEQ            // >=
	LEQ            // <=
	GTR            // >
	LSS            // <
	AND            // and
	OR             // or
	IN             // in
	NOT_IN         // not in
	EQL            // eq
	LEQL           // ==
	NEARBY         // nearby
	INTERSECTS     // intersects
	NOT_EQ         // not eq
	NOT_NEARBY     // not nearby
	LNEQ           // !=
	NOT_INTERSECTS // not intersects
)

var keywords = map[string]Token{
	"trigger": TRIGGER,
	"set":     SET,
	"when":    WHEN,
	"repeat":  REPEAT,
	"reset":   RESET,

	"=":              ASSIGN,
	";":              SEMICOLON,
	"(":              LPAREN,
	")":              RPAREN,
	",":              COMMA,
	">=":             GEQ,
	"<=":             LEQ,
	">":              GTR,
	"<":              LSS,
	"rem":            REM,
	"or":             OR,
	"and":            AND,
	"[":              LBRACK,
	"]":              RBRACK,
	"{":              LBRACE,
	"}":              RBRACE,
	"/":              QUO,
	"*":              MUL,
	"-":              SUB,
	"+":              ADD,
	":":              COLON,
	"eq":             EQL,
	"==":             LEQL,
	".":              PERIOD,
	"not eq":         NOT_EQ,
	"!=":             LNEQ,
	"in":             IN,
	"range":          RANGE,
	"@":              IDENT,
	"not in":         NOT_IN,
	"nearby":         NEARBY,
	"not nearby":     NOT_NEARBY,
	"intersects":     INTERSECTS,
	"not intersects": NOT_INTERSECTS,

	"time":    TIME,
	"date":    DATE,
	"weekday": WEEKDAY,
	"month":   MONTH,
}

var keywordStrings = map[Token]string{}

func init() {
	for str, id := range keywords {
		keywordStrings[id] = str
	}
}

func KeywordString(id Token) string {
	str, _ := keywordStrings[id]
	return str
}
