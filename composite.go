package geoqlparser

func (s *parser) parseArrayExpr() (expr Expr, err error) {
	if !s.except(LBRACK) {
		return nil, s.error()
	}

	var arrayExpr *ArrayTyp
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
			arrayExpr = &ArrayTyp{lpos: startPos, List: make([]Expr, 0), Kind: ILLEGAL}
		}

		switch expr.(type) {
		default:
			err = s.error()
		case *Selector:
			err = checkKind(arrayExpr.Kind, SELECTOR)
			arrayExpr.Kind = SELECTOR
		case *DurationTyp:
			err = checkKind(arrayExpr.Kind, DURATION)
			arrayExpr.Kind = DURATION
		case *SpeedTyp:
			err = checkKind(arrayExpr.Kind, SPEED)
			arrayExpr.Kind = SPEED
		case *PressureTyp:
			err = checkKind(arrayExpr.Kind, PRESSURE)
			arrayExpr.Kind = PRESSURE
		case *TemperatureTyp:
			err = checkKind(arrayExpr.Kind, TEMPERATURE)
			arrayExpr.Kind = TEMPERATURE
		case *DistanceTyp:
			err = checkKind(arrayExpr.Kind, DISTANCE)
			arrayExpr.Kind = DISTANCE
		case *PercentTyp:
			err = checkKind(arrayExpr.Kind, PERCENT)
			arrayExpr.Kind = PERCENT
		case *IntTyp:
			err = checkKind(arrayExpr.Kind, INT)
			arrayExpr.Kind = INT
		case *FloatTyp:
			err = checkKind(arrayExpr.Kind, FLOAT)
			arrayExpr.Kind = FLOAT
		case *StringTyp:
			err = checkKind(arrayExpr.Kind, STRING)
			arrayExpr.Kind = STRING
		case *Ref:
			err = checkKind(arrayExpr.Kind, IDENT)
			arrayExpr.Kind = IDENT
		case *Range:
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
	arrayExpr.rpos = s.t.Offset()
	s.next()
	return arrayExpr, nil
}

func (s *parser) parseRangeExpr(low Expr) (expr Expr, err error) {
	_, isRangeExpr := low.(*Range)
	if isRangeExpr {
		return nil, s.error()
	}
	startPos := s.t.Offset()
	high, err := s.parseUnaryExpr()
	if err != nil {
		return
	}
	return &Range{
		Low:  low,
		High: high,
		lpos: startPos,
		rpos: s.t.Offset(),
	}, nil
}

func (s *parser) parseSelectorExpr() (expr Expr, err error) {
	selector := &Selector{Ident: s.lit, lpos: s.t.Offset()}

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

func (s *parser) parseSelectorProps(selector *Selector) error {
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
