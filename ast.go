package geoqlparser

type Statement interface {
	isStatement()
}

// Trigger represents a TRIGGER statement.
type Trigger struct {
	Vars   map[string]interface{}
	When   int
	Repeat int
	Reset  int
}

func (t *Trigger) initVars() {
	if t.Vars != nil {
		return
	}
	t.Vars = make(map[string]interface{})
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

func (t *Trigger) isStatement() {}
