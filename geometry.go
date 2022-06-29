package geoqlparser

import "strconv"

func (s *parser) parseGeometryMultiObject() (expr Expr, err error) {
	geotyp := s.tok
	s.next()
	if !s.except(LBRACK) {
		return nil, s.error()
	}
	s.next()
	multiobj := &GeometryMultiObjectTyp{
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
		case *GeometryPointTyp:
			if multiobj.Kind != GEOMETRY_MULTIPOINT {
				return nil, s.error()
			}
			multiobj.Val = append(multiobj.Val, typ)
		case *GeometryLineTyp:
			if multiobj.Kind != GEOMETRY_MULTILINE {
				return nil, s.error()
			}
			multiobj.Val = append(multiobj.Val, typ)
		case *GeometryPolygonTyp:
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
	collection := &GeometryCollectionTyp{
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
		point := &GeometryPointTyp{Val: aa, lpos: sp, rpos: s.t.Offset() - 1}
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
		line := &GeometryLineTyp{Val: bb, lpos: sp, rpos: s.t.Offset() - 1}
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
		return &GeometryPolygonTyp{Val: cc, lpos: sp, rpos: s.t.Offset() - 1}, nil
	}
	err = s.error()
	return
}
