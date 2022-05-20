package geoqlparser

import (
	"testing"
	"time"
)

func TestParseTriggerStmtWhen(t *testing.T) {
	testCases := []struct {
		str string
		err bool
	}{
		{
			str: "trigger when * tracker > 50",
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
		_ = trigger
	}
}

func TestParseTriggerStmtRepeat(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want Repeat
	}{
		{
			str:  "trigger when * repeat once; reset after 1s;",
			want: Repeat{V: 1},
		},
		{
			str:  "trigger when * repeat 1/5m", // short
			want: Repeat{V: 1, Interval: 5 * time.Minute},
		},
		{
			str:  "trigger when * repeat 10 times interval 5s",
			want: Repeat{V: 10, Interval: 5 * time.Second},
		},
		{
			str:  "trigger when * repeat 1 times",
			want: Repeat{V: 1},
		},
		{
			str: "trigger when * repeat 1 times interval",
			err: true,
		},
		{
			str: "trigger when * repeat 5.5 times",
			err: true,
		},
		{
			str: "trigger when * repeat every",
			err: true,
		},
		{
			str:  "trigger when * repeat;",
			want: DefaultRepeatVal,
		},
		{
			str: "trigger when * repeat 1 hoho",
			err: true,
		},
		{
			str: "trigger when * repeat 0x44 times",
			err: true,
		},
		{
			str: "trigger when * repeat 1/34moi",
			err: true,
		},
		{
			str: "trigger when * repeat 1 times someident",
			err: true,
		},
	}
	for i, tc := range testCases {
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
		if have, want := trigger.Repeat.V, tc.want.V; have != want {
			t.Fatalf("%d. have %d, want %d repeat times", i, have, want)
		}
		if have, want := trigger.Repeat.Interval, tc.want.Interval; have != want {
			t.Fatalf("%d. have %d, want %d repeat interval", i, have, want)
		}
	}
}

func TestParseTriggerStmtReset(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want time.Duration
	}{
		{
			str:  "trigger when * reset after 1h",
			want: 1 * time.Hour,
		},
		{
			str:  "trigger when * reset after 5m",
			want: 5 * time.Minute,
		},
		{
			str:  "trigger when * reset after 400s",
			want: 400 * time.Second,
		},
		{
			str: `
TRIGGER
VARS
   a={1,2,3}
   b=[1,2,3]
WHEN *
RESET AFTER 24h;
REPEAT
`,
			want: 24 * time.Hour,
		},
		{
			str: "trigger when * reset 1h",
			err: true,
		},
		{
			str: "trigger when * reset after 9MO",
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
		if have, want := trigger.Reset.V, tc.want; have != want {
			t.Fatalf("have %s, want %s", have, want)
		}
	}
}

func TestParseTriggerStmtFloatListVal(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want map[string]map[float64]struct{}
	}{
		{
			str: "trigger vars a={1.1,2.2,3.3,3.3,5.5} when *",
			want: map[string]map[float64]struct{}{
				"a": {
					1.1: {}, 2.2: {}, 3.3: {}, 5.5: {},
				},
			},
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
			t.Fatalf("have %d, want %d int vars", have, want)
		}
		for varname, want := range tc.want {
			vals, found := trigger.Vars[varname]
			if !found {
				t.Fatalf("variable not found %s", varname)
			}
			vars := vals.(ListFloatVal)
			if have, want := len(vars.V), len(want); have != want {
				t.Fatalf("have %d, want %d", have, want)
			}
			for a := range want {
				_, ok := vars.V[a]
				if !ok {
					t.Fatalf("list %s, item not found %f", varname, a)
				}
			}
		}
	}
}

func TestParseTriggerStmtIntListVal(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want map[string]map[int]struct{}
	}{
		{
			str: "trigger vars a={1,2,3,3,5} when *",
			want: map[string]map[int]struct{}{
				"a": {
					1: {}, 2: {}, 3: {}, 5: {},
				},
			},
		},
		{
			str: "trigger vars a={} when *",
		},
		{
			str: "trigger vars a={0x11} when *",
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
			t.Fatalf("have %d, want %d int vars", have, want)
		}
		for varname, want := range tc.want {
			vals, found := trigger.Vars[varname]
			if !found {
				t.Fatalf("variable not found %s", varname)
			}
			vars := vals.(ListIntVal)
			if have, want := len(vars.V), len(want); have != want {
				t.Fatalf("have %d, want %d", have, want)
			}
			for a := range want {
				_, ok := vars.V[a]
				if !ok {
					t.Fatalf("list %s, item not found %d", varname, a)
				}
			}
		}
	}
}

func TestParseTriggerStmtStringListVal(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want map[string]map[string]struct{}
	}{
		{
			str: "trigger vars a={\"70c960f3-4b56-4d71-a04a-2c62a714f4af\", \"one\"} when *",
			want: map[string]map[string]struct{}{
				"a": {
					"70c960f3-4b56-4d71-a04a-2c62a714f4af": {},
					"one":                                  {},
				},
			},
		},
		{
			str: "trigger vars a={\"one\", \"two\",,,,} when *",
			want: map[string]map[string]struct{}{
				"a": {
					"one": {},
					"two": {},
				},
			},
		},
		{
			str: "trigger vars a={\"one\", 1, \"two\"} when *",
			err: true,
		},
		{
			str: "trigger vars a={1, \"two\"} when *",
			err: true,
		},
		{
			str: "trigger vars a={1.1, \"one\", 1, \"two\"} when *",
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
