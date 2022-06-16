package geoqlparser

const (
	ILLEGAL Token = iota
	EOF

	TRIGGER // trigger
	WHEN    // when
	SET     // set
	REPEAT  // repeat
	RESET   // reset

	INT                   // 1
	FLOAT                 // 1.1
	STRING                // "1"
	SPEED                 // 10kmp
	TIME                  // 11:11, 11:11:11
	DATE                  // 2030-10-02
	DATETIME              // 2030-10-02T11:11:11
	DURATION              // 1h, 20s, 7h3m45s, 7h3m, 3m
	TEMPERATURE           // -30C, +30C, -40F
	PRESSURE              // 2.2bar, 2.2psi
	DISTANCE              // 1km, 2m
	PERCENT               // 11%
	IDENT                 // @ident
	BOOLEAN               // true | false
	SELECTOR              // index, count, speed, etc
	GEOMETRY_POINT        // point
	GEOMETRY_MULTIPOINT   // multipoint
	GEOMETRY_LINE         // line
	GEOMETRY_MULTILINE    // multiline
	GEOMETRY_POLYGON      // polygon
	GEOMETRY_MULTIPOLYGON // multipolygon
	GEOMETRY_CIRCLE       // multipolygon
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

	QUO           // /
	MUL           // *
	SUB           // -
	ADD           // +
	GEQ           // >=
	LEQ           // <=
	GTR           // >
	LSS           // <
	LAND          // &&
	LOR           // ||
	AND           // and
	OR            // or
	IN            // in
	NOT_IN        // not in
	EQL           // eq
	LEQL          // ==
	BETWEEN       // between
	NEARBY        // nearby
	INTERSECTS    // intersects
	WITHIN        // within
	INCREASE      // increase
	DECREASE      // decrease
	NEQ           // not eq
	NOTBETWEEN    // not between
	NOTNEARBY     // not nearby
	LNEQ          // !=
	NOTINTERSECTS // not intersects
	NOTWITHIN     // not within
)

var keywords = map[string]Token{
	"trigger":        TRIGGER,
	"set":            SET,
	"when":           WHEN,
	"repeat":         REPEAT,
	"reset":          RESET,
	"between":        BETWEEN,
	"not between":    NOTBETWEEN,
	"=":              ASSIGN,
	";":              SEMICOLON,
	"(":              LPAREN,
	")":              RPAREN,
	",":              COMMA,
	">=":             GEQ,
	"<=":             LEQ,
	">":              GTR,
	"<":              LSS,
	"&&":             LAND,
	"||":             LOR,
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
	"not eq":         NEQ,
	"!=":             LNEQ,
	"in":             IN,
	"@":              IDENT,
	"not in":         NOT_IN,
	"nearby":         NEARBY,
	"not nearby":     NOTNEARBY,
	"intersects":     INTERSECTS,
	"not intersects": NOTINTERSECTS,
	"within":         WITHIN,
	"not within":     NOTWITHIN,
	"increase":       INCREASE,
	"decrease":       DECREASE,
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
