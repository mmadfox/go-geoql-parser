package geoqlparser

import (
	"testing"
)

func TestParseTriggerStmtStringListVal(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want map[string]map[string]struct{}
	}{
		{
			str: "trigger vars a={\"70c960f3-4b56-4d71-a04a-2c62a714f4af\", \"one\"} when",
			want: map[string]map[string]struct{}{
				"a": {
					"70c960f3-4b56-4d71-a04a-2c62a714f4af": {},
					"one":                                  {},
				},
			},
		},
		{
			str: "trigger vars a={\"one\", \"two\",,,,} when",
			want: map[string]map[string]struct{}{
				"a": {
					"one": {},
					"two": {},
				},
			},
		},
		{
			str: "trigger vars a={\"one\", 1, \"two\"} when",
			err: true,
		},
		{
			str: "trigger vars a={1, \"two\"} when",
			err: true,
		},
		{
			str: "trigger vars a={1.1, \"one\", 1, \"two\"} when",
			err: true,
		},
	}
	for _, tc := range testCases {
		stmt, err := Parse(tc.str)
		if tc.err {
			if err == nil {
				t.Fatalf("got nil, expected error")
			} else {
				continue
			}
		} else if !tc.err && err != nil {
			t.Fatal(err)
		}
		trigger := stmt.(*Trigger)
		if have, want := len(trigger.Vars), len(tc.want); have != want {
			t.Fatalf("have %d, want %d string vars", have, want)
		}
		for varname, want := range tc.want {
			vals, found := trigger.Vars[varname]
			if !found {
				t.Fatalf("variable not found %s", varname)
			}
			vars := vals.(ListStringVal)
			for a := range want {
				_, ok := vars.V[a]
				if !ok {
					t.Fatalf("list %s, item not found %s", varname, a)
				}
			}
		}
	}
}

//func TestParseTriggerStatement(t *testing.T) {
//	testCases := []struct {
//		str string
//		err bool
//	}{
//		{
//			str: "trigger when",
//		},
//		{
//			str: "trigger vars when",
//		},
//		{
//			str: "trigger repeat",
//			err: true,
//		},
//		{
//			str: "trigger vars group=1 groub=1.0 when repeat reset",
//		},
//		{
//			str: "trigger when reset",
//		},
//		{
//			str: "trigger when reset repeat",
//		},
//	}
//	for _, tc := range testCases {
//		stmt, err := Parse(tc.str)
//		if tc.err {
//			if err == nil {
//				t.Fatalf("got nil, expected error")
//			} else {
//				continue
//			}
//		} else if !tc.err && err != nil {
//			t.Fatal(err)
//		}
//		_ = stmt
//	}
//}
