package geoqlparser

import (
	"bytes"
	"errors"
	"fmt"
)

type SelectorType int

const (
	Int SelectorType = iota + 1
	Float
	String
	Boolean
	ArrayInt
	ArrayFloat
	ArrayString
)

type Dictionary map[string]SelectorType

func Dict() Dictionary { return make(Dictionary) }

func (d Dictionary) lookup(selectorName string) (SelectorType, error) {
	typ, ok := d[selectorName]
	if !ok {
		return -1, fmt.Errorf("selector type %s not declared", selectorName)
	}
	return typ, nil
}

var (
	opFloat       = &FloatTyp{}
	opInt         = &IntTyp{}
	opString      = &StringTyp{}
	opBoolean     = &BooleanTyp{}
	opArrayFloat  = &ArrayTyp{List: []Expr{opFloat}}
	opArrayInt    = &ArrayTyp{List: []Expr{opInt}}
	opArrayString = &ArrayTyp{List: []Expr{opString}}
	opRangeInt    = &Range{Low: opInt}
	opRangeFloat  = &Range{Low: opFloat}
)

type ruleTyp uint

const (
	isUnknown ruleTyp = iota
	isInt
	isFloat
	isString
	isRangeInt
	isRangeFloat
	isArrayInt
	isArrayFloat
	isArrayString
	isGeometry
	isBoolean
)

var rules = map[Token]map[ruleTyp][]ruleTyp{
	OR: {
		isBoolean: {isBoolean},
	},
	AND: {
		isBoolean: {isBoolean},
	},
	EQL: {
		isInt:         {isInt, isFloat},
		isFloat:       {isInt, isFloat},
		isString:      {isString},
		isBoolean:     {isBoolean},
		isArrayString: {isArrayString},
		isArrayInt:    {isArrayInt, isArrayFloat, isGeometry},
		isArrayFloat:  {isArrayFloat, isArrayInt, isGeometry},
		isGeometry:    {isArrayInt, isArrayFloat, isGeometry},
	},
	LEQL: {
		isInt:         {isInt, isFloat},
		isFloat:       {isInt, isFloat},
		isString:      {isString},
		isBoolean:     {isBoolean},
		isArrayString: {isArrayString},
		isArrayInt:    {isArrayInt, isArrayFloat, isGeometry},
		isArrayFloat:  {isArrayFloat, isArrayInt, isGeometry},
		isGeometry:    {isArrayInt, isArrayFloat, isGeometry},
	},
	NOT_EQ: {
		isInt:         {isInt, isFloat},
		isFloat:       {isInt, isFloat},
		isString:      {isString},
		isBoolean:     {isBoolean},
		isArrayString: {isArrayString},
		isArrayInt:    {isArrayInt, isArrayFloat, isGeometry},
		isArrayFloat:  {isArrayFloat, isArrayInt, isGeometry},
		isGeometry:    {isArrayInt, isArrayFloat, isGeometry},
	},
	LNEQ: {
		isInt:         {isInt, isFloat},
		isFloat:       {isInt, isFloat},
		isString:      {isString},
		isBoolean:     {isBoolean},
		isArrayString: {isArrayString},
		isArrayInt:    {isArrayInt, isArrayFloat, isGeometry},
		isArrayFloat:  {isArrayFloat, isArrayInt, isGeometry},
		isGeometry:    {isArrayInt, isArrayFloat, isGeometry},
	},
	GEQ: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	LEQ: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	GTR: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	LSS: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	QUO: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	MUL: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	SUB: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	ADD: {
		isInt:    {isInt, isFloat},
		isFloat:  {isInt, isFloat},
		isString: {isString},
	},
	REM: {
		isInt:   {isInt, isFloat},
		isFloat: {isInt, isFloat},
	},
	IN: {
		isInt:      {isArrayInt, isRangeInt, isRangeFloat, isGeometry},
		isFloat:    {isArrayFloat, isRangeFloat, isRangeInt, isGeometry},
		isString:   {isArrayString},
		isGeometry: {isGeometry},
	},
	NOT_IN: {
		isInt:      {isArrayInt, isRangeInt, isRangeFloat, isGeometry},
		isFloat:    {isArrayFloat, isRangeFloat, isRangeInt, isGeometry},
		isString:   {isArrayString},
		isGeometry: {isGeometry},
	},
	NEARBY: {
		isGeometry: {isGeometry},
		isFloat:    {isGeometry},
		isInt:      {isGeometry},
	},
	NOT_NEARBY: {
		isGeometry: {isGeometry},
		isFloat:    {isGeometry},
		isInt:      {isGeometry},
	},
	INTERSECTS: {
		isGeometry: {isGeometry},
		isFloat:    {isGeometry},
		isInt:      {isGeometry},
	},
	NOT_INTERSECTS: {
		isGeometry: {isGeometry},
		isFloat:    {isGeometry},
		isInt:      {isGeometry},
	},
}

func CheckType(stmt Statement, dict Dictionary) (err error) {
	switch typ := stmt.(type) {
	case *Trigger:
		ch := &checker{trigger: typ, dict: dict}
		err = ch.check()
	}
	return
}

type checker struct {
	trigger *Trigger
	dict    Dictionary
}

func (tc *checker) check() error {
	_, err := tc.walk(tc.trigger.When)
	return err
}

func (tc *checker) error(left, right Expr, op Token, msg string, mism string) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(msg)
	buf.WriteString(": ")
	left.format(buf, "", true)
	buf.WriteString(" ")
	buf.WriteString(KeywordString(op))
	buf.WriteString(" ")
	right.format(buf, "", true)
	if len(mism) > 0 {
		buf.WriteString(" ")
		buf.WriteString(mism)
	}
	return errors.New(buf.String())
}

func (tc *checker) eval(left, right Expr, op Token) (expr Expr, err error) {
	rule, ok := rules[op]
	if !ok {
		err = tc.error(left, right, op, "invalid operator", "")
		return
	}
	var checkOk bool
	var l, r ruleTyp
loop:
	for leftTyp, rightTypes := range rule {
		if lok := tc.is(leftTyp, left); !lok {
			continue
		}
		for i := 0; i < len(rightTypes); i++ {
			rt := rightTypes[i]
			if rok := tc.is(rt, right); rok {
				checkOk = true
				l = leftTyp
				r = rightTypes[i]
				break loop
			}
		}
	}
	if checkOk {
		if expr = tc.toExpr(l, r, op); expr == nil {
			err = tc.error(left, right, op, "can't find internal expression", "")
		}
	} else {
		err = tc.error(left, right, op, "invalid operator", "(mismatched types)")
	}
	return
}

func (tc *checker) walk(expr Expr) (Expr, error) {
	switch typ := expr.(type) {
	case *ParenExpr:
		return tc.walk(typ.Expr)
	case *BinaryExpr:
		var left, right Expr
		var err error
		left, err = tc.walk(typ.Left)
		if err != nil {
			return nil, err
		}
		right, err = tc.walk(typ.Right)
		if err != nil {
			return nil, err
		}
		return tc.eval(left, right, typ.Op)
	}
	return expr, nil
}

func (tc *checker) is(rt ruleTyp, in Expr) (ok bool) {
	switch rt {
	case isGeometry:
		ok = tc.isGeometry(in)
	case isBoolean:
		ok = tc.isBoolean(in)
	case isInt:
		ok = tc.isInt(in)
	case isFloat:
		ok = tc.isFloat(in)
	case isString:
		ok = tc.isString(in)
	case isArrayInt:
		ok = tc.isArray(in, INT)
	case isArrayFloat:
		ok = tc.isArray(in, FLOAT)
	case isArrayString:
		ok = tc.isArray(in, STRING)
	case isRangeInt:
		ok = tc.isRange(in, INT)
	case isRangeFloat:
		ok = tc.isRange(in, FLOAT)
	}
	return
}

func (tc *checker) toExpr(left, right ruleTyp, op Token) (expr Expr) {
	switch op {
	case ADD:
		switch {
		case left == isFloat && right == isInt:
			expr = opFloat
		case left == isFloat && right == isFloat:
			expr = opFloat
		case left == isInt && right == isFloat:
			expr = opFloat
		case left == isInt && right == isInt:
			expr = opInt
		case left == isString && right == isString:
			expr = opString
		}
	case QUO, MUL, SUB, REM:
		switch {
		case left == isFloat && right == isInt:
			expr = opFloat
		case left == isFloat && right == isFloat:
			expr = opFloat
		case left == isInt && right == isFloat:
			expr = opFloat
		case left == isInt && right == isInt:
			expr = opInt
		}
	case AND, OR, EQL, LEQL, NOT_EQ, LNEQ, GEQ, LEQ, GTR, LSS:
		expr = opBoolean
	case IN, NOT_IN:
		expr = opBoolean
	case NEARBY, NOT_NEARBY:
		expr = opBoolean
	case INTERSECTS, NOT_INTERSECTS:
		expr = opBoolean
	}
	return
}

func (tc *checker) isBoolean(expr Expr) (ok bool) {
	switch typ := expr.(type) {
	case *BooleanTyp:
		ok = true
	case *Ref:
		assign, err := tc.trigger.findAssign(typ.ID)
		if err != nil {
			return
		}
		switch assign.Right.(type) {
		case *BooleanTyp:
			ok = true
		}

	case *Selector:
		selector, err := tc.getSelectorType(typ.Ident)
		if err != nil {
			return
		}
		switch selector.(type) {
		case *BooleanTyp:
			ok = true
		}
	}
	return
}

func (tc *checker) isString(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *StringTyp:
		ok = true
	case *Ref:
		assign, err := tc.trigger.findAssign(typ.ID)
		if err != nil {
			return
		}
		switch assign.Right.(type) {
		case *StringTyp:
			ok = true
		}

	case *Selector:
		selector, err := tc.getSelectorType(typ.Ident)
		if err != nil {
			return
		}
		switch t := selector.(type) {
		case *ArrayTyp:
			if len(t.List) > 0 {
				_, ok = t.List[0].(*StringTyp)
			}
		case *StringTyp:
			ok = true
		}
	}
	return
}

func (tc *checker) isRange(in Expr, typ Token) (ok bool) {
	range_, ok := in.(*Range)
	if !ok {
		return
	}
	switch typ {
	case INT:
		return tc.isInt(range_.Low)
	case FLOAT:
		return tc.isFloat(range_.Low)
	case STRING:
		return tc.isString(range_.Low)
	}
	return
}

func (tc *checker) isArray(in Expr, typ Token) (ok bool) {
	var item Expr
	switch t := in.(type) {
	default:
		return
	case *Selector:
		item = t
	case *ArrayTyp:
		if len(t.List) == 0 {
			return
		}
		item = t.List[0]
	case *Ref:
		assign, err := tc.trigger.findAssign(t.ID)
		if err != nil {
			return
		}
		switch tt := assign.Right.(type) {
		default:
			return
		case *ArrayTyp:
			if len(tt.List) == 0 {
				return
			}
			item = tt.List[0]
		}
	}

	switch typ {
	case BOOLEAN:
		return tc.isBoolean(item)
	case INT:
		return tc.isInt(item)
	case FLOAT:
		return tc.isFloat(item)
	case STRING:
		return tc.isString(item)
	}
	return
}

func (tc *checker) isGeometry(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *GeometryCollectionTyp, *GeometryLineTyp, *GeometryPointTyp,
		*GeometryMultiObjectTyp:
		ok = true
	case *Ref:
		assign, err := tc.trigger.findAssign(typ.ID)
		if err != nil {
			return
		}
		switch assign.Right.(type) {
		case *GeometryCollectionTyp, *GeometryLineTyp, *GeometryPointTyp,
			*GeometryMultiObjectTyp:
			ok = true
		}
	}
	return
}

func (tc *checker) isInt(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *IntTyp:
		ok = true
	case *Ref:
		assign, err := tc.trigger.findAssign(typ.ID)
		if err != nil {
			return
		}
		switch assign.Right.(type) {
		case *IntTyp:
			ok = true
		}

	case *Selector:
		selector, err := tc.getSelectorType(typ.Ident)
		if err != nil {
			return
		}
		switch selector.(type) {
		case *IntTyp:
			ok = true
		}
	}
	return
}

func (tc *checker) isNumber(in Expr) (ok bool) {
	return tc.isInt(in) || tc.isFloat(in)
}

func (tc *checker) isFloat(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *FloatTyp, *PercentTyp, *PressureTyp, *DistanceTyp, *SpeedTyp,
		*TemperatureTyp:
		ok = true
	case *Ref:
		assign, err := tc.trigger.findAssign(typ.ID)
		if err != nil {
			return
		}
		switch assign.Right.(type) {
		case *FloatTyp, *PercentTyp, *PressureTyp, *DistanceTyp, *SpeedTyp, *TemperatureTyp:
			ok = true
		}

	case *Selector:
		selector, err := tc.getSelectorType(typ.Ident)
		if err != nil {
			return
		}
		if array, isArray := selector.(*ArrayTyp); isArray {
			selector = array.List[0]
		}
		switch selector.(type) {
		case *FloatTyp:
			ok = true
		}
	}
	return
}

func (tc *checker) getSelectorType(selectorName string) (expr Expr, err error) {
	typ, err := tc.dict.lookup(selectorName)
	if err != nil {
		return nil, err
	}
	switch typ {
	default:
		err = fmt.Errorf("cannot find selector type '%s'", selectorName)
	case Int:
		expr = opInt
	case Float:
		expr = opFloat
	case String:
		expr = opString
	case Boolean:
		expr = opBoolean
	case ArrayInt:
		expr = opArrayInt
	case ArrayFloat:
		expr = opArrayFloat
	case ArrayString:
		expr = opArrayString
	}
	return
}
