package geoqlparser

type Visitor interface {
	Visit(expr Expr) Visitor
}

func Walk(v Visitor, expr Expr) {
	if v.Visit(expr) == nil {
		return
	}
	switch typ := expr.(type) {
	case *SelectorExpr:
		if typ.Props != nil {
			for i := 0; i < len(typ.Props); i++ {
				Walk(v, typ.Props[i])
			}
		}
	case *RangeExpr:
		Walk(v, typ.Low)
		Walk(v, typ.High)
	case *ParenExpr:
		Walk(v, typ.Expr)
	case *BinaryExpr:
		Walk(v, typ.Left)
		Walk(v, typ.Right)
	case *ArrayExpr:
		for i := 0; i < len(typ.List); i++ {
			Walk(v, typ.List[i])
		}
	case *TriggerStmt:
		if typ.Vars != nil {
			for _, expr := range typ.Vars {
				Walk(v, expr)
			}
		}
		Walk(v, typ.When)
		Walk(v, typ.RepeatCount)
		Walk(v, typ.RepeatInterval)
		Walk(v, typ.ResetAfter)
	}
	v.Visit(nil)
}

type visitor func(expr Expr) bool

func (f visitor) Visit(expr Expr) Visitor {
	if f(expr) {
		return f
	}
	return nil
}

func Visit(expr Expr, f func(Expr) bool) {
	Walk(visitor(f), expr)
}
