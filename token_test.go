package geoqlparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenizer_Scan(t *testing.T) {
	testCases := []struct {
		name string
		want []Token
		str  string
	}{
		{
			name: "UNUSED",
			want: []Token{UNUSED, UNUSED, UNUSED},
			str:  "a b c",
		},
		{
			name: "TRIGGER,WHEN,VARS,REPEAT,RESET,AFTER",
			want: []Token{TRIGGER, WHEN, VARS, REPEAT, RESET, AFTER},
			str:  "TRIGGER WHEN VARS REPEAT RESET AFTER",
		},
		{
			name: "INT,FLOAT,STRING",
			want: []Token{INT, FLOAT, STRING},
			str:  "1 1.1 \"ok\"",
		},
		{
			name: "ASSIGN,SEMICOLON,LPAREN,RPAREN,COMMA,LBRACK,RBRACK",
			want: []Token{ASSIGN, SEMICOLON, LPAREN, RPAREN, COMMA, LBRACK, RBRACK},
			str:  "= ; ( ) , [ ]",
		},
		{
			name: "GEQ,LEQ,NEQ,GTR,LSS",
			want: []Token{GEQ, LEQ, NEQ, GTR, LSS},
			str:  ">= <= != > <",
		},
		{
			name: "LAND,AND,LOR,OR,NOT",
			want: []Token{LAND, AND, LOR, OR, NOT},
			str:  "&& and || or not",
		},
		{
			name: "RESET,AFTER,INT,UNUSED",
			want: []Token{RESET, AFTER, INT, UNUSED},
			str:  "reset after 24H",
		},
		{
			name: "REPEAT,EVERY,INT,UNUSED",
			want: []Token{REPEAT, EVERY, INT, UNUSED},
			str:  "repeat every 24H",
		},
		{
			name: "REPEAT,INT,TIMES,INTERVAL,INT,UNUSED",
			want: []Token{REPEAT, INT, TIMES, INTERVAL, INT, UNUSED},
			str:  "repeat 25 times interval 10s",
		},
	}
	for _, tc := range testCases {
		tokenizer := NewTokenizer(strings.NewReader(tc.str))
		for i := 0; i < len(tc.want); i++ {
			name := fmt.Sprintf("%s_%s", tc.name, KeywordString(tc.want[i]))
			t.Run(name, func(t *testing.T) {
				tok, _ := tokenizer.Scan()
				if have, want := tok, tc.want[i]; have != want {
					t.Fatalf("have %d, want %d", have, want)
				}
			})
		}
	}
}
