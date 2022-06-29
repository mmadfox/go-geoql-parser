package geoqlparser

import (
	"fmt"
	"testing"
)

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
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
	}
}

func TestParseSelector(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "with wildcard",
			s:    "when selector{*}",
			assert: func(t *Trigger) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *Selector:
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
			assert: func(t *Trigger) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *Selector:
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
			assert: func(t *Trigger) (err error) {
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *Selector:
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
							dist, ok := typ.Props[0].(*DistanceTyp)
							if !ok {
								err = fmt.Errorf("got %T, expected *DistanceTyp",
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
			assert: func(t *Trigger) (err error) {
				var foundSelectors int
				var foundArgs int
				var wildcard bool
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *Selector:
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
			assert: func(t *Trigger) (err error) {
				var found int
				Visit(t, func(expr Expr) bool {
					switch typ := expr.(type) {
					case *Selector:
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
		runAndTestTriggerStmt(t, tc)
	}
}
