package geoqlparser

type Visitor interface {
	Visit(expr Expr) (w Visitor)
}
