package geoqlparser

import (
	"fmt"
	"testing"
	"time"
)

type parserTestCase1 struct {
	name   string
	s      string
	err    bool
	assert func(t *TriggerStmt) (err error)
}

func TestParseRepeatAndReset(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "repeat 1 every 1s",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 1 every 1s`,
			assert: func(t *TriggerStmt) (err error) {
				if err = assertDuration(t.RepeatInterval, 1*time.Second); err != nil {
					return err
				}
				return assertInt(t.RepeatCount, 1)
			},
		},
		{
			name: "error: repeat 2.2",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 2.2`,
			err:  true,
		},
		{
			name: "error: repeat 30s",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 30s`,
			err:  true,
		},
		{
			name: "empty repeat",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat`,
		},
		{
			name: "without repeat and reset",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar`,
			assert: func(t *TriggerStmt) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil TriggerStmt.RepeatInterval", t.RepeatInterval)
				}
				if t.RepeatCount != nil {
					err = fmt.Errorf("got %T, expected nil TriggerStmt.RepeatCount", t.RepeatCount)
				}
				if t.ResetAfter != nil {
					err = fmt.Errorf("got %T, expected nil TriggerStmt.ResetAfter", t.ResetAfter)
				}
				return
			},
		},
		{
			name: "reset after 24h",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar reset after 24h`,
			assert: func(t *TriggerStmt) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil TriggerStmt.RepeatInterval", t.RepeatInterval)
				}
				if t.RepeatCount != nil {
					err = fmt.Errorf("got %T, expected nil TriggerStmt.RepeatCount", t.RepeatCount)
				}
				return assertDuration(t.ResetAfter, 24*time.Hour)
			},
		},
		{
			name: "repeat 1 reset after 24h",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 1 reset after 24h`,
			assert: func(t *TriggerStmt) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil", t.RepeatInterval)
				}
				if err = assertDuration(t.ResetAfter, 24*time.Hour); err != nil {
					return err
				}
				return assertInt(t.RepeatCount, 1)
			},
		},
		{
			name: "repeat 1",
			s:    `trigger when tracker_osi*tracker_miu > 300 repeat 1`,
			assert: func(t *TriggerStmt) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil", t.RepeatInterval)
				}
				return assertInt(t.RepeatCount, 1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := Parse(tc.s)
			if tc.err {
				if err == nil {
					t.Fatalf("%s: got nil, expected error", tc.name)
				} else {
					return
				}
			} else if !tc.err && err != nil {
				t.Fatal(err)
			}
			trigger := stmt.(*TriggerStmt)
			if tc.assert == nil {
				return
			}
			if err := tc.assert(trigger); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParseVars(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "assign selector with args and props",
			s:    `trigger set a=selector{"one","two"}:1km;  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 38},
			}),
		},
		{
			name: "assign selector with args",
			s:    `trigger set a=selector{"one","two"};  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 34},
			}),
		},
		{
			name: "assign array of range",
			s:    `trigger set a=[1 .. 2, 5 .. 9];  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 29},
			}),
		},
		{
			name: "assign array of time",
			s:    `trigger set a=[11:11:11, 10:10:10];  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 33},
			}),
		},
		{
			name: "assign geometry polygon",
			s:    `trigger set a=polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]];  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 90},
			}),
		},
		{
			name: "assign geometry line",
			s:    `trigger set a=line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5]];  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 65},
			}),
		},
		{
			name: "assign geometry point",
			s:    `trigger set a=point[-1.1, 1.1];  when @a`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 29},
			}),
		},
		{
			name: "assign duration",
			s:    `trigger set a=7h3m45s  b=3m  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 20},
				"b": {23, 23, 24, 25, 26},
			}),
		},
		{
			name: "assign negative duration",
			s:    `trigger set a=-7h3m45s  b=3m  when @a == @b`,
			err:  true,
		},
		{
			name: "assign temperature",
			s:    `trigger set a=+3C  b=-55F  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 16},
				"b": {19, 19, 20, 21, 24},
			}),
		},
		{
			name: "assign temperature without sign",
			s:    `trigger set a=3C  b=0F  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 15},
				"b": {18, 18, 19, 20, 21},
			}),
		},
		{
			name: "assign distance",
			s:    `trigger set a=300km  b=4M  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 18},
				"b": {21, 21, 22, 23, 24},
			}),
		},
		{
			name: "assign negative distance",
			s:    `trigger set a=-300km  b=4M  when @a == @b`,
			err:  true,
		},
		{
			name: "assign pressure",
			s:    `trigger set a=12Bar  b=40000Psi  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 18},
				"b": {21, 21, 22, 23, 30},
			}),
		},
		{
			name: "assign dateTime",
			s:    `trigger set a=2030-10-02T11:11:11  b=2030-10-02T11:11:11  when @a == @b`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 32},
				"b": {35, 35, 36, 37, 55},
			}),
		},
		{
			name: "assign negative dateTime",
			s:    `trigger set a=-2030-10-02T11:11:11  b=2030-10-02T11:11:11  when @a == @b`,
			err:  true,
		},
		{
			name: "assign percent",
			s:    `trigger set aaaa=10%; b=0.1%; when @aaaa == @b`,
			assert: assertVars(map[string][5]Pos{
				"aaaa": {12, 15, 16, 17, 19},
				"b":    {22, 22, 23, 24, 27},
			}),
		},
		{
			name: "assign speed",
			s:    `trigger set aaaa=1Kph; bbbb=30Mph; when @aaaa == @bbbb`,
			assert: assertVars(map[string][5]Pos{
				"aaaa": {12, 15, 16, 17, 20},
				"bbbb": {23, 26, 27, 28, 32},
			}),
		},
		{
			name: "assign negative speed",
			s:    `trigger set aaaa=-1Kph; bbbb=-30Mph; when @aaaa == @bbbb`,
			err:  true,
		},
		{
			name: "assign string",
			s:    `trigger set aaaa="some string"; bbbb="bbbb"; when @aaaa == @bbbb`,
			assert: assertVars(map[string][5]Pos{
				"aaaa": {12, 15, 16, 17, 29},
				"bbbb": {32, 35, 36, 37, 42},
			}),
		},
		{
			name: "assign float",
			s:    `trigger set aa=1.1; bbb=200.1; c=-3123.345768; when @aa > 100 and @bbb < 10 or @c == 300`,
			assert: assertVars(map[string][5]Pos{
				"aa":  {12, 13, 14, 15, 17},
				"bbb": {20, 22, 23, 24, 28},
				"c":   {31, 31, 32, 33, 44},
			}),
		},
		{
			name: "assign int",
			s:    `trigger set a=1; b=2; c=-100; when @a > 100 and @b < 10 or @c == 300`,
			assert: assertVars(map[string][5]Pos{
				"a": {12, 12, 13, 14, 14},
				"b": {17, 17, 18, 19, 19},
				"c": {22, 22, 23, 24, 27},
			}),
		},
		{
			name: "assign duplicate",
			s:    `trigger set a="some text"; a="a"; when @a in "yes" and @b == "no"`,
			err:  true,
		},
		{
			name: "assign variable",
			s:    `trigger set a=@a;  when *`,
			err:  true,
		},
		{
			name: "assign variable",
			s:    `trigger set a=[@a, @b];  when *`,
			err:  true,
		},
		{
			name: "assign variable",
			s:    `trigger set a:1;  when *`,
			err:  true,
		},
		{
			name: "assign variable",
			s:    `trigger set a=1`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := Parse(tc.s)
			if tc.err {
				if err == nil {
					t.Fatalf("%s: got nil, expected error", tc.name)
				} else {
					return
				}
			} else if !tc.err && err != nil {
				t.Fatal(err)
			}
			trigger := stmt.(*TriggerStmt)
			if err := tc.assert(trigger); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "array with different types",
			s:    `when selector in [1, 30Km]`,
			err:  true,
		},
		{
			name: "illegal array",
			s:    `when selector in [}]`,
			err:  true,
		},
		{
			name: "empty array",
			s:    `when selector in []`,
			err:  true,
		},
		{
			name: "array of array",
			s:    `when selector in [[1,2], [2,3]]`,
			err:  true,
		},
		{
			name:   "array of months",
			s:      `when selector in [Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec ]`,
			assert: assertArray(MONTH, 1, 12, [][2]Pos{{17, 77}}),
		},
		{
			name:   "array of weekdays",
			s:      `when selector in [Sun, Mon, Tue, Wed, Thu, Fri, Sat]`,
			assert: assertArray(WEEKDAY, 1, 7, [][2]Pos{{17, 51}}),
		},
		{
			name:   "array of ranges",
			s:      `when selector in [1 .. 1, 2 .. 10]`,
			assert: assertArray(RANGE, 1, 2, [][2]Pos{{17, 33}}),
		},
		{
			name:   "array of vars",
			s:      `when selector in [@somevar, @somevar2, @somevar3]`,
			assert: assertArray(IDENT, 1, 3, [][2]Pos{{17, 48}}),
		},
		{
			name:   "array of percents",
			s:      `when selector in [100%, 0.001%]`,
			assert: assertArray(PERCENT, 1, 2, [][2]Pos{{17, 30}}),
		},
		{
			name:   "array of distance",
			s:      `when selector in [100M, 5Km]`,
			assert: assertArray(DISTANCE, 1, 2, [][2]Pos{{17, 27}}),
		},
		{
			name:   "array of temperature",
			s:      `when selector in [19C, 30F]`,
			assert: assertArray(TEMPERATURE, 1, 2, [][2]Pos{{17, 26}}),
		},
		{
			name:   "array of pressure",
			s:      `when selector in [50Psi, 1Bar]`,
			assert: assertArray(PRESSURE, 1, 2, [][2]Pos{{17, 29}}),
		},
		{
			name:   "array of speed",
			s:      `when selector in [50mph, 1kph]`,
			assert: assertArray(SPEED, 1, 2, [][2]Pos{{17, 29}}),
		},
		{
			name:   "array of duration",
			s:      `when selector in [1h, 20s, 7h3m45s, 7h3m, 3m]`,
			assert: assertArray(DURATION, 1, 5, [][2]Pos{{17, 44}}),
		},
		{
			name:   "array of int",
			s:      `when selector in [1,2,3,-4]`,
			assert: assertArray(INT, 1, 4, [][2]Pos{{17, 26}}),
		},
		{
			name:   "array of float",
			s:      `when selector in [1.1,22.2,-3.0,1.4]`,
			assert: assertArray(FLOAT, 1, 4, [][2]Pos{{17, 35}}),
		},
		{
			name:   "array of string",
			s:      `when selector in ["one", "two"]`,
			assert: assertArray(STRING, 1, 2, [][2]Pos{{17, 30}}),
		},
		{
			name:   "array of dateTime",
			s:      `when selector in [2030-10-02T11:11:11, 2060-10-02T11:11:11]`,
			assert: assertArray(DATETIME, 1, 2, [][2]Pos{{17, 58}}),
		},
		{
			name:   "array of date",
			s:      `when selector in [2030-10-02, 2030-10-02 , 2030-10-02 ]`,
			assert: assertArray(DATE, 1, 3, [][2]Pos{{17, 54}}),
		},
		{
			name:   "array of time",
			s:      `when selector in [11:11:11, 11:11, 9:11AM, 3:04Pm ]`,
			assert: assertArray(TIME, 1, 4, [][2]Pos{{17, 50}}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := Parse(tc.s)
			if tc.err {
				if err == nil {
					t.Fatalf("got nil, expected error")
				} else {
					return
				}
			} else if !tc.err && err != nil {
				t.Fatal(err)
			}
			trigger := stmt.(*TriggerStmt)
			if err := tc.assert(trigger); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParseSelector(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "with wildcard",
			s:    "when selector{*}",
			assert: func(t *TriggerStmt) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "selector":
							if !typ.Wildcard {
								err = fmt.Errorf("got false, expected true wildcard flag")
								return false
							}
						}
					}
					return true
				})
				return
			},
		},
		{
			name: "with args and distances props",
			s:    `when selector{"one", "two"}:1km,3km,6km`,
			assert: func(t *TriggerStmt) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "selector":
							if len(typ.Props) != 3 {
								err = fmt.Errorf("got %d, expected 3 selector prop",
									len(typ.Props))
								return false
							}
						}
					}
					return true
				})
				return
			},
		},
		{
			name: "without args with only distance props",
			s:    "when selector:1km",
			assert: func(t *TriggerStmt) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "selector":
							if !typ.Wildcard {
								err = fmt.Errorf("got %v, expected true",
									typ.Wildcard)
								return false
							}
							if len(typ.Props) != 1 {
								err = fmt.Errorf("got %d, expected 1 selector prop",
									len(typ.Props))
								return false
							}
							dist, ok := typ.Props[0].(*DistanceLit)
							if !ok {
								err = fmt.Errorf("got %T, expected *DistanceLit",
									typ.Props[0])
								return false
							}
							if dist.Val != 1 {
								err = fmt.Errorf("got 1, expected %f%s",
									dist.Val, dist.U)
								return false
							}
						}
					}
					return true
				})
				return
			},
		},
		{
			name: "args and wildcard",
			s:    `when selector{*, "one", "two"}`,
			assert: func(t *TriggerStmt) (err error) {
				var foundSelectors int
				var foundArgs int
				var wildcard bool
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "selector":
							wildcard = typ.Wildcard
							foundSelectors++
							foundArgs = len(typ.Args)
						}
					}
					return true
				})
				if foundSelectors != 1 {
					err = fmt.Errorf("got %d, expected 2 slectors", foundSelectors)
				}
				if foundArgs != 2 {
					err = fmt.Errorf("got %d, expected 2 slectors arguments", foundArgs)
				}
				if !wildcard {
					err = fmt.Errorf("got false, expected true for slectors wildcard")
				}
				return
			},
		},
		{
			name: "only ident",
			s:    "when trigger_one > 0 and trigger_two < 0",
			assert: func(t *TriggerStmt) (err error) {
				var found int
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "trigger_one", "trigger_two":
							found++
						}
					}
					return true
				})
				if found != 2 {
					err = fmt.Errorf("got %d, expected 2 slectors", found)
				}
				return
			},
		},
		{
			name: "illegal args",
			s:    "when some_selector{1km, 2km}",
			err:  true,
		},
		{
			name: "illegal props with short form",
			s:    "when selector:#",
			err:  true,
		},
		{
			name: "illegal props",
			s:    `when selector{"one", "two"}:>,>`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := Parse(tc.s)
			if tc.err {
				if err == nil {
					t.Fatalf("got nil, expected error")
				} else {
					return
				}
			} else if !tc.err && err != nil {
				t.Fatal(err)
			}
			trigger := stmt.(*TriggerStmt)
			if err := tc.assert(trigger); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertVars(positions map[string][5]Pos) func(t *TriggerStmt) (err error) {
	return func(t *TriggerStmt) (err error) {
		vars := make(map[string]struct{})
		Visit(t, func(expr Expr) bool {
			varlit, ok := expr.(*RefLit)
			if !ok {
				return true
			}
			vars[varlit.ID] = struct{}{}
			return true
		})
		if len(vars) == 0 {
			return fmt.Errorf("no variables found")
		}
		for id := range vars {
			assignStmt, err := t.findVar(id)
			if err != nil {
				return err
			}
			// test positions
			posSet, ok := positions[id]
			if !ok {
				return fmt.Errorf("can't find position for variable %s", id)
			}
			// ident left pos
			if have, want := assignStmt.Left.Pos(), posSet[0]; have != want {
				return fmt.Errorf("var %s: got %d, want %d AssignStmt.Left.Pos()", id, have, want)
			}
			// ident right pos
			if have, want := assignStmt.Left.End(), posSet[1]; have != want {
				return fmt.Errorf("var %s: got %d, want %d AssignStmt.Left.End()", id, have, want)
			}
			// assign pos
			if have, want := assignStmt.TokPos, posSet[2]; have != want {
				return fmt.Errorf("var %s: got %d, want %d AssignStmt.TokPos", id, have, want)
			}
			// right operand left pos
			if have, want := assignStmt.Right.Pos(), posSet[3]; have != want {
				return fmt.Errorf("var %s: got %d, want %d AssignStmt.Right.Pos()", id, have, want)
			}
			// right operand left end
			if have, want := assignStmt.Right.End(), posSet[4]; have != want {
				return fmt.Errorf("var %s: got %d, want %d AssignStmt.Right.End()", id, have, want)
			}
		}
		return
	}
}

func assertInt(expr Expr, val int) error {
	lit, ok := expr.(*IntLit)
	if !ok {
		return fmt.Errorf("got %T, expected *IntLit", expr)
	}
	if lit.Val != val {
		return fmt.Errorf("got %d, expected %d", lit.Val, val)
	}
	return nil
}

func assertDuration(expr Expr, val time.Duration) error {
	lit, ok := expr.(*DurationLit)
	if !ok {
		return fmt.Errorf("got %T, expected *IntLit", expr)
	}
	if lit.Val != val {
		return fmt.Errorf("got %d, expected %d", lit.Val, val)
	}
	return nil
}

func assertArray(kind Token, match int, totalElements int, positions [][2]Pos) func(t *TriggerStmt) (err error) {
	return func(t *TriggerStmt) (err error) {
		var found int
		Visit(t, func(expr Expr) bool {
			array, ok := expr.(*ArrayExpr)
			if !ok {
				return true
			}
			if have, want := array.Kind, kind; have != want {
				err = fmt.Errorf("got %s, expected %s type", KeywordString(have), KeywordString(want))
				return false
			}
			if have, want := len(array.List), totalElements; have != want {
				err = fmt.Errorf("got %d, expected %d array len", want, have)
				return false
			}
			for i := 0; i < len(array.List); i++ {
				item := array.List[i]
				var ok bool
				switch kind {
				case WEEKDAY, MONTH, DATE, DATETIME, TIME:
					_, ok = item.(*CalendarLit)
				case RANGE:
					_, ok = item.(*RangeExpr)
				case IDENT:
					_, ok = item.(*RefLit)
				case PERCENT:
					_, ok = item.(*PercentLit)
				case DISTANCE:
					_, ok = item.(*DistanceLit)
				case TEMPERATURE:
					_, ok = item.(*TemperatureLit)
				case PRESSURE:
					_, ok = item.(*PressureLit)
				case SPEED:
					_, ok = item.(*SpeedLit)
				case DURATION:
					_, ok = item.(*DurationLit)
				case STRING:
					_, ok = item.(*StringLit)
				case FLOAT:
					_, ok = item.(*FloatLit)
				case INT:
					_, ok = item.(*IntLit)
				}
				if !ok {
					err = fmt.Errorf("got %T, expected %s", item, KeywordString(kind))
					return false
				}
			}
			pos := positions[found]
			if have, want := array.Pos(), pos[0]; have != want {
				err = fmt.Errorf("got %d, expected %d start pos", have, want)
				return false
			}
			if have, want := array.End(), pos[1]; have != want {
				err = fmt.Errorf("got %d, expected %d end pos", have, want)
				return false
			}
			found++
			return true
		})
		if err == nil && found != match {
			err = fmt.Errorf("got %d, expected %d array of %s",
				found, match, KeywordString(kind))
		}
		return
	}
}
