package geoqlparser

import (
	"fmt"
	"strconv"
)

func (s *parser) parseDateExpr() (expr Expr, err error) {
	lpos := s.t.Offset()
	s.next()
	if !s.except(LBRACK) {
		s.err = fmt.Errorf("invalid date format: got date without body, expected date[YYYY-MM-DD]")
		return nil, s.error()
	}
	var arrayExpr *ArrayTyp
	var rangeExpr *Range
	isDay := func(n int) bool { return n == 3 }
	isMonth := func(n int) bool { return n == 1 }
	isArrayTyp := func() bool { return arrayExpr != nil }
	isRangeTyp := func() bool { return rangeExpr != nil }
	for {
		var y, m, d int
		s.next()
		if s.except(EOF, RBRACK) {
			if y == 0 {
				s.err = fmt.Errorf("invalid date format: got date[], expected date[YYYY-MM-DD]")
				return nil, s.error()
			}
			break
		}
		if !s.except(INT, SUB, COMMA, RANGE) {
			return nil, s.error()
		}
		y, err = strconv.Atoi(s.lit)
		if err != nil {
			s.err = err
			return nil, s.error()
		}
		var sep int
		for i := 0; i <= 3; i++ {
			s.next()
			if s.except(SUB) {
				sep++
				continue
			}
			if !s.except(SUB, INT) || s.except(EOF) {
				err = fmt.Errorf("invalid date format: expected date[YYYY-MM-DD]")
			}
			switch {
			case isMonth(i):
				m, err = strconv.Atoi(s.lit)
			case isDay(i):
				d, err = strconv.Atoi(s.lit)
			}
			if err != nil {
				break
			}
		}
		if err != nil {
			s.err = err
			return nil, s.error()
		}
		if d < 1 || d > 31 {
			s.err = fmt.Errorf("invalid day format: got %d, expected 1-31", d)
			return nil, s.error()
		}
		if m < 1 || m > 12 {
			s.err = fmt.Errorf("invalid month format: got %d, expected 1-12", m)
			return nil, s.error()
		}
		if y < 2022 || y > 2200 {
			s.err = fmt.Errorf("invalid year format: got %d, expected 2022-2200", y)
			return nil, s.error()
		}
		if sep < 1 {
			s.err = fmt.Errorf("invalid date format")
			return nil, s.error()
		}
		s.next()
		date := &DateTyp{
			Year:  y,
			Month: m,
			Day:   d,
			lpos:  lpos,
			rpos:  s.t.Offset(),
		}
		switch s.tok {
		case COMMA:
			if rangeExpr != nil {
				return nil, s.error()
			}
			if arrayExpr == nil {
				arrayExpr = &ArrayTyp{Kind: DATE, List: make([]Expr, 0)}
			}
			arrayExpr.List = append(arrayExpr.List, date)
		case RANGE, PERIOD:
			if arrayExpr != nil {
				return nil, s.error()
			}
			if rangeExpr == nil {
				rangeExpr = &Range{}
			}
			rangeExpr.Low = date
		}
		if s.except(RBRACK) {
			switch {
			default:
				s.next()
				return date, nil
			case isArrayTyp():
				arrayExpr.List = append(arrayExpr.List, date)
				arrayExpr.lpos = lpos
				arrayExpr.rpos = s.t.Offset()
				s.next()
				return arrayExpr, nil
			case isRangeTyp():
				rangeExpr.High = date
				rangeExpr.lpos = lpos
				rangeExpr.rpos = s.t.Offset()
				s.next()
				return rangeExpr, nil
			}
		}
	}
	return
}

func (s *parser) parseTimeExpr() (expr Expr, err error) {
	lpos := s.t.Offset()
	s.next()
	if !s.except(LBRACK) {
		s.err = fmt.Errorf("invalid time format: got time without body, expected time[HH:MM:SS] or time[HH:MM]")
		return nil, s.error()
	}
	var arrayExpr *ArrayTyp
	var rangeExpr *Range
	isSeconds := func(n int) bool { return n == 3 }
	isMinutes := func(n int) bool { return n == 1 }
	isArrayTyp := func() bool { return arrayExpr != nil }
	isRangeTyp := func() bool { return rangeExpr != nil }
	for {
		var h, m, c int
		s.next()
		if s.except(EOF, RBRACK) {
			if h == 0 {
				s.err = fmt.Errorf("invalid time format: got time[], expected time[HH:MM:SS]")
				return nil, s.error()
			}
			break
		}
		if !s.except(INT, COLON, COMMA, RANGE) {
			return nil, s.error()
		}
		h, err = strconv.Atoi(s.lit)
		if err != nil {
			s.err = err
			return nil, s.error()
		}
		var sep int
		for i := 0; i <= 3; i++ {
			s.next()
			if s.except(COLON) {
				sep++
				continue
			}
			if s.except(RBRACK, COMMA, RANGE, SELECTOR) {
				s.t.Reset()
				break
			}
			if !s.except(COLON, INT) || s.except(EOF) {
				err = fmt.Errorf("invalid time format: expected time[HH:MM:SS]")
			}
			switch {
			case isMinutes(i):
				m, err = strconv.Atoi(s.lit)
			case isSeconds(i):
				c, err = strconv.Atoi(s.lit)
			}
			if err != nil {
				break
			}
		}
		if err != nil {
			s.err = err
			return nil, s.error()
		}
		if h < 0 || h > 24 {
			s.err = fmt.Errorf("invalid hour: got %d, expected 0-24", h)
			return nil, s.error()
		}
		if m < 0 || m > 59 {
			s.err = fmt.Errorf("invalid minutes: got %d, expected 0-59", m)
			return nil, s.error()
		}
		if c < 0 || c > 59 {
			s.err = fmt.Errorf("invalid seconds: got %d, expected 0-59", c)
			return nil, s.error()
		}
		if sep < 1 {
			s.err = fmt.Errorf("invalid time format")
			return nil, s.error()
		}
		s.next()
		time_ := &TimeTyp{
			Hours:   h,
			Minutes: m,
			Seconds: c,
			lpos:    lpos,
			rpos:    s.t.Offset(),
		}

		if s.except(SELECTOR) {
			if !isTimeUnit(s.lit) {
				s.err = fmt.Errorf("invalid unit time: got %s, expected AM or PM", s.lit)
				return nil, s.error()
			}
			time_.U = unitFromString(s.lit)
			time_.rpos = s.t.Offset() + time_.U.size()
			s.next()
		}

		switch s.tok {
		case COMMA:
			if rangeExpr != nil {
				return nil, s.error()
			}
			if arrayExpr == nil {
				arrayExpr = &ArrayTyp{Kind: TIME, List: make([]Expr, 0)}
			}
			arrayExpr.List = append(arrayExpr.List, time_)
		case RANGE, PERIOD:
			if arrayExpr != nil {
				return nil, s.error()
			}
			if rangeExpr == nil {
				rangeExpr = &Range{}
			}
			rangeExpr.Low = time_
		}
		if s.except(RBRACK) {
			switch {
			default:
				s.next()
				return time_, nil
			case isArrayTyp():
				arrayExpr.List = append(arrayExpr.List, time_)
				arrayExpr.lpos = lpos
				arrayExpr.rpos = s.t.Offset()
				s.next()
				return arrayExpr, nil
			case isRangeTyp():
				rangeExpr.High = time_
				rangeExpr.lpos = lpos
				rangeExpr.rpos = s.t.Offset()
				s.next()
				return rangeExpr, nil
			}
		}
	}
	return
}

func (s *parser) parseWeekdayExpr() (expr Expr, err error) {
	lpos := s.t.Offset()
	s.next()
	if !s.except(LBRACK) {
		s.err = fmt.Errorf("invalid weekday format: got without body, expected weekday[Mon]")
		return nil, s.error()
	}
	var arrayExpr *ArrayTyp
	var rangeExpr *Range
	isArrayTyp := func() bool { return arrayExpr != nil }
	isRangeTyp := func() bool { return rangeExpr != nil }
	for {
		var weekday int
		s.next()
		if s.except(EOF, RBRACK) {
			if weekday == 0 {
				s.err = fmt.Errorf("invalid weekday format: got weekday[], expected weekday[Mon]")
				return nil, s.error()
			}
			break
		}
		if !s.except(SELECTOR) {
			return nil, s.error()
		}
		switch s.lit {
		default:
			s.err = fmt.Errorf("invalid weekday format got %s, expected weekday[Mon, ...]", s.lit)
			return nil, s.error()
		case "sun":
			weekday = 0
		case "mon":
			weekday = 1
		case "tue":
			weekday = 2
		case "wed":
			weekday = 3
		case "thu":
			weekday = 4
		case "fri":
			weekday = 5
		case "sat":
			weekday = 6
		}
		s.next()
		wd := &WeekdayTyp{Val: weekday, lpos: lpos, rpos: s.t.Offset()}
		switch s.tok {
		case COMMA:
			if rangeExpr != nil {
				return nil, s.error()
			}
			if arrayExpr == nil {
				arrayExpr = &ArrayTyp{Kind: WEEKDAY, List: make([]Expr, 0)}
			}
			arrayExpr.List = append(arrayExpr.List, wd)
		case RANGE, PERIOD:
			if arrayExpr != nil {
				return nil, s.error()
			}
			if rangeExpr == nil {
				rangeExpr = &Range{}
			}
			rangeExpr.Low = wd
		}
		if s.except(RBRACK) {
			switch {
			default:
				s.next()
				return wd, nil
			case isArrayTyp():
				arrayExpr.List = append(arrayExpr.List, wd)
				arrayExpr.lpos = lpos
				arrayExpr.rpos = s.t.Offset()
				s.next()
				return arrayExpr, nil
			case isRangeTyp():
				rangeExpr.High = wd
				rangeExpr.lpos = lpos
				rangeExpr.rpos = s.t.Offset()
				s.next()
				return rangeExpr, nil
			}
		}
	}
	return
}

func (s *parser) parseMonthExpr() (expr Expr, err error) {
	lpos := s.t.Offset()
	s.next()
	if !s.except(LBRACK) {
		s.err = fmt.Errorf("invalid month format: got without body, expectedl month[Jan]")
		return nil, s.error()
	}
	var arrayExpr *ArrayTyp
	var rangeExpr *Range
	isArrayTyp := func() bool { return arrayExpr != nil }
	isRangeTyp := func() bool { return rangeExpr != nil }
	for {
		var month int
		s.next()
		if s.except(EOF, RBRACK) {
			if month == 0 {
				s.err = fmt.Errorf("invalid month format: got month[], expected month[Jan]")
				return nil, s.error()
			}
			break
		}
		if !s.except(SELECTOR) {
			return nil, s.error()
		}
		switch s.lit {
		default:
			s.err = fmt.Errorf("invalid month format got %s, expected month[Mon, ...]", s.lit)
			return nil, s.error()
		case "jan":
			month = 1
		case "feb":
			month = 2
		case "mar":
			month = 3
		case "apr":
			month = 4
		case "may":
			month = 5
		case "jun":
			month = 6
		case "jul":
			month = 7
		case "aug":
			month = 8
		case "sep":
			month = 9
		case "oct":
			month = 10
		case "nov":
			month = 11
		case "dec":
			month = 12
		}
		s.next()
		wd := &MonthTyp{Val: month, lpos: lpos, rpos: s.t.Offset()}
		switch s.tok {
		case COMMA:
			if rangeExpr != nil {
				return nil, s.error()
			}
			if arrayExpr == nil {
				arrayExpr = &ArrayTyp{Kind: MONTH, List: make([]Expr, 0)}
			}
			arrayExpr.List = append(arrayExpr.List, wd)
		case RANGE, PERIOD:
			if arrayExpr != nil {
				return nil, s.error()
			}
			if rangeExpr == nil {
				rangeExpr = &Range{}
			}
			rangeExpr.Low = wd
		}
		if s.except(RBRACK) {
			switch {
			default:
				s.next()
				return wd, nil
			case isArrayTyp():
				arrayExpr.List = append(arrayExpr.List, wd)
				arrayExpr.lpos = lpos
				arrayExpr.rpos = s.t.Offset()
				s.next()
				return arrayExpr, nil
			case isRangeTyp():
				rangeExpr.High = wd
				rangeExpr.lpos = lpos
				rangeExpr.rpos = s.t.Offset()
				s.next()
				return rangeExpr, nil
			}
		}
	}
	return
}
