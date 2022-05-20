package geoqlparser

import (
	"testing"
)

func TestParseRadiusVal(t *testing.T) {
	testCases := []struct {
		str  string
		err  bool
		want RadiusVal
	}{
		{
			str:  "r10000m",
			want: RadiusVal{V: 10000, U: Meter},
		},
		{
			str:  "R1000000km",
			want: RadiusVal{V: 1000000, U: Kilometer},
		},
		{
			str:  "0",
			want: RadiusVal{},
		},
		{
			str: "",
			err: true,
		},
		{
			str: "r1000000000000000000000000000m",
			err: true,
		},
		{
			str: "some2234234",
			err: true,
		},
		{
			str: "rrrrrrr1111111m",
			err: true,
		},
		{
			str: "r5B",
			err: true,
		},
	}
	for _, tc := range testCases {
		have, err := toRadiusVal(tc.str)
		if tc.err {
			if err == nil {
				t.Fatalf("got nil, expected error")
			} else {
				continue
			}
		} else if !tc.err && err != nil {
			t.Fatal(err)
		}
		if have != tc.want {
			t.Fatalf("have %v, want %v", have, tc.want)
		}
	}
}
