package geoqlparser

import (
	"errors"
)

var errIllegalExpr = errors.New("illegal expression")

type ApplyBinaryExprFunc func(left Expr, right Expr, op Token) (bool, error)

func ApplyBinaryExpr(stmt Statement, fn ApplyBinaryExprFunc) (ok bool, err error) {
	switch typ := stmt.(type) {
	case *Trigger:
		if typ.WhenFlat != nil {
			return walkBinaryExprFlat(typ.WhenFlat, fn)
		} else {
			expr, err := walkBinaryExpr(typ.When, fn)
			if err != nil {
				return false, err
			}
			boolVal, ok := expr.(*BooleanExpr)
			if !ok {
				return false, errIllegalExpr
			}
			return boolVal.V, nil
		}
	}
	return
}

func walkBinaryExprFlat(expressions []Expr, fn ApplyBinaryExprFunc) (ok bool, err error) {
	for i := 0; i < len(expressions); i += 2 {
		expr := expressions[i]
		be, bok := expr.(*BinaryExpr)
		if !bok {
			return false, errIllegalExpr
		}
		res, er := fn(be.Left, be.Right, be.Op)
		if er != nil {
			return false, er
		}
		if i > 0 {
			op := expressions[i-1]
			ope, opk := op.(*OpExpr)
			if !opk {
				return false, errIllegalExpr
			}
			switch ope.Op {
			case AND, LAND:
				ok = ok && res
			case OR, LOR:
				ok = ok || res
			}
		} else {
			ok = res
		}
	}
	return
}

func walkBinaryExpr(expr Expr, fn ApplyBinaryExprFunc) (Expr, error) {
	switch node := expr.(type) {
	case *ParenExpr:
		return walkBinaryExpr(node.Expr, fn)
	case *BinaryExpr:
		left, err := walkBinaryExpr(node.Left, fn)
		if err != nil {
			return nil, err
		}
		right, err := walkBinaryExpr(node.Right, fn)
		if err != nil {
			return nil, err
		}
		switch node.Op {
		case AND, LAND:
			lv, lok := left.(*BooleanExpr)
			rv, rok := right.(*BooleanExpr)
			if !lok || !rok {
				return nil, errIllegalExpr
			}
			return &BooleanExpr{V: lv.V && rv.V}, nil
		case OR, LOR:
			lv, lok := left.(*BooleanExpr)
			rv, rok := right.(*BooleanExpr)
			if !lok || !rok {
				return nil, errIllegalExpr
			}
			return &BooleanExpr{V: lv.V || rv.V}, nil
		default:
			ok, err := fn(left, right, node.Op)
			if err != nil {
				return nil, err
			}
			return &BooleanExpr{V: ok}, nil
		}
	}
	return expr, nil
}
