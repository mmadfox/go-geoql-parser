package geoqlparser

import (
	"fmt"
	"testing"
)

type parserTestCase1 struct {
	name   string
	s      string
	err    bool
	assert func(t *TriggerStmt) (err error)
}

func TestParseVars(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "assign geometry",
			s: `trigger 
					set
						a = point[-1.1, 1.1];
						b = multipoint[point[-1.1, 1.1], point[-1.1, 1.1]];
						c = line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5]];
						d = polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]];
						x = point[-1.1, 1.1]:12km;
						u = multiline[
								line[[1.1,1.1], [1.1,1.1], [1.1,1.1]], 
								line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
						];
						o = multipolygon[
								polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]], 
								polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
						];
						m = collection[
					        point[-1.1, 1.1],
							multipoint[point[-1.1, 1.1], point[-1.1, 1.1]],
							line[[1.1, 1.1], [2.1, 3.1], [3.1, 5.5], [5.5, 5.5]],
							polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
							multiline[
								line[[1.1,1.1], [1.1,1.1], [1.1,1.1]], 
								line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
						    ],
							multipolygon[
								polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]], 
								polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
							]
						];
					when 
						@a intersects tracker_coords 
						and @b intersects tracker_coords 
						and @d intersects tracker_coords 
						and @u intersects tracker_coords 
						and @m intersects tracker_coords 
						and @o intersects tracker_coords 
						and @x intersects tracker_coords 
						and @c intersects tracker_coords`,
			assert: assertVars(),
		},
		{
			name:   "assign int",
			s:      `trigger set a=1; b=2; c=100; when @a > 100 and @b < 10 or @c == 300`,
			assert: assertVars(),
		},
		{
			name:   "assign float",
			s:      `trigger set a=1.1; b=2.1; c=-100.9; when @a > 100 and @b < 10 or @c == 300`,
			assert: assertVars(),
		},
		{
			name:   "assign string",
			s:      `trigger set a= "some text"; b="a"; when @a in "yes" and @b == "no"`,
			assert: assertVars(),
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
			name: "start pos + end pos",
			s:    `when selector{"one",*}:1,2,3,4 > 100 and selector2 in 1 .. 5000`,
			assert: func(t *TriggerStmt) (err error) {
				var found int
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *SelectorExpr:
						switch typ.Ident {
						case "selector2":
							if want, have := Pos(41), typ.Pos(); want != have {
								err = fmt.Errorf("got %d, expected %d pos", have, want)
								return false
							}
							if want, have := Pos(50), typ.End(); want != have {
								err = fmt.Errorf("got %d, expected %d pos", have, want)
								return false
							}
							found++
						case "selector":
							if want, have := Pos(5), typ.Pos(); want != have {
								err = fmt.Errorf("got %d, expected %d pos", have, want)
								return false
							}
							if want, have := Pos(30), typ.End(); want != have {
								err = fmt.Errorf("got %d, expected %d pos", have, want)
								return false
							}
							found++
						}
					}
					return true
				})
				if found != 2 {
					err = fmt.Errorf("got %d, expected 2 selector", found)
				}
				return
			},
		},
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

func TestParseTriggerStmtWhen(t *testing.T) {
	testCases := []struct {
		str string
		err bool
	}{
		{
			str: `
TRIGGER
SET
	somepoly = polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	somepoint = point[1.1,1.1]
	somecircle = point[1.2, 3.3]:500M
	someline  = line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	someline2  = line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]:1km
	somemultiline = multiline[
        line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	]
	somemultipoly = multipolygon[
	   	polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	]
    somemultipoint = multipoint[
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1]:400M
    ]
	somecoll = collection[
        multipoint[
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1]:400M
    ],
multipolygon[
	   	polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
		polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]]
	],
multiline[
        line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
		line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]
	],
polygon[[[1.1,1.1], [1.1,1.1], [1.1,1.1]], [[1.1,1.1], [1.1,1.1], [1.1,1.1]]],
point[1.1,1.1],
line[[1.1,1.1], [1.1,1.1], [1.1,1.1]],
line[[1.1,1.1], [1.1,1.1], [1.1,1.1]]:1km
    ]
WHEN
	tracker_point3 % 2 == 0
	and tracker_cords intersects @someplace
	and tracker_point1 / tracker_point2 * 100 > 20%
	and tracker_week in Sun .. Fri
	and tracker_time in 9:01AM .. 12:12PM 
	and tracker_temperature in 12Bar .. 44Psi
	and (tracker_speed in 10kph .. 40kph
	or tracker_speed in [10kph .. 40kph, 10kph .. 40kph, 10kph .. 40kph])
repeat 5 times 10s interval 
reset after 1h 
`,
		},
	}
	// TODO:
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
		trigger := stmt.(*TriggerStmt)
		_ = trigger
	}
}

func assertVars() func(t *TriggerStmt) (err error) {
	return func(t *TriggerStmt) (err error) {
		vars := make(map[string]struct{})
		Visit(t, func(expr Expr) bool {
			varlit, ok := expr.(*VarLit)
			if !ok {
				return true
			}
			vars[varlit.ID] = struct{}{}
			return true
		})
		if len(vars) == 0 {
			err = fmt.Errorf("no variables found")
		}
		for id := range vars {
			_, ok := t.Set[id]
			if !ok {
				err = fmt.Errorf("var %s not assigned", id)
				return
			}
		}
		return
	}
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
				case WEEKDAY, MONTH:
					_, ok = item.(*CalendarLit)
				case RANGE:
					_, ok = item.(*RangeExpr)
				case IDENT:
					_, ok = item.(*VarLit)
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
				case TIME:
					_, ok = item.(*TimeLit)
				case DATE:
					_, ok = item.(*DateLit)
				case DATETIME:
					_, ok = item.(*DateTimeLit)
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
