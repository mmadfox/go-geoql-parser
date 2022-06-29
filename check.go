package geoqlparser

import (
	"errors"
	"fmt"
)

type SelectorType int

const (
	Int SelectorType = iota + 1
	Float
	String
	DateTime
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

func Check(stmt Statement, dict Dictionary) error {
	switch typ := stmt.(type) {
	default:
		return fmt.Errorf("walk: unknown statement %T", stmt)
	case *Trigger:
		checker := &triggerChecker{ctx: typ, dict: dict}
		return checker.check()
	}
}

type triggerChecker struct {
	ctx  *Trigger
	dict Dictionary
}

func (tc *triggerChecker) check() error {
	_, err := tc.walk(tc.ctx.When)
	return err
}

func (tc *triggerChecker) eval(left, right Expr, op Token) (expr Expr, err error) {
	switch op {
	default:
		err = fmt.Errorf("todo eval error")
	case OR, AND:
		switch {
		default:
			err = fmt.Errorf("todo: eval boolean")
		case tc.isBoolean(left):
			if tc.isBoolean(right) {
				expr = opBoolean
			} else {
				err = fmt.Errorf("todo: eval boolean")
			}
		}
	case EQL, LEQL, NOT_EQ, LNEQ:
		switch {
		default:
			err = fmt.Errorf("todo eval EQL, LEQL, NOT_EQ, LNEQ")
		case tc.isBoolean(left):
			switch {
			default:
				err = fmt.Errorf("todo eval")
			case tc.isBoolean(right):
				expr = opBoolean
			}
		case tc.isString(left):
			switch {
			default:
				err = fmt.Errorf("todo eval")
			case tc.isString(right):
				expr = opBoolean
			}
		case tc.isInt(left):
			switch {
			default:
				err = fmt.Errorf("todo eval int")
			case tc.isInt(right):
				expr = opBoolean
			case tc.isFloat(right):
				expr = opBoolean
			}
		}
	case GEQ, LEQ, GTR, LSS:
		switch {
		default:
			err = fmt.Errorf("todo: eval GEQ, LEQ, GTR, LSS")
		case tc.isInt(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval GEQ, LEQ, GTR, LSS")
			case tc.isInt(right):
				expr = opBoolean
			case tc.isFloat(right):
				expr = opBoolean
			}
		case tc.isFloat(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval GEQ, LEQ, GTR, LSS")
			case tc.isInt(right):
				expr = opBoolean
			case tc.isFloat(right):
				expr = opBoolean
			}
		}

	case QUO, MUL, SUB, ADD, REM:
		switch {
		default:
			err = fmt.Errorf("todo: eval QUO, MUL, SUB, ADD, REM")
		case tc.isString(left):
			if op == ADD && tc.isString(right) {
				expr = opString
			} else {
				err = fmt.Errorf("todo: eval string")
			}
		case tc.isInt(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval int")
			case tc.isInt(right):
				expr = opInt
			case tc.isFloat(right):
				expr = opFloat
			}
		case tc.isFloat(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval int")
			case tc.isInt(right):
				expr = opFloat
			case tc.isFloat(right):
				expr = opFloat
			}
		}
	case IN, NOT_IN:
		switch {
		default:
			err = fmt.Errorf("todo: IN, NOT_IN")
		case tc.isString(left):
			if tc.isArray(right, STRING) {
				expr = opArrayString
			} else {
				err = fmt.Errorf("todo: eval int")
			}
		case tc.isFloat(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval int")
			case tc.isArray(right, FLOAT):
				expr = opArrayFloat
			case tc.isArray(right, INT):
				expr = opArrayFloat
			case tc.isRange(right, FLOAT):
				expr = opRangeFloat
			case tc.isRange(right, INT):
				expr = opRangeFloat
			}
		case tc.isInt(left):
			switch {
			default:
				err = fmt.Errorf("todo: eval int")
			case tc.isArray(right, FLOAT):
				expr = opArrayFloat
			case tc.isArray(right, INT):
				expr = opArrayInt
			case tc.isRange(right, FLOAT):
				expr = opRangeFloat
			case tc.isRange(right, INT):
				expr = opRangeInt
			}
		}
	}
	return
}

func (tc *triggerChecker) walk(expr Expr) (Expr, error) {
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

func (tc *triggerChecker) isBoolean(expr Expr) (ok bool) {
	switch typ := expr.(type) {
	case *BooleanTyp:
		ok = true
	case *Ref:
		assign, err := tc.ctx.findVar(typ.ID)
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

func (tc *triggerChecker) isString(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *StringTyp:
		ok = true
	case *Ref:
		assign, err := tc.ctx.findVar(typ.ID)
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
		switch selector.(type) {
		case *StringTyp:
			ok = true
		}
	}
	return
}

func (tc *triggerChecker) isRange(in Expr, typ Token) (ok bool) {
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

func (tc *triggerChecker) isArray(in Expr, typ Token) (ok bool) {
	array, ok := in.(*ArrayTyp)
	if !ok {
		return
	}
	if len(array.List) == 0 {
		return
	}
	item := array.List[0]
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

func (tc *triggerChecker) isInt(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *IntTyp:
		ok = true
	case *Range:
		_, ok = typ.Low.(*IntTyp)
	case *Ref:
		assign, err := tc.ctx.findVar(typ.ID)
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

func (tc *triggerChecker) isFloat(in Expr) (ok bool) {
	switch typ := in.(type) {
	case *FloatTyp, *PercentTyp, *PressureTyp, *DistanceTyp, *SpeedTyp, *TemperatureTyp:
		ok = true
	case *Range:
		switch typ.Low.(type) {
		case *FloatTyp, *PercentTyp, *PressureTyp, *DistanceTyp, *SpeedTyp, *TemperatureTyp:
			ok = true
		}
	case *Ref:
		assign, err := tc.ctx.findVar(typ.ID)
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
		switch selector.(type) {
		case *FloatTyp, *PercentTyp, *PressureTyp, *DistanceTyp, *SpeedTyp, *TemperatureTyp:
			ok = true
		}
	}
	return
}

func (tc *triggerChecker) getSelectorType(selectorName string) (expr Expr, err error) {
	typ, err := tc.dict.lookup(selectorName)
	if err != nil {
		return nil, err
	}
	switch typ {
	default:
		err = errors.New("todo")
	case Int:
		expr = opInt
	case Float:
		expr = opFloat
	case String:
		expr = opString
	case ArrayInt:
		expr = opArrayInt
	case ArrayFloat:
		expr = opArrayFloat
	case ArrayString:
		expr = opArrayString
	}
	return
}

type CheckError struct {
	Op    Token
	Left  Expr
	Eight Expr
	Msg   string
}

func (e *CheckError) Error() string {
	return ""
}
