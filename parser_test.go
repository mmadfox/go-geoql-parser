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
	assert func(t *Trigger) (err error)
}

func TestParseRepeatAndReset(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "repeat 1 every 1s",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 1 every 1s`,
			assert: func(t *Trigger) (err error) {
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
			assert: func(t *Trigger) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil Trigger.RepeatInterval", t.RepeatInterval)
				}
				if t.RepeatCount != nil {
					err = fmt.Errorf("got %T, expected nil Trigger.RepeatCount", t.RepeatCount)
				}
				if t.ResetAfter != nil {
					err = fmt.Errorf("got %T, expected nil Trigger.ResetAfter", t.ResetAfter)
				}
				return
			},
		},
		{
			name: "reset after 24h",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar reset after 24h`,
			assert: func(t *Trigger) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil Trigger.RepeatInterval", t.RepeatInterval)
				}
				if t.RepeatCount != nil {
					err = fmt.Errorf("got %T, expected nil Trigger.RepeatCount", t.RepeatCount)
				}
				return assertDuration(t.ResetAfter, 24*time.Hour)
			},
		},
		{
			name: "repeat 1 reset after 24h",
			s:    `trigger when tracker_osi*tracker_miu >= 300Bar repeat 1 reset after 24h`,
			assert: func(t *Trigger) (err error) {
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
			assert: func(t *Trigger) (err error) {
				if t.RepeatInterval != nil {
					err = fmt.Errorf("got %T, expected nil", t.RepeatInterval)
				}
				return assertInt(t.RepeatCount, 1)
			},
		},
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
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
		runAndTestTriggerStmt(t, tc)
	}
}

func assertVars(positions map[string][5]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		vars := make(map[string]struct{})
		Visit(t, func(expr Expr) bool {
			varlit, ok := expr.(*Ref)
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
			assignStmt, err := t.findAssign(id)
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
				return fmt.Errorf("var %s: got %d, want %d Assign.Left.Pos()", id, have, want)
			}
			// ident right pos
			if have, want := assignStmt.Left.End(), posSet[1]; have != want {
				return fmt.Errorf("var %s: got %d, want %d Assign.Left.End()", id, have, want)
			}
			// assign pos
			if have, want := assignStmt.TokPos, posSet[2]; have != want {
				return fmt.Errorf("var %s: got %d, want %d Assign.TokPos", id, have, want)
			}
			// right operand left pos
			if have, want := assignStmt.Right.Pos(), posSet[3]; have != want {
				return fmt.Errorf("var %s: got %d, want %d Assign.Right.Pos()", id, have, want)
			}
			// right operand left end
			if have, want := assignStmt.Right.End(), posSet[4]; have != want {
				return fmt.Errorf("var %s: got %d, want %d Assign.Right.End()", id, have, want)
			}
		}
		return
	}
}

func assertInt(expr Expr, val int) error {
	lit, ok := expr.(*IntTyp)
	if !ok {
		return fmt.Errorf("got %T, expected *IntTyp", expr)
	}
	if lit.Val != val {
		return fmt.Errorf("got %d, expected %d", lit.Val, val)
	}
	return nil
}

func assertDuration(expr Expr, val time.Duration) error {
	lit, ok := expr.(*DurationTyp)
	if !ok {
		return fmt.Errorf("got %T, expected *IntTyp", expr)
	}
	if lit.Val != val {
		return fmt.Errorf("got %d, expected %d", lit.Val, val)
	}
	return nil
}

func assertRange(kind Token, match int, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t, func(expr Expr) bool {
			rangeExpr, ok := expr.(*Range)
			if !ok {
				return true
			}
			if !isType(kind, rangeExpr.Low) {
				err = fmt.Errorf("*Range.Low: %T", rangeExpr.Low)
				return false
			}
			if !isType(kind, rangeExpr.High) {
				err = fmt.Errorf("*Range.High: %T", rangeExpr.High)
				return false
			}
			pos := positions[found]
			if have, want := rangeExpr.Pos(), pos[0]; have != want {
				err = fmt.Errorf("*Range.lpos: got %d, expected %d", have, want)
				return false
			}
			if have, want := rangeExpr.End(), pos[1]; have != want {
				err = fmt.Errorf("*Range.rpos: got %d, expected %d", have, want)
				return false
			}
			found++
			return true
		})
		if err == nil && found != match {
			err = fmt.Errorf("got %d, expected %d *Range",
				found, match)
		}
		return
	}
}

func assertArray(kind Token, match int, totalElements int, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t, func(expr Expr) bool {
			array, ok := expr.(*ArrayTyp)
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
				if !isType(kind, item) {
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

func runAndTestTriggerStmt(t *testing.T, tc parserTestCase1) {
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
		trigger := stmt.(*Trigger)
		if tc.assert != nil {
			if err = tc.assert(trigger); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func isType(kind Token, expr Expr) (ok bool) {
	switch kind {
	case DATE:
		_, ok = expr.(*DateTyp)
	case TIME:
		_, ok = expr.(*TimeTyp)
	case WEEKDAY:
		_, ok = expr.(*WeekdayTyp)
	case MONTH:
		_, ok = expr.(*MonthTyp)
	case RANGE:
		_, ok = expr.(*Range)
	case IDENT:
		_, ok = expr.(*Ref)
	case PERCENT:
		_, ok = expr.(*PercentTyp)
	case DISTANCE:
		_, ok = expr.(*DistanceTyp)
	case TEMPERATURE:
		_, ok = expr.(*TemperatureTyp)
	case PRESSURE:
		_, ok = expr.(*PressureTyp)
	case SPEED:
		_, ok = expr.(*SpeedTyp)
	case DURATION:
		_, ok = expr.(*DurationTyp)
	case STRING:
		_, ok = expr.(*StringTyp)
	case FLOAT:
		_, ok = expr.(*FloatTyp)
	case INT:
		_, ok = expr.(*IntTyp)
	}
	return
}
