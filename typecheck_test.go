package geoqlparser

import (
	"fmt"
	"testing"
)

var describeSelectors Dictionary

func init() {
	describeSelectors = Dict()
	describeSelectors["s_int"] = Int
	describeSelectors["s_float"] = Float
	describeSelectors["s_string"] = String
	describeSelectors["s_bool"] = Boolean
	describeSelectors["s_int_arr"] = ArrayInt
	describeSelectors["s_float_arr"] = ArrayFloat
	describeSelectors["s_string_arr"] = ArrayString
}

type checkSpec struct {
	n    string
	l, r string
	e    bool
	v    string
}

type checkTestCase struct {
	name string
	ops  []Token
	want []checkSpec
	err  bool
}

func checkAndTest(t *testing.T, name, str string, isErr bool) {
	stmt, err := Parse(str)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckType(stmt, describeSelectors)
	if err != nil {
		if isErr {
			return
		}
		t.Fatal(err)
	}
	if isErr && err == nil {
		t.Fatalf("%s [%s] got nil, expected error ", name, str)
	}
}

func TestCheck(t *testing.T) {
	testCases := []checkTestCase{
		{
			name: "aOps",
			ops:  []Token{REM, QUO, MUL, SUB},
			want: []checkSpec{
				{n: "int_int_ok", l: "1", r: "1"},
				{n: "int_float_ok", l: "1", r: "2.2"},
				{n: "int_sint", l: "1", r: "s_int"},
				{n: "int_rint", l: "1", r: "@ref", v: "1"},
				{n: "int_bool_bad", l: "1", r: "false", e: true},
				{n: "int_string_bad", l: "1", r: `"on"`, e: true},
				{n: "int_sint", l: "1", r: "selector_not_described", e: true},

				{n: "float_int_ok", l: "1.1", r: "2"},
				{n: "float_float_ok", l: "1.1", r: "2"},
				{n: "float_sfloat", l: "1.1", r: "s_float"},
				{n: "float_rfloat", l: "1.1", r: "@ref", v: "1.1"},
				{n: "float_string_bad", l: "1.1", r: `"on"`, e: true},
				{n: "float_bool_bad", l: "1.1", r: "true", e: true},
			},
		},
		{
			name: "equalOps",
			ops:  []Token{EQL, LEQL, NOT_EQ, LNEQ},
			want: []checkSpec{
				{n: "int_int_ok", l: "1", r: "1"},
				{n: "int_float_ok", l: "1", r: "2.2"},
				{n: "int_sint", l: "1", r: "s_int"},
				{n: "int_rint", l: "1", r: "@ref", v: "1"},
				{n: "int_bool_bad", l: "1", r: "false", e: true},
				{n: "int_string_bad", l: "1", r: `"on"`, e: true},
				{n: "int_sint", l: "1", r: "selector_not_described", e: true},

				{n: "float_int_ok", l: "1.1", r: "2"},
				{n: "float_float_ok", l: "1.1", r: "2"},
				{n: "float_sfloat", l: "1.1", r: "s_float"},
				{n: "float_rfloat", l: "1.1", r: "@ref", v: "1.1"},
				{n: "float_string_bad", l: "1.1", r: `"on"`, e: true},
				{n: "float_bool_bad", l: "1.1", r: "true", e: true},

				{n: "string_string_ok", l: `"on"`, r: `"on"`},
				{n: "string_sstring_ok", l: `"on"`, r: `s_string`},
				{n: "string_rstring_ok", l: `"on"`, r: "@ref", v: `"some string"`},
				{n: "string_int_bad", l: `"on"`, r: "1", e: true},
				{n: "string_float_bad", l: `"on"`, r: "1.1", e: true},
				{n: "string_bool_bad", l: `"on"`, r: "true", e: true},

				{n: "bool_bool_ok", l: "true", r: "true"},
				{n: "bool_sbool_ok", l: "true", r: "s_bool"},
				{n: "bool_rbool_ok", l: "true", r: "@ref", v: "true"},
				{n: "bool_int_bad", l: "true", r: "1", e: true},
				{n: "bool_float_bad", l: "true", r: "1.1", e: true},
				{n: "bool_string_bad", l: "true", r: `"on"`, e: true},

				{n: "astring_astring_ok", l: `["on"]`, r: `["on"]`},

				{n: "aint_aint_ok", l: "[1,2]", r: "[1,2]"},
				{n: "aint_afloat_ok", l: "[1,2]", r: "[1.1,2.1]"},
				{n: "aint_geometry_ok", l: "[1,2]", r: "point[1.1,2.1]"},

				{n: "afloat_aint_ok", l: "[1.1,2.1]", r: "[1,2]"},
				{n: "afloat_afloat_ok", l: "[1.1,2.1]", r: "[1.1,2.1]"},
				{n: "afloat_geometry_ok", l: "[1.1,2.1]", r: "point[1.1,2.1]"},

				{n: "geometry_aint_ok", l: "point[1.1,2.1]", r: "[1,2]"},
				{n: "geometry_afloat_ok", l: "point[1.1,2.1]", r: "[1.1,2.1]"},
				{n: "geometry_geometry_ok", l: "point[1.1,2.1]", r: "point[1.1,2.1]"},

				{n: "astring_aselector_ok", l: `["one", "two"]`, r: "s_string_arr"},
				{n: "astring_astring_ok", l: `["one", "two"]`, r: `["one", "two"]`},
				{n: "astring_aref_ok", l: `["one", "two"]`, r: "@ref", v: `["one", "two"]`},
				{n: "astring_string_ok", l: `["one", "two"]`, r: `"one"`, e: true},
				{n: "astring_aref_bed", l: `["one", "two"]`, r: "@ref", v: "some_selector", e: true},
				{n: "astring_aint_bad", l: `["one", "two"]`, r: `[1, 2, 3]`, e: true},
				{n: "astring_afloat_bad", l: `["one", "two"]`, r: `[1.1, 2.2]`, e: true},
				{n: "astring_float_bad", l: `["one", "two"]`, r: `2.2`, e: true},
				{n: "astring_int_bad", l: `["one", "two"]`, r: `2`, e: true},
				{n: "astring_geometry_bad", l: `["one", "two"]`, r: "point[1,2]", e: true},
				{n: "astring_rstring_bad", l: `["one", "two"]`, r: `"one" .. "two"`, e: true},
				{n: "astring_bool_bad", l: `["one", "two"]`, r: `true`, e: true},

				{n: "aint_afloat_ok", l: "[1,2]", r: "[1.1,2.2]"},
				{n: "aint_aint_ok", l: "[1,2]", r: "[1,2]"},
				{n: "aint_vint_ok", l: "[1,2]", r: "@ref", v: "[1,2]"},
				{n: "aint_sint_ok", l: "[1,2]", r: "s_int"},
				{n: "aint_rint_bad", l: "[1,2]", r: "1 .. 1", e: true},
				{n: "aint_float_bad", l: "[1,2]", r: "1.1", e: true},
				{n: "aint_string_bad", l: "[1,2]", r: `"one"`, e: true},
				{n: "aint_bool_bad", l: "[1,2]", r: `true`, e: true},
				{n: "aint_int_bad", l: "[1,2]", r: "1", e: true},
			},
		},
		{
			name: "inOps",
			ops:  []Token{IN, NOT_IN},
			want: []checkSpec{
				{n: "int_aint_ok", l: "1", r: "[1,2,3]"},
				{n: "int_int_ok", l: "1", r: "1 .. 2"},
				{n: "int_float_ok", l: "1", r: "1.1 .. 2.2"},
				{n: "int_point_ok", l: "1", r: "point[1,3]"},
				{n: "aint_int_bad", l: "[1,2,3]", r: "1", e: true},
				{n: "geometry_geometry_ok", l: "point[1,1]", r: "point[1,1]"},
			},
		},
		{
			name: "addOps",
			ops:  []Token{ADD},
			want: []checkSpec{
				{n: "int_int_ok", l: "1", r: "1"},
				{n: "int_float_ok", l: "1", r: "2.2"},
				{n: "int_sint", l: "1", r: "s_int"},
				{n: "int_rint", l: "1", r: "@ref", v: "1"},
				{n: "int_bool_bad", l: "1", r: "false", e: true},
				{n: "int_string_bad", l: "1", r: `"on"`, e: true},
				{n: "int_sint", l: "1", r: "selector_not_described", e: true},

				{n: "float_int_ok", l: "1.1", r: "2"},
				{n: "float_float_ok", l: "1.1", r: "2"},
				{n: "float_sfloat", l: "1.1", r: "s_float"},
				{n: "float_rfloat", l: "1.1", r: "@ref", v: "1.1"},
				{n: "float_string_bad", l: "1.1", r: `"on"`, e: true},
				{n: "float_bool_bad", l: "1.1", r: "true", e: true},

				{n: "string_string_ok", l: `"on"`, r: `"on"`},
				{n: "string_sstring_ok", l: `"on"`, r: `s_string`},
				{n: "string_rstring_ok", l: `"on"`, r: "@ref", v: `"some string"`},
				{n: "string_int_bad", l: `"on"`, r: "1", e: true},
				{n: "string_float_bad", l: `"on"`, r: "1.1", e: true},
				{n: "string_bool_bad", l: `"on"`, r: "true", e: true},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < len(tc.ops); i++ {
				op := tc.ops[i]
				for _, spec := range tc.want {
					specName := fmt.Sprintf("%s[%s]", spec.n, KeywordString(op))
					t.Run(specName, func(t *testing.T) {
						var bs string
						if len(spec.v) > 0 {
							bs = fmt.Sprintf("set ref=%s;", spec.v)
						}
						str := fmt.Sprintf("trigger %s  when %s %s %s",
							bs, spec.l, KeywordString(op), spec.r)
						checkAndTest(t, spec.n, str, spec.e)
					})
				}
			}
		})
	}
}
