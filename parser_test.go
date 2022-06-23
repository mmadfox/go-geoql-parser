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
	somecircle = point[1.2, 3.3]:500m
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
point[1.1,1.1]:400m
    ]
	somecoll = collection[
        multipoint[
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1],
point[1.1,1.1]:400m
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
