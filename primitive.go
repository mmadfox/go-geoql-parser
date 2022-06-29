package geoqlparser

import (
	"strconv"
	"time"
)

func (s *parser) parseVarExpr() (expr Expr, err error) {
	s.next()
	expr = &Ref{ID: s.t.TokenText(), lpos: s.t.Offset()}
	s.next()
	return
}

func (s *parser) parseStringLit() (expr Expr, err error) {
	text := s.t.TokenText()
	expr = &StringTyp{
		Val:  trim(text),
		lpos: s.t.Offset(),
		rpos: s.t.Offset() + Pos(len(text)-1),
	}
	s.next()
	return
}

func (s *parser) parseFloatTypes() (expr Expr, err error) {
	s.lpos = s.t.Offset()
	val, err := strconv.ParseFloat(s.lit, 64)
	if err != nil {
		return nil, s.error()
	}
	expr, err = s.parseAllTypes(val)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		if s.isSignMinus() {
			val = -val
		}
		if s.isSignPlus() || s.isSignMinus() {
			s.lpos -= 1
		}
		expr = &FloatTyp{Val: val, lpos: s.lpos, rpos: s.t.Offset() - 1}
	}
	return
}

func (s *parser) parseIntTypes() (expr Expr, err error) {
	s.lpos = s.t.Offset()
	intval, err := strconv.Atoi(s.lit)
	if err != nil {
		return nil, s.error()
	}
	expr, err = s.parseAllTypes(float64(intval))
	if err != nil {
		return nil, err
	}
	if expr == nil {
		if s.isSignMinus() {
			intval = -intval
		}
		if s.isSignPlus() || s.isSignMinus() {
			s.lpos -= 1
		}
		expr = &IntTyp{Val: intval, lpos: s.lpos, rpos: s.t.Offset() - 1}
	}
	return
}

func (s *parser) parseAllTypes(v float64) (expr Expr, err error) {
	plit := s.lit
	s.rpos = s.t.Offset() - 1
	s.next()
	unit := s.t.TokenText()
	litlen := Pos(len(plit))
	switch {
	case isPercentUnit(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		s.rpos += litlen + 1
		expr = &PercentTyp{Val: v, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isPressureUnit(unit):
		u := unitFromString(s.lit)
		s.rpos += litlen + u.size()
		expr = &PressureTyp{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isDistanceUnit(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		u := unitFromString(unit)
		s.rpos += litlen + u.size()
		expr = &DistanceTyp{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isSpeedUnit(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		u := unitFromString(s.lit)
		s.rpos += litlen + u.size()
		expr = &SpeedTyp{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isTemperatureUnit(unit):
		var ts Sign
		if s.isSignMinus() {
			ts = Minus
		}
		if s.isSignPlus() {
			ts = Plus
		}
		if s.isSignPlus() || s.isSignMinus() {
			s.lpos -= 1
		}
		u := unitFromString(s.lit)
		s.rpos += litlen + u.size()
		expr = &TemperatureTyp{Val: v, U: u, Vec: ts, lpos: s.lpos, rpos: s.rpos}
		s.next()
	default:
		var fr rune
		if len(s.lit) > 0 {
			fr = rune(s.lit[0])
		}
		switch fr {
		case 'h', 'm', 's':
			if s.isSignMinus() {
				s.err = errNegativeValue
				return nil, s.error()
			}
			dur, er := time.ParseDuration(plit + s.lit)
			if er != nil {
				return nil, s.error()
			}
			s.rpos = s.lpos + Pos(len(plit+s.lit)-1)
			expr = &DurationTyp{Val: dur, lpos: s.lpos, rpos: s.rpos}
			s.next()
		}
	}
	return
}

func (s *parser) parseDistance() (dist *DistanceTyp, err error) {
	s.next()
	var ok bool
	switch s.tok {
	case INT:
		re, err := s.parseIntTypes()
		if err != nil {
			return nil, err
		}
		dist, ok = re.(*DistanceTyp)
		if !ok {
			return nil, s.error()
		}
	case FLOAT:
		re, err := s.parseFloatTypes()
		if err != nil {
			return nil, err
		}
		dist, ok = re.(*DistanceTyp)
		if !ok {
			return nil, s.error()
		}
	}
	return
}

func (s *parser) parseBooleanLit() (expr Expr, err error) {
	switch s.lit {
	default:
		return nil, s.error()
	case "true", "up":
		expr = &BooleanTyp{Val: true, lpos: s.t.Offset()}
	case "false", "down":
		expr = &BooleanTyp{Val: false, lpos: s.t.Offset()}
	}
	s.next()
	return
}
