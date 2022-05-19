package geoqlparser

import "time"

type Statement interface {
	isStatement()
}

// Trigger represents a TRIGGER statement.
type Trigger struct {
	Vars   map[string]interface{}
	When   Expr
	Repeat Repeat
	Reset  DurVal
}

func (t *Trigger) initVars() {
	if t.Vars != nil {
		return
	}
	t.Vars = make(map[string]interface{})
}

type Repeat struct {
	V        int
	Interval time.Duration
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
