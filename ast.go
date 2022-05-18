package geoqlparser

type Statement interface {
	isStatement()
}

// Trigger represents a TRIGGER statement.
type Trigger struct {
	Vars   map[string]*Variable
	When   int
	Repeat int
	Reset  int
}

type Expr interface {
	isExpr()
}

type BinaryExpr struct {
	Op    Token
	Left  Expr
	Right Expr
}

type ParenExpr struct {
	Expr Expr
}

type ExprList []Expr

type ValType int

const (
	StrVal = ValType(iota)
	IntVal
	FloatVal
)

type Val struct {
	Type ValType
	Data []byte
}

type Variable struct {
	Ident string
	Value *Val
}
