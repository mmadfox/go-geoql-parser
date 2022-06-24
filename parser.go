package geoqlparser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var errNegativeValue = errors.New("value cannot be negative")

func Parse(gql string) (Statement, error) {
	r := strings.NewReader(gql)
	t := NewTokenizer(r)
	s := newParser(t, r)
	return s.parse0()
}

type parser struct {
	r    *strings.Reader
	t    *Tokenizer
	tok  Token
	lit  string
	err  error
	sign Token
	lpos Pos
	rpos Pos
}

func (s *parser) next() {
	s.tok, s.lit = s.t.Scan()
}

func (s *parser) parse0() (stmt Statement, err error) {
	s.next()
	switch s.tok {
	case TRIGGER, WHEN:
		if s.except(TRIGGER) {
			s.next()
		}
		return s.parseTriggerStmt()
	default:
		err = s.error()
	}
	return
}

func (s *parser) parseTriggerStmt() (stmt *TriggerStmt, err error) {
	if !s.except(WHEN, SET) {
		return nil, s.error()
	}
	stmt = new(TriggerStmt)
	stmt.lpos = s.t.Offset()
	if s.except(SET) {
		if err = s.parseSet(stmt); err != nil {
			return nil, err
		}
		s.next()
		if !s.except(WHEN) {
			return nil, s.error()
		}
	}
	if s.except(WHEN) {
		if err = s.parseWhen(stmt); err != nil {
			return nil, err
		}
	}
	if s.except(REPEAT) {
		if err = s.parseRepeat(stmt); err != nil {
			return nil, err
		}
	}
	if s.except(RESET) {
		if err = s.parseReset(stmt); err != nil {
			return nil, err
		}
	}
	stmt.rpos = s.t.Offset()
	return stmt, nil
}

func (s *parser) except(in ...Token) bool {
	for i := 0; i < len(in); i++ {
		if s.tok == in[i] {
			return true
		}
	}
	return false
}

func (s *parser) parseReset(stmt *TriggerStmt) (err error) {
	s.next()
	s.next()
	var dur Expr
	switch s.tok {
	case INT:
		dur, err = s.parseIntTypes()
	case FLOAT:
		dur, err = s.parseFloatTypes()
	}
	if err != nil {
		return err
	}
	_, ok := dur.(*DurationLit)
	if !ok {
		return s.error()
	}
	stmt.ResetAfter = dur
	return
}

func (s *parser) parseRepeat(stmt *TriggerStmt) (err error) {
	s.next()
	if s.except(EOF, RESET) {
		return
	}

	if !s.except(INT) {
		return s.error()
	}

	repeatCount, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	if _, ok := repeatCount.(*IntLit); !ok {
		return s.error()
	}

	stmt.RepeatCount = repeatCount
	if s.except(EOF, RESET) {
		return
	}

	if !s.except(SELECTOR) {
		return s.error()
	}
	s.next()
	dur, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	if _, ok := dur.(*DurationLit); !ok {
		return s.error()
	}
	stmt.RepeatInterval = dur
	return
}

func (s *parser) parseSet(stmt *TriggerStmt) error {
	s.next()
	for {
		if s.except(WHEN) {
			s.t.Reset()
			break
		}
		if s.except(EOF) {
			break
		}
		if !s.except(SELECTOR) {
			return s.error()
		}
		ident := IdentLit{Val: s.t.TokenText(), lpos: s.t.Offset()}
		ident.rpos = s.t.Offset() + Pos(len(ident.Val)-1)
		s.t.Unwind()
		s.next()
		if !s.except(ASSIGN) {
			return s.error()
		}
		tokPos := s.t.Offset() - 1
		expr, err := s.parseUnaryExpr()
		if err != nil {
			return err
		}
		if s.except(SEMICOLON) {
			s.next()
		}
		switch typ := expr.(type) {
		case *RefLit:
			return s.error()
		case *ArrayExpr:
			if typ.Kind == IDENT {
				return s.error()
			}
		}
		stmt.initVars()
		va := &AssignStmt{
			Left:   &ident,
			TokPos: tokPos,
			Right:  expr,
		}
		if err := stmt.SetVar(va); err != nil {
			s.err = err
			return s.error()
		}
	}
	return nil
}

func (s *parser) parseWhen(stmt *TriggerStmt) error {
	if !s.except(WHEN) {
		return s.error()
	}
	expr, err := s.parseBinaryExpr(1)
	if err != nil {
		return err
	}
	stmt.When = expr
	return nil
}

func (s *parser) parseBinaryExpr(oprec0 int) (Expr, error) {
	left, err := s.parseUnaryExpr()
	if err != nil {
		return nil, err
	}
	if s.except(EOF) || s.except(RPAREN) {
		return left, nil
	}
	for {
		if s.except(RPAREN) {
			return left, nil
		}
		op, oprec, pos := s.tok, s.tok.Precedence(), s.t.Offset()
		if oprec < oprec0 {
			return left, nil
		}

		right, err := s.parseBinaryExpr(oprec + 1)
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Right: right, Op: op, OpPos: pos}
	}
}

func (s *parser) parseUnaryExpr() (expr Expr, err error) {
	if s.t.Err() != nil {
		s.err = s.t.Err()
		return nil, s.error()
	}

	s.next()

	if s.except(SUB, ADD) {
		switch s.tok {
		case SUB:
			s.sign = SUB
		case ADD:
			s.sign = ADD
		}
		s.next()
	}

	switch s.tok {
	default:
		s.err = fmt.Errorf("illegal expression")
		err = s.error()
	case SELECTOR:
		expr, err = s.parseSelectorExpr()
	case MUL:
		expr, err = s.parseWildcardLit()
	case IDENT:
		expr, err = s.parseVarExpr()
	case LBRACK:
		expr, err = s.parseArrayExpr()
	case LPAREN:
		expr, err = s.parseParenExpr()
	case FLOAT:
		expr, err = s.parseFloatTypes()
	case INT:
		expr, err = s.parseIntTypes()
	case STRING:
		expr, err = s.parseStringLit()
	case GEOMETRY_POINT, GEOMETRY_LINE, GEOMETRY_POLYGON:
		expr, err = s.parseGeometryExpr()
	case GEOMETRY_MULTIPOINT, GEOMETRY_MULTILINE, GEOMETRY_MULTIPOLYGON:
		expr, err = s.parseGeometryMultiObject()
	case GEOMETRY_COLLECTION:
		expr, err = s.parseGeometryCollectionExpr()
	case BOOLEAN:
		expr, err = s.parseBooleanLit()
	case WEEKDAY, MONTH:
		expr, err = s.parseCalendarLit()
	}
	if err == nil {
		switch s.tok {
		case RANGE:
			expr, err = s.parseRangeExpr(expr)
		}
	}
	s.resetSign()
	return
}

func (s *parser) parseVarExpr() (expr Expr, err error) {
	s.next()
	expr = &RefLit{ID: s.t.TokenText(), lpos: s.t.Offset()}
	s.next()
	return
}

func (s *parser) parseCalendarLit() (expr Expr, err error) {
	var val int
	switch s.tok {
	case WEEKDAY:
		switch s.lit {
		case "sun":
			val = 0
		case "mon":
			val = 1
		case "tue":
			val = 2
		case "wed":
			val = 3
		case "thu":
			val = 4
		case "fri":
			val = 5
		case "sat":
			val = 6
		}
	case MONTH:
		switch s.lit {
		case "jan":
			val = 1
		case "feb":
			val = 2
		case "mar":
			val = 3
		case "apr":
			val = 4
		case "may":
			val = 5
		case "jun":
			val = 6
		case "jul":
			val = 7
		case "aug":
			val = 8
		case "sep":
			val = 9
		case "oct":
			val = 10
		case "nov":
			val = 11
		case "dec":
			val = 12
		}
	}
	expr = &CalendarLit{Val: val, lpos: s.t.Offset(), Kind: s.tok}
	s.next()
	return
}

func (s *parser) parseRangeExpr(low Expr) (expr Expr, err error) {
	_, isRangeExpr := low.(*RangeExpr)
	if isRangeExpr {
		return nil, s.error()
	}
	startPos := s.t.Offset()
	high, err := s.parseUnaryExpr()
	if err != nil {
		return
	}
	return &RangeExpr{
		Low:  low,
		High: high,
		lpos: startPos,
		rpos: s.t.Offset(),
	}, nil
}

func (s *parser) parseWildcardLit() (expr Expr, err error) {
	return &WildcardLit{lpos: s.t.Offset()}, nil
}

func (s *parser) parseBooleanLit() (expr Expr, err error) {
	switch s.lit {
	default:
		return nil, s.error()
	case "true", "up":
		expr = &BooleanLit{Val: true, lpos: s.t.Offset()}
	case "false", "down":
		expr = &BooleanLit{Val: false, lpos: s.t.Offset()}
	}
	s.next()
	return
}

func (s *parser) parseArrayExpr() (expr Expr, err error) {
	if !s.except(LBRACK) {
		return nil, s.error()
	}

	var arrayExpr *ArrayExpr
	startPos := s.t.Offset()
	checkKind := func(k1, k2 Token) error {
		if k1 != ILLEGAL && k1 != k2 {
			return s.error()
		}
		return nil
	}

	for {
		expr, err = s.parseUnaryExpr()
		if err != nil {
			return nil, err
		}

		if s.except(COMMA) {
			s.t.Unwind()
		}

		if arrayExpr == nil {
			arrayExpr = &ArrayExpr{lpos: startPos, List: make([]Expr, 0), Kind: ILLEGAL}
		}

		switch typ := expr.(type) {
		default:
			err = s.error()
		case *DurationLit:
			err = checkKind(arrayExpr.Kind, DURATION)
			arrayExpr.Kind = DURATION
		case *SpeedLit:
			err = checkKind(arrayExpr.Kind, SPEED)
			arrayExpr.Kind = SPEED
		case *PressureLit:
			err = checkKind(arrayExpr.Kind, PRESSURE)
			arrayExpr.Kind = PRESSURE
		case *TemperatureLit:
			err = checkKind(arrayExpr.Kind, TEMPERATURE)
			arrayExpr.Kind = TEMPERATURE
		case *DistanceLit:
			err = checkKind(arrayExpr.Kind, DISTANCE)
			arrayExpr.Kind = DISTANCE
		case *PercentLit:
			err = checkKind(arrayExpr.Kind, PERCENT)
			arrayExpr.Kind = PERCENT
		case *IntLit:
			err = checkKind(arrayExpr.Kind, INT)
			arrayExpr.Kind = INT
		case *FloatLit:
			err = checkKind(arrayExpr.Kind, FLOAT)
			arrayExpr.Kind = FLOAT
		case *StringLit:
			err = checkKind(arrayExpr.Kind, STRING)
			arrayExpr.Kind = STRING
		case *RefLit:
			err = checkKind(arrayExpr.Kind, IDENT)
			arrayExpr.Kind = IDENT
		case *RangeExpr:
			err = checkKind(arrayExpr.Kind, RANGE)
			arrayExpr.Kind = RANGE
		case *CalendarLit:
			switch typ.Kind {
			case DATE:
				err = checkKind(arrayExpr.Kind, DATE)
				arrayExpr.Kind = DATE
			case TIME:
				err = checkKind(arrayExpr.Kind, TIME)
				arrayExpr.Kind = TIME
			case DATETIME:
				err = checkKind(arrayExpr.Kind, DATETIME)
				arrayExpr.Kind = DATETIME
			case WEEKDAY:
				err = checkKind(arrayExpr.Kind, WEEKDAY)
				arrayExpr.Kind = WEEKDAY
			case MONTH:
				err = checkKind(arrayExpr.Kind, MONTH)
				arrayExpr.Kind = MONTH
			}
		}
		if err != nil {
			return nil, err
		}

		arrayExpr.List = append(arrayExpr.List, expr)

		if s.except(RBRACK) {
			break
		}
	}
	arrayExpr.rpos = s.t.Offset()
	s.next()
	return arrayExpr, nil
}

func (s *parser) parseGeometryMultiObject() (expr Expr, err error) {
	geotyp := s.tok
	s.next()
	if !s.except(LBRACK) {
		return nil, s.error()
	}
	s.next()
	multiobj := &GeometryMultiObject{
		Kind: geotyp,
		Val:  make([]Expr, 0),
		lpos: s.t.Offset(),
	}
	for {
		if !isGeometryToken(s.tok) {
			break
		}
		object, oer := s.parseGeometryExpr()
		if oer != nil {
			return nil, oer
		}
		switch typ := object.(type) {
		default:
			return nil, s.error()
		case *GeometryPointExpr:
			if multiobj.Kind != GEOMETRY_MULTIPOINT {
				return nil, s.error()
			}
			multiobj.Val = append(multiobj.Val, typ)
		case *GeometryLineExpr:
			if multiobj.Kind != GEOMETRY_MULTILINE {
				return nil, s.error()
			}
			multiobj.Val = append(multiobj.Val, typ)
		case *GeometryPolygonExpr:
			if multiobj.Kind != GEOMETRY_MULTIPOLYGON {
				return nil, s.error()
			}
			multiobj.Val = append(multiobj.Val, typ)
		}
		if s.except(RBRACK) {
			s.next()
			break
		}
		if s.except(COMMA) {
			s.next()
		}
	}
	multiobj.rpos = s.t.Offset()
	if len(multiobj.Val) == 0 {
		return nil, s.error()
	}
	return multiobj, nil
}

func (s *parser) parseGeometryCollectionExpr() (expr Expr, err error) {
	collection := &GeometryCollectionExpr{
		Objects: make([]Expr, 0),
		lpos:    s.t.Offset(),
	}
	s.next()
	if !s.except(LBRACK) {
		return nil, s.error()
	}
	s.next()
	for {
		if !isGeometryToken(s.tok) {
			break
		}
		var object Expr
		var oer error
		switch s.tok {
		default:
			object, oer = s.parseGeometryExpr()
		case GEOMETRY_MULTIPOINT, GEOMETRY_MULTILINE, GEOMETRY_MULTIPOLYGON:
			object, oer = s.parseGeometryMultiObject()
		}
		if oer != nil {
			return nil, oer
		}
		collection.Objects = append(collection.Objects, object)
		if s.except(RBRACK) {
			s.next()
			break
		}
		if s.except(COMMA) {
			s.next()
		}
	}
	collection.rpos = s.t.Offset()
	if len(collection.Objects) == 0 {
		return nil, s.error()
	}
	return collection, nil
}

func (s *parser) parseGeometryExpr() (expr Expr, err error) {
	geomtyp := s.tok
	sp := s.t.Offset()
	s.next()
	if !s.except(LBRACK) {
		return nil, s.error()
	}
	var path int
	var x, y float64
	var pi uint8
	var aa [2]float64
	var bb [][2]float64
	var cc [][][2]float64
	open := true
	q := 1
	for {
		s.next()
		if s.except(EOF) {
			break
		}
		if !s.except(LBRACK, SUB, COMMA, FLOAT, INT, RBRACK) {
			break
		}
		if s.except(LBRACK) {
			path++
			q++
			open = true
			continue
		}
		if s.except(RBRACK) {
			if open {
				switch path {
				case 0:
					aa = [2]float64{x, y}
				case 1:
					if bb == nil {
						bb = make([][2]float64, 0)
					}
					bb = append(bb, [2]float64{x, y})
				case 2:
					if cc == nil {
						cc = make([][][2]float64, 0)
					}
					if bb == nil {
						bb = make([][2]float64, 0)
					}
					bb = append(bb, [2]float64{x, y})
				}
			}
			if path > 0 {
				path--
			}
			if path == 0 && bb != nil && cc != nil {
				if len(bb) > 0 {
					cc = append(cc, bb)
				}
				bb = make([][2]float64, 0)
			}
			open = false
			pi = 0
			x = 0
			y = 0
			q--
			if q <= 0 {
				s.next()
				break
			}
			continue
		}
		if s.except(COMMA) {
			if open {
				pi = 1
			}
			continue
		}
		if s.except(SUB) {
			s.sign = SUB
			s.next()
		}

		var val float64
		switch s.tok {
		case FLOAT:
			val, err = strconv.ParseFloat(s.lit, 64)
			if err != nil {
				s.err = err
				return nil, s.error()
			}
		case INT:
			ival, err := strconv.Atoi(s.lit)
			if err != nil {
				s.err = err
				return nil, s.error()
			}
			val = float64(ival)
		default:
			return nil, s.error()
		}

		if s.isSignMinus() {
			val = -val
		}
		switch pi {
		case 0:
			x = val
		case 1:
			y = val
		default:
			return nil, s.error()
		}
		s.resetSign()
	}
	switch geomtyp {
	case GEOMETRY_POINT:
		point := &GeometryPointExpr{Val: aa, lpos: sp, rpos: s.t.Offset() - 1}
		if !s.except(COLON) {
			return point, nil
		}
		// point with margin
		radius, err := s.parseDistance()
		if err != nil {
			return nil, err
		}
		point.Radius = radius
		point.rpos = s.t.Offset()
		return point, nil
	case GEOMETRY_LINE:
		line := &GeometryLineExpr{Val: bb, lpos: sp, rpos: s.t.Offset() - 1}
		if !s.except(COLON) {
			return line, nil
		}
		// line with margin
		margin, err := s.parseDistance()
		if err != nil {
			return nil, err
		}
		line.Margin = margin
		line.rpos = s.t.Offset()
		return line, nil
	case GEOMETRY_POLYGON:
		return &GeometryPolygonExpr{Val: cc, lpos: sp, rpos: s.t.Offset() - 1}, nil
	}
	err = s.error()
	return
}

func (s *parser) resetSign() {
	s.sign = ILLEGAL
}

func (s *parser) isSignMinus() bool {
	return s.sign == SUB
}

func (s *parser) isSignPlus() bool {
	return s.sign == ADD
}

func (s *parser) parseDistance() (dist *DistanceLit, err error) {
	s.next()
	var ok bool
	switch s.tok {
	case INT:
		re, err := s.parseIntTypes()
		if err != nil {
			return nil, err
		}
		dist, ok = re.(*DistanceLit)
		if !ok {
			return nil, s.error()
		}
	case FLOAT:
		re, err := s.parseFloatTypes()
		if err != nil {
			return nil, err
		}
		dist, ok = re.(*DistanceLit)
		if !ok {
			return nil, s.error()
		}
	}
	return
}

func (s *parser) parseParenExpr() (expr Expr, err error) {
	lp := s.t.Offset()
	expr, err = s.parseBinaryExpr(s.tok.Precedence())
	if err != nil {
		return nil, err
	}
	rp := s.t.Offset()
	s.next()
	return &ParenExpr{Expr: expr, lpos: lp, rpos: rp}, nil
}

func (s *parser) parseStringLit() (expr Expr, err error) {
	text := s.t.TokenText()
	expr = &StringLit{
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
		expr = &FloatLit{Val: val, lpos: s.lpos, rpos: s.t.Offset() - 1}
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
		expr = &IntLit{Val: intval, lpos: s.lpos, rpos: s.t.Offset() - 1}
	}
	return
}

func (s *parser) parseDateTime(prefix string) (expr Expr, err error) {
	lpos := s.t.Offset() - Pos(len(prefix))
	in, err := strconv.Atoi(prefix)
	if err != nil {
		s.err = err
		return nil, s.error()
	}
	var (
		val   int
		year  int
		month time.Month
		day   int
		hour  int
		min   int
		sec   int
	)
	var isOnlyTime bool
	index := 1
	step := 8
	if s.lit == ":" {
		step = 3
		hour = in
		isOnlyTime = true
	} else {
		year = in
	}
	for i := 0; i < step; i++ {
		s.next()
		if s.except(EOF) {
			break
		}
		if len(s.lit) > 0 && s.lit[0] == 't' {
			s.lit = s.lit[1:]
			s.tok = INT
		}
		if !s.except(INT, COLON, SUB) {
			s.t.Reset()
			break
		}
		if s.except(SUB, COLON) {
			continue
		}
		s.rpos = s.t.Offset() + 1
		val, err = strconv.Atoi(s.lit)
		if err != nil {
			s.err = err
			return nil, s.error()
		}
		if isOnlyTime {
			switch index {
			case 1:
				min = val
			case 2:
				sec = val
			}
		} else {
			switch index {
			case 1:
				month = time.Month(val)
			case 2:
				day = val
			case 3:
				hour = val
			case 4:
				min = val
			case 5:
				sec = val
			}
		}
		index++
	}

	switch isOnlyTime {
	case true:
		switch index {
		default:
			err = s.error()
		case 3, 2:
			expr = &CalendarLit{
				Kind:    TIME,
				Hours:   hour,
				Minutes: min,
				Seconds: sec,
				lpos:    lpos,
				rpos:    s.rpos,
			}
		}
	case false:
		switch index {
		default:
			err = s.error()
		case 3:
			expr = &CalendarLit{
				Kind:  DATE,
				Year:  year,
				Month: month,
				Day:   day,
				lpos:  lpos,
				rpos:  s.rpos,
			}
		case 6:
			expr = &CalendarLit{
				Kind:    DATETIME,
				Year:    year,
				Month:   month,
				Day:     day,
				Hours:   hour,
				Minutes: min,
				Seconds: sec,
				lpos:    lpos,
				rpos:    s.rpos,
			}
		}
	}
	s.next()
	if isTimeUnitPostfix(s.lit) {
		u := unitFromString(s.lit)
		switch typ := expr.(type) {
		case *CalendarLit:
			typ.U = u
			typ.rpos += u.size()
		}
		s.next()
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
		expr = &PercentLit{Val: v, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isDateTimePrefix(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		return s.parseDateTime(plit)
	case isPressureUnit(unit):
		u := unitFromString(s.lit)
		s.rpos += litlen + u.size()
		expr = &PressureLit{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isDistanceUnit(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		u := unitFromString(unit)
		s.rpos += litlen + u.size()
		expr = &DistanceLit{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
		s.next()
	case isSpeedUnit(unit):
		if s.isSignMinus() {
			s.err = errNegativeValue
			return nil, s.error()
		}
		u := unitFromString(s.lit)
		s.rpos += litlen + u.size()
		expr = &SpeedLit{Val: v, U: u, lpos: s.lpos, rpos: s.rpos}
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
		expr = &TemperatureLit{Val: v, U: u, Vec: ts, lpos: s.lpos, rpos: s.rpos}
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
			expr = &DurationLit{Val: dur, lpos: s.lpos, rpos: s.rpos}
			s.next()
		}
	}
	return
}

func (s *parser) parseSelectorProps(selector *SelectorExpr) error {
	for {
		prop, err := s.parseUnaryExpr()
		if err != nil {
			return err
		}
		if selector.Props == nil {
			selector.Props = make([]Expr, 0)
		}
		selector.Props = append(selector.Props, prop)
		if s.except(EOF) {
			break
		}
		if !s.except(COMMA) {
			break
		}
	}
	return nil
}

func (s *parser) parseSelectorExpr() (expr Expr, err error) {
	selector := &SelectorExpr{Ident: s.lit, lpos: s.t.Offset()}

	s.next()

	if !s.except(LBRACE, COLON) {
		selector.Wildcard = true
		selector.calculateEnd(s.t.Offset())
		return selector, nil
	}

	if s.except(COLON) {
		err = s.parseSelectorProps(selector)
		if err != nil {
			return nil, err
		}
		selector.Wildcard = true
		selector.calculateEnd(s.t.Offset())
		return selector, nil
	}

	if s.except(LBRACE) {
		var i int
		for {
			s.next()
			if s.except(EOF, RBRACE) {
				break
			}
			if s.except(COMMA) {
				continue
			}
			if s.except(MUL) {
				selector.Wildcard = true
				continue
			}
			i++
			if !s.except(STRING) {
				return nil, s.error()
			}
			if selector.Args == nil {
				selector.Args = make(map[string]struct{})
			}
			selector.Args[trim(s.t.TokenText())] = struct{}{}
		}
		if i == 0 {
			selector.Wildcard = true
		}
	}

	s.next()

	if s.except(COLON) {
		err = s.parseSelectorProps(selector)
		if err != nil {
			return nil, err
		}
	}
	selector.calculateEnd(s.t.Offset())
	return selector, nil
}

func (s *parser) error() error {
	err := Error{
		Offset: s.t.s.Offset,
		Err:    s.err,
		Lit:    s.t.lit,
	}
	_, er := s.r.Seek(0, io.SeekStart)
	if er == nil {
		buf := make([]byte, s.t.s.Offset)
		_, _ = s.r.Read(buf)
		err.Msg = string(buf)
	}
	return &err
}

func newParser(t *Tokenizer, r *strings.Reader) *parser {
	return &parser{t: t, r: r}
}

type Error struct {
	Offset int
	Err    error
	Msg    string
	Lit    string
}

func (e *Error) Error() string {
	var ctx string
	if e.Err != nil {
		ctx = "error: " + e.Err.Error()
	}
	return fmt.Sprintf("syntax error at position offset=%d, near=%s\n```\n%s ...^\n```\n%s",
		e.Offset, e.Lit, strings.TrimSpace(e.Msg), ctx)
}

func isDateTimePrefix(s string) (ok bool) {
	switch s {
	case "-", ":":
		ok = true
	}
	return
}

func isGeometryToken(tok Token) (ok bool) {
	switch tok {
	case GEOMETRY_POINT, GEOMETRY_MULTIPOINT,
		GEOMETRY_LINE, GEOMETRY_MULTILINE,
		GEOMETRY_POLYGON, GEOMETRY_MULTIPOLYGON:
		ok = true
	}
	return
}

func trim(lit string) string {
	lit = strings.TrimLeft(lit, `"`)
	lit = strings.TrimRight(lit, `"`)
	return lit
}
