package geoqlparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenizerScan(t *testing.T) {
	testCases := []struct {
		name string
		want []Token
		str  string
	}{
		{
			name: "CALENDAR",
			want: []Token{DATE, TIME, WEEKDAY, MONTH},
			str:  "date time weekday month",
		},
		{
			name: "GEOMETRY",
			want: []Token{GEOMETRY_POINT, GEOMETRY_MULTIPOINT,
				GEOMETRY_LINE, GEOMETRY_MULTILINE,
				GEOMETRY_POLYGON, GEOMETRY_MULTIPOLYGON, GEOMETRY_COLLECTION},
			str: "point multipoint line multiline polygon multipolygon collection",
		},
		{
			name: "SELECTOR RANGE",
			want: []Token{SELECTOR, SELECTOR, SELECTOR, RANGE},
			str:  "someField speed index ..",
		},
		{
			name: "BOOLEAN",
			want: []Token{BOOLEAN, BOOLEAN, BOOLEAN, BOOLEAN},
			str:  "true up false down",
		},
		{
			name: "NEARBY, NOTNEARBY, INTERSECTS, NOTINTERSECTS, WITHIN, NOTWITHIN",
			want: []Token{NEARBY, NOT_NEARBY, INTERSECTS, NOT_INTERSECTS},
			str:  "NEARBY NOT NEARBY INTERSECTS  NOT INTERSECTS",
		},
		{
			name: "EQL, LEQL, NEQ, IN, NOT_IN SUB ADD",
			want: []Token{EQL, LEQL, NOT_EQ, LNEQ, IN, NOT_IN, SUB, ADD},
			str:  "eq == not eq != in not in - +",
		},
		{
			name: "TRIGGER,WHEN,SET,REPEAT,RESET,AFTER",
			want: []Token{TRIGGER, WHEN, SET, REPEAT, RESET},
			str:  "TRIGGER WHEN SET REPEAT RESET ",
		},
		{
			name: "INT,FLOAT,STRING",
			want: []Token{INT, FLOAT, STRING},
			str:  "1 1.1 \"ok\"",
		},
		{
			name: "ASSIGN,SEMICOLON,LPAREN,RPAREN,COMMA,LBRACK,RBRACK,QUO",
			want: []Token{ASSIGN, SEMICOLON, LPAREN, RPAREN, COMMA, LBRACK, RBRACK, QUO},
			str:  "= ; ( ) , [ ] /",
		},
		{
			name: "GEQ,LEQ,LNEQ,GTR,LSS",
			want: []Token{GEQ, LEQ, LNEQ, GTR, LSS},
			str:  ">= <= != > <",
		},
		{
			name: "LAND,AND,LOR,OR",
			want: []Token{AND, OR},
			str:  "and or",
		},
		{
			name: "REPEAT,EVERY,INT,UNUSED",
			want: []Token{REPEAT, INT},
			str:  "repeat  24H",
		},
		{
			name: "ILLEGAL",
			want: []Token{ILLEGAL, ILLEGAL, ILLEGAL},
			str:  "!! |> &",
		},
		{
			name: "ILLEGAL",
			want: []Token{ILLEGAL},
			str:  "&%",
		},
	}
	for _, tc := range testCases {
		tokenizer := NewTokenizer(strings.NewReader(tc.str))
		for i := 0; i < len(tc.want); i++ {
			name := fmt.Sprintf("%s_%s", tc.name, KeywordString(tc.want[i]))
			t.Run(name, func(t *testing.T) {
				tok, _ := tokenizer.Scan()
				if have, want := tok, tc.want[i]; have != want {
					t.Fatalf("have %s, want %s", KeywordString(have), KeywordString(want))
				}
			})
		}
	}
}
