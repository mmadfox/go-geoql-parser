package geoqlparser

import (
	"fmt"
	"testing"
)

func TestParseMonth(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name:   "valid array month",
			s:      `when month[jan, jul]`,
			assert: assertArray(MONTH, 1, 2, [][2]Pos{{5, 19}}),
		},
		{
			name:   "valid range month",
			s:      `when month[jan .. jul]`,
			assert: assertRange(MONTH, 1, [][2]Pos{{5, 21}}),
		},
		{
			name:   "valid month",
			s:      `when month[jan]`,
			assert: assertMonth([][1]int{{1}}, [][2]Pos{{5, 14}}),
		},
		{
			name: "invalid month",
			s:    `when month[mom]`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
	}
}

func TestParseWeekday(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name:   "valid array weekday",
			s:      `when weekday[mon, fri]`,
			assert: assertArray(WEEKDAY, 1, 2, [][2]Pos{{5, 21}}),
		},
		{
			name:   "valid range weekday",
			s:      `when weekday[mon .. fri]`,
			assert: assertRange(WEEKDAY, 1, [][2]Pos{{5, 23}}),
		},
		{
			name:   "valid weekday",
			s:      `when weekday[mon]`,
			assert: assertWeekday([][1]int{{1}}, [][2]Pos{{5, 16}}),
		},
		{
			name: "invalid weekday",
			s:    `when weekday[mom]`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
	}
}

func TestParseTimeLit(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name:   "valid time range time[9:12PM .. 9:12AM]",
			s:      `when time[9:12PM .. 9:12AM]`,
			assert: assertRange(TIME, 1, [][2]Pos{{5, 26}}),
		},
		{
			name:   "valid time array time[9:12PM,9:12AM,9:12:01PM]",
			s:      `when time[9:12PM, 9:12AM, 9:12:01PM]`,
			assert: assertArray(TIME, 1, 3, [][2]Pos{{5, 35}}),
		},
		{
			name:   "valid time time[9:12PM]",
			s:      `when time[9:12PM]`,
			assert: assertTime([][3]int{{9, 12, 0}}, PM, [][2]Pos{{5, 16}}),
		},
		{
			name:   "valid time time[11:01]",
			s:      `when time[11:01]`,
			assert: assertTime([][3]int{{11, 1, 0}}, 0, [][2]Pos{{5, 15}}),
		},
		{
			name:   "valid time time[11:01:01]",
			s:      `when time[11:01:01]`,
			assert: assertTime([][3]int{{11, 1, 1}}, 0, [][2]Pos{{5, 18}}),
		},
		{
			name: "invalid time format",
			s:    `when time`,
			err:  true,
		},
		{
			name: "invalid time format time[]",
			s:    `when time[]`,
			err:  true,
		},
		{
			name: "invalid time format time[1,1,1]",
			s:    `when time[1,1,1]`,
			err:  true,
		},
		{
			name: "invalid time format time[1-1]",
			s:    `when time["one" .. one]`,
			err:  true,
		},
		{
			name: "invalid time format time[1:1.1]",
			s:    `when time[1:1.1]`,
			err:  true,
		},
		{
			name: "invalid time format time[30:1:1]",
			s:    `when time[30:1:1]`,
			err:  true,
		},
		{
			name: "invalid time format time[23:60:1]",
			s:    `when time[23:60:1]`,
			err:  true,
		},
		{
			name: "invalid time format time[23:40:60]",
			s:    `when time[23:40:60]`,
			err:  true,
		},
		{
			name: "invalid time format",
			s:    `when time[23:40:60m]`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
	}
}

func TestParseDateLit(t *testing.T) {
	testCases := []parserTestCase1{
		{
			name: "valid range dates",
			s:    `when date[2022-01-01 .. 2022-01-01]`,
			assert: assertRange(DATE, 1, [][2]Pos{
				{5, 34},
			}),
		},
		{
			name: "valid array of dates",
			s:    `when date[2022-01-01, 2022-01-01, 2022-01-01]`,
			assert: assertArray(DATE, 1, 3, [][2]Pos{
				{5, 44},
			}),
		},
		{
			name: "valid date 2022-01-01",
			s:    `when date[2022-01-01] > 1`,
			assert: assertDate(
				[][3]int{{2022, 01, 01}}, [][2]Pos{{5, 20}}),
		},
		{
			name: "valid date 2022-01-01 and 2025-01-01",
			s:    `when date[2022-01-01] > 1 and date[2025-01-01] > 1`,
			assert: assertDate(
				[][3]int{
					{2022, 01, 01},
					{2025, 01, 01},
				}, [][2]Pos{
					{5, 20},
					{30, 45},
				}),
		},
		{
			name: "valid range dates",
			s:    `when date[2022-01-01 .. 2022-01-01, 2022-01-01 ]`,
			err:  true,
		},
		{
			name: "valid range dates",
			s:    `when date[,2022-01-01 .. 2022-01-01]`,
			err:  true,
		},
		{
			name: "valid range dates",
			s:    `when date[ .. 2022-01-01]`,
			err:  true,
		},
		{
			name: "invalid array of dates",
			s:    `when date[2022-01-01 ..  2022-01-01, 2022-01-01]`,
			err:  true,
		},
		{
			name: "invalid date format",
			s:    `when  date[2044/01/01] eq 1`,
			err:  true,
		},
		{
			name: "invalid date format",
			s:    `when  date[2044-1.1-1.1] eq 1`,
			err:  true,
		},
		{
			name: "invalid year 2000",
			s:    `when date[2000-01-01] > 1`,
			err:  true,
		},
		{
			name: "invalid month 13",
			s:    `when date[2032-13-01] > 1`,
			err:  true,
		},
		{
			name: "invalid month 0",
			s:    `when date[2032-00-01] > 1`,
			err:  true,
		},
		{
			name: "invalid day 32",
			s:    `when date[2032-13-32] > 1`,
			err:  true,
		},
		{
			name: "invalid day 00",
			s:    `when date[2032-13-00] > 1`,
			err:  true,
		},
		{
			name: "invalid date format: date[2032-13]",
			s:    `when date[2032-13] > 1`,
			err:  true,
		},
		{
			name: "invalid date format: date[]",
			s:    `when date[] > 1`,
			err:  true,
		},
		{
			name: "invalid date format: date",
			s:    `when date > 1`,
			err:  true,
		},
		{
			name: "invalid date format: date[",
			s:    `when  date[`,
			err:  true,
		},
	}
	for _, tc := range testCases {
		runAndTestTriggerStmt(t, tc)
	}
}

func assertMonth(expect [][1]int, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t.When, func(expr Expr) bool {
			month, ok := expr.(*MonthTyp)
			if !ok {
				return true
			}
			wi := expect[found]
			pos := positions[found]
			if have, want := month.Val, wi[0]; have != want {
				err = fmt.Errorf("got %d, want %d *MonthTyp.Val", have, want)
				return false
			}
			if have, want := month.lpos, pos[0]; have != want {
				err = fmt.Errorf("got %d, want %d *MonthTyp.lpos", have, want)
				return false
			}
			if have, want := month.rpos, pos[1]; have != want {
				err = fmt.Errorf("got %d, want %d *MonthTyp.rpos", have, want)
				return false
			}
			found++
			return false
		})
		if err == nil && found == 0 {
			err = fmt.Errorf("no test found")
		}
		return
	}
}

func assertWeekday(expect [][1]int, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t.When, func(expr Expr) bool {
			weekday, ok := expr.(*WeekdayTyp)
			if !ok {
				return true
			}
			wi := expect[found]
			pos := positions[found]
			if have, want := weekday.Val, wi[0]; have != want {
				err = fmt.Errorf("got %d, want %d WeekdayLit.Val", have, want)
				return false
			}
			if have, want := weekday.lpos, pos[0]; have != want {
				err = fmt.Errorf("got %d, want %d WeekdayLit.lpos", have, want)
				return false
			}
			if have, want := weekday.rpos, pos[1]; have != want {
				err = fmt.Errorf("got %d, want %d WeekdayLit.rpos", have, want)
				return false
			}
			found++
			return false
		})
		if err == nil && found == 0 {
			err = fmt.Errorf("no test found")
		}
		return
	}
}

func assertTime(expect [][3]int, unit Unit, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t.When, func(expr Expr) bool {
			time_, ok := expr.(*TimeTyp)
			if !ok {
				return true
			}
			wi := expect[found]
			pos := positions[found]
			if have, want := time_.Hours, wi[0]; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.Hours", have, want)
				return false
			}
			if have, want := time_.Minutes, wi[1]; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.Minutes", have, want)
				return false
			}
			if have, want := time_.Seconds, wi[2]; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.Seconds", have, want)
				return false
			}
			if have, want := time_.U, unit; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.Unit", have, want)
				return false
			}
			if have, want := time_.lpos, pos[0]; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.lpos", have, want)
				return false
			}
			if have, want := time_.rpos, pos[1]; have != want {
				err = fmt.Errorf("got %d, want %d TimeTyp.rpos", have, want)
				return false
			}
			found++
			return false
		})
		if err == nil && found == 0 {
			err = fmt.Errorf("no test found")
		}
		return
	}
}

func assertDate(expect [][3]int, positions [][2]Pos) func(t *Trigger) (err error) {
	return func(t *Trigger) (err error) {
		var found int
		Visit(t.When, func(expr Expr) bool {
			date, ok := expr.(*DateTyp)
			if !ok {
				return true
			}
			wi := expect[found]
			pos := positions[found]
			if have, want := date.Year, wi[0]; have != want {
				err = fmt.Errorf("got %d, want %d DateTyp.Year", have, want)
				return false
			}
			if have, want := date.Month, wi[1]; have != want {
				err = fmt.Errorf("got %d, want %d DateTyp.MonthTyp", have, want)
				return false
			}
			if have, want := date.Day, wi[2]; have != want {
				err = fmt.Errorf("got %d, want %d DateTyp.Day", have, want)
				return false
			}
			if have, want := date.lpos, pos[0]; have != want {
				err = fmt.Errorf("got %d, want %d DateTyp.lpos", have, want)
				return false
			}
			if have, want := date.rpos, pos[1]; have != want {
				err = fmt.Errorf("got %d, want %d DateTyp.rpos", have, want)
				return false
			}
			found++
			return false
		})
		if err == nil && found == 0 {
			err = fmt.Errorf("no test found")
		}
		return
	}
}
