package geoqlparser

import (
	"fmt"
)

func wantTokens(op Token, tokens ...Token) bool {
	for i := 0; i < len(tokens); i++ {
		if op == tokens[i] {
			return true
		}
	}
	return false
}

func makeFlat(t *Trigger, expr Expr) Expr {
	switch node := expr.(type) {
	case *ParenExpr:
		return makeFlat(t, node.Expr)
	case *BinaryExpr:
		left := makeFlat(t, node.Left)
		if wantTokens(node.Op, AND, LAND, OR, LOR) {
			t.WhenFlat = append(t.WhenFlat, &OpExpr{Op: node.Op, Pos: node.Pos})
		}
		right := makeFlat(t, node.Right)
		_, lok := left.(*BinaryExpr)
		_, rok := right.(*BinaryExpr)
		if !lok && !rok && left == node.Left && right == node.Right {
			t.WhenFlat = append(t.WhenFlat, node)
		}
	}
	return expr
}

func ToFlat(stmt Statement) (Statement, error) {
	switch typ := stmt.(type) {
	case *Trigger:
		makeFlat(typ, typ.When)
		if len(typ.WhenFlat)%2 == 0 {
			return nil, fmt.Errorf("ToFlat(%T) => have %d even number of expressions instead of odd ones",
				stmt, len(typ.WhenFlat))
		}
		if len(typ.WhenFlat) > 0 {
			typ.When = nil
		}
		return stmt, nil
	}
	return stmt, nil
}
