package geoqlparser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func Parse(gql string) (Statement, error) {
	r := strings.NewReader(gql)
	t := NewTokenizer(r)
	s := newParser(t, r)
	return s.parse0()
}

type parser struct {
	r   *strings.Reader
	t   *Tokenizer
	pos Pos
	tok Token
	lit string
	err error
	neg bool
}

func (s *parser) next() {
	s.tok, s.lit = s.t.Scan()
	s.pos = s.t.Offset()
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
	stmt.Pos = s.t.Offset()
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
	if !s.except(INT) {
		return s.error()
	}
	intVal, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	iiv, ok := intVal.(*IntLit)
	if !ok {
		return s.error()
	}
	if iiv.Val < 2 {
		return s.error()
	}
	if s.lit != "times" {
		s.err = fmt.Errorf("got %s, expected times", s.lit)
		return s.error()
	}
	s.next()
	dur, err := s.parseIntTypes()
	if err != nil {
		return err
	}
	if s.lit != "interval" {
		s.err = fmt.Errorf("got %s, expected interval", s.lit)
		return s.error()
	}
	_, ok = dur.(*DurationLit)
	if !ok {
		return s.error()
	}
	s.next()
	stmt.RepeatInterval = dur
	stmt.RepeatCount = intVal
	return
}

func (s *parser) parseSet(stmt *TriggerStmt) error {
	var varname string
	s.next()
	for {
		if s.except(EOF) {
			break
		}
		if s.except(WHEN) {
			s.t.Reset()
			break
		}
		varname = s.t.TokenText()
		s.t.Unwind()
		s.next()
		if !s.except(ASSIGN) {
			return s.error()
		}
		expr, err := s.parseUnaryExpr()
		if err != nil {
			return err
		}
		switch typ := expr.(type) {
		case *VarLit:
			return s.error()
		case *ArrayExpr:
			if typ.Kind == IDENT {
				return s.error()
			}
		}
		stmt.initVars()
		_, found := stmt.Set[varname]
		if found {
			s.err = fmt.Errorf("variable %s already exists", varname)
			return s.error()
		}
		stmt.Set[varname] = expr
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
		left = &BinaryExpr{Left: left, Right: right, Op: op, Pos: pos}
	}
}

func (s *parser) parseUnaryExpr() (expr Expr, err error) {
	if s.t.Err() != nil {
		s.err = s.t.Err()
		return nil, s.error()
	}

	s.next()

	if s.except(SUB, ADD) {
		if s.tok == SUB {
			s.neg = true
		}
		s.next()
	}

	switch s.tok {
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
	case GEOMETRY_POINT, GEOMETRY_MULTIPOINT,
		GEOMETRY_LINE, GEOMETRY_MULTILINE,
		GEOMETRY_POLYGON, GEOMETRY_CIRCLE, GEOMETRY_MULTIPOLYGON:
		expr, err = s.parseGeometryExpr()
	case GEOMETRY_COLLECTION:
		expr, err = s.parseGeometryCollectionExpr()
	case BOOLEAN:
		expr, err = s.parseBooleanLit()
	case WEEKDAY, MONTH:
		expr, err = s.parseCalendarLit()
	}

	switch s.tok {
	case RANGE:
		expr, err = s.parseRangeExpr(expr)
	}

	s.neg = false
	return
}

func (s *parser) parseVarExpr() (expr Expr, err error) {
	s.next()
	expr = &VarLit{ID: s.t.TokenText(), Pos: s.t.Offset()}
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
	expr = &CalendarLit{Val: val, Pos: s.t.Offset(), Kind: s.tok}
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
		Low:      low,
		High:     high,
		StartPos: startPos,
		EndPos:   s.t.Offset(),
	}, nil
}

func (s *parser) parseWildcardLit() (expr Expr, err error) {
	return &WildcardLit{Pos: s.t.Offset()}, nil
}

func (s *parser) parseBooleanLit() (expr Expr, err error) {
	switch s.lit {
	default:
		return nil, s.error()
	case "true", "up":
		expr = &BooleanLit{Val: true, Pos: s.t.Offset()}
	case "false", "down":
		expr = &BooleanLit{Val: false, Pos: s.t.Offset()}
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
			arrayExpr = &ArrayExpr{StartPos: startPos, List: make([]Expr, 0), Kind: ILLEGAL}
		}

		switch expr.(type) {
		default:
			err = s.error()
		case *DateTimeLit:
			err = checkKind(arrayExpr.Kind, DATETIME)
			arrayExpr.Kind = DATETIME
		case *TimeLit:
			err = checkKind(arrayExpr.Kind, TIME)
			arrayExpr.Kind = TIME
		case *DateLit:
			err = checkKind(arrayExpr.Kind, DATE)
			arrayExpr.Kind = DATE
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
		case *VarLit:
			err = checkKind(arrayExpr.Kind, IDENT)
			arrayExpr.Kind = IDENT
		case *RangeExpr:
			err = checkKind(arrayExpr.Kind, RANGE)
			arrayExpr.Kind = RANGE
		}
		if err != nil {
			return nil, err
		}

		arrayExpr.List = append(arrayExpr.List, expr)

		if s.except(RBRACK) {
			break
		}
	}
	arrayExpr.EndPos = s.t.Offset()
	s.next()
	return arrayExpr, nil
}

func (s *parser) parseGeometryCollectionExpr() (expr Expr, err error) {
	collection := &GeometryCollectionExpr{
		Objects:  make([]Expr, 0),
		StartPos: s.t.Offset(),
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
		object, oer := s.parseGeometryExpr()
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
	collection.EndPos = s.t.Offset()
	if len(collection.Objects) == 0 {
		return nil, s.error()
	}
	return collection, nil
}

func (s *parser) parseGeometryExpr() (expr Expr, err error) {
	geojsontyp := s.tok
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
			s.neg = true
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

		if s.neg {
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
		s.neg = false
	}
	switch geojsontyp {
	case GEOMETRY_POINT:
		return &GeometryPointExpr{Val: aa, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_MULTIPOINT:
		return &GeometryMultiPointExpr{Val: bb, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_LINE:
		return &GeometryLineExpr{Val: bb, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_MULTILINE:
		return &GeometryMultiLineExpr{Val: cc, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_POLYGON:
		return &GeometryPolygonExpr{Val: bb, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_MULTIPOLYGON:
		return &GeometryMultiPolygonExpr{Val: cc, StartPos: sp, EndPos: s.t.Offset()}, nil
	case GEOMETRY_CIRCLE:
		if !s.except(COLON) {
			return nil, s.error()
		}
		s.next()
		switch s.tok {
		case INT:
			re, err := s.parseIntTypes()
			if err != nil {
				return nil, err
			}
			distLit, ok := re.(*DistanceLit)
			if !ok {
				return nil, s.error()
			}
			return &GeometryCircleExpr{Val: aa, Radius: distLit, StartPos: sp, EndPos: s.t.Offset()}, nil
		case FLOAT:
			re, err := s.parseFloatTypes()
			if err != nil {
				return nil, err
			}
			distLit, ok := re.(*DistanceLit)
			if !ok {
				return nil, s.error()
			}
			return &GeometryCircleExpr{Val: aa, Radius: distLit, StartPos: sp, EndPos: s.t.Offset()}, nil
		}
	}
	err = s.error()
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
	return &ParenExpr{Expr: expr, StartPos: lp, EndPos: rp}, nil
}

func (s *parser) parseStringLit() (expr Expr, err error) {
	expr = &StringLit{Val: trim(s.t.TokenText()), Pos: s.t.Offset()}
	s.next()
	return
}

func (s *parser) parseFloatTypes() (expr Expr, err error) {
	val, err := strconv.ParseFloat(s.lit, 64)
	if err != nil {
		return nil, s.error()
	}
	expr, err = s.parseAllTypes(val)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		if s.neg {
			val = -val
		}
		expr = &FloatLit{Val: val, Pos: s.t.Offset()}
	}
	return
}

func (s *parser) parseIntTypes() (expr Expr, err error) {
	intval, err := strconv.Atoi(s.lit)
	if err != nil {
		return nil, s.error()
	}
	expr, err = s.parseAllTypes(float64(intval))
	if err != nil {
		return nil, err
	}
	if expr == nil {
		if s.neg {
			intval = -intval
		}
		expr = &IntLit{Val: intval, Pos: s.t.Offset()}
	}

	return
}

func (s *parser) parseDateTime(prefix string) (expr Expr, err error) {
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
			expr = &TimeLit{
				Hour:    hour,
				Minute:  min,
				Seconds: sec,
				Pos:     s.t.Offset(),
			}
		}
	case false:
		switch index {
		default:
			err = s.error()
		case 3:
			expr = &DateLit{
				Year:  year,
				Month: month,
				Day:   day,
				Pos:   s.t.Offset(),
			}
		case 6:
			expr = &DateTimeLit{
				Year:    year,
				Month:   month,
				Day:     day,
				Hours:   hour,
				Minutes: min,
				Seconds: sec,
				Pos:     s.t.Offset(),
			}
		}
	}
	s.next()
	if isTimeUnitPostfix(s.lit) {
		u := unitFromString(s.lit)
		switch typ := expr.(type) {
		case *TimeLit:
			typ.U = u
		case *DateTimeLit:
			typ.U = u
		}
		s.next()
	}
	return
}

func (s *parser) parseAllTypes(v float64) (expr Expr, err error) {
	plit := s.lit
	s.next()
	unit := s.t.TokenText()
	switch {
	case isPercentUnit(unit):
		expr = &PercentLit{Val: v, Pos: s.t.Offset()}
		s.next()
	case isDateTimePrefix(unit):
		return s.parseDateTime(plit)
	case isPressureUnit(unit):
		u := unitFromString(s.lit)
		expr = &PressureLit{Val: v, U: u, Pos: s.t.Offset()}
		s.next()
	case isDistanceUnit(unit):
		u := unitFromString(s.lit)
		expr = &DistanceLit{Val: v, U: u, Pos: s.t.Offset()}
		s.next()
	case isSpeedUnit(unit):
		u := unitFromString(s.lit)
		expr = &SpeedLit{Val: v, U: u, Pos: s.t.Offset()}
		s.next()
	case isTemperatureUnit(unit):
		ts := Positive
		if s.neg {
			ts = Negative
		}
		u := unitFromString(s.lit)
		expr = &TemperatureLit{Val: v, U: u, Vec: ts, Pos: s.t.Offset()}
		s.next()
	default:
		var fr rune
		if len(s.lit) > 0 {
			fr = rune(s.lit[0])
		}
		switch fr {
		case 'h', 'm', 's':
			dur, er := time.ParseDuration(plit + s.lit)
			if er != nil {
				return nil, s.error()
			}
			expr = &DurationLit{Val: dur, Pos: s.t.Offset()}
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
		if s.except(EOF) {
			break
		}
		if selector.Props == nil {
			selector.Props = make([]Expr, 0)
		}
		selector.Props = append(selector.Props, prop)
		if !s.except(COMMA) {
			break
		}
	}
	return nil
}

func (s *parser) parseSelectorExpr() (expr Expr, err error) {
	selExpr := &SelectorExpr{
		Ident:    s.lit,
		StartPos: s.t.Offset(),
	}
	s.next()
	if !s.except(LBRACE, COLON) {
		selExpr.Wildcard = true
		selExpr.EndPos = s.t.Offset()
		return selExpr, nil
	}

	if s.except(COLON) {
		err = s.parseSelectorProps(selExpr)
		if err != nil {
			return nil, err
		}
		selExpr.EndPos = s.t.Offset()
		return selExpr, nil
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
				selExpr.Wildcard = true
				continue
			}
			i++
			if !s.except(STRING) {
				return nil, s.error()
			}
			if selExpr.Args == nil {
				selExpr.Args = make(map[string]struct{})
			}
			selExpr.Args[trim(s.t.TokenText())] = struct{}{}
		}
		if i == 0 {
			selExpr.Wildcard = true
		}
	}

	s.next()

	if s.except(COLON) {
		err = s.parseSelectorProps(selExpr)
		if err != nil {
			return nil, err
		}
		selExpr.EndPos = s.t.Offset()
		return selExpr, nil
	}

	selExpr.EndPos = s.t.Offset()
	return selExpr, nil
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

func Format(stmt Statement, b *bytes.Buffer) {
	switch typ := stmt.(type) {
	case *TriggerStmt:
		formatTriggerStmt(typ, b)
	}
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
		GEOMETRY_POLYGON, GEOMETRY_MULTIPOLYGON,
		GEOMETRY_CIRCLE:
		ok = true
	}
	return
}

func trim(lit string) string {
	lit = strings.TrimLeft(lit, `"`)
	lit = strings.TrimRight(lit, `"`)
	return lit
}
