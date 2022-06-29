package geoqlparser

import (
	"testing"
)

var testDict Dictionary

func init() {
	testDict = Dict()
	testDict["selector1_int_typ"] = Int
	testDict["selector2_int_typ"] = Int
	testDict["float_typ"] = Float
	testDict["string_typ"] = String
	testDict["array_int_typ"] = ArrayInt
	testDict["array_float_typ"] = ArrayFloat
	testDict["array_string_typ"] = ArrayString
	testDict["date_time_typ"] = DateTime
}

func TestCheck(t *testing.T) {
	testCases := []struct {
		s   string
		err bool
	}{
		{
			s: `when (15 / (7 - (1 + 1)) * 3 - (2 + (1 + 1))) >= 40 `,
		},
		{
			s: `when selector1_int_typ+selector2_int_typ > 300`,
		},
		{
			s:   `when selector1_int_typ+selector2_int_typ > "string"`,
			err: true,
		},
		{
			s: `when 1+1 > 400`,
		},
		{
			s:   `when 1+1 > "one"`,
			err: true,
		},
	}
	for i, tc := range testCases {
		stmt, err := Parse(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		err = Check(stmt, testDict)
		if tc.err {
			if err == nil {
				t.Fatalf("testCase-%d: got nil, expected error", i)
			} else {
				return
			}
		} else if !tc.err && err != nil {
			t.Fatal(err)
		}
	}
}
