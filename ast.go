package geoqlparser

import "time"

type Statement interface {
	isStatement()
}

func (t *Trigger) isStatement() {}

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

type WildcardLit struct {
	Pos Pos
}

type BaseSelectorLit struct {
	Ident     Token
	Args      map[string]struct{}
	Vars      map[string]struct{}
	Qualifier Qualifier
	Wildcard  bool
	Pos       Pos
}

type TrackerSelectorLit struct {
	Ident    Token
	Args     map[string]struct{}
	Vars     map[string]struct{}
	Wildcard bool
	Radius   RadiusVal
	Pos      Pos
}

type ExprList []Expr

func (n *BinaryExpr) isExpr()         {}
func (n *ParenExpr) isExpr()          {}
func (n *BaseSelectorLit) isExpr()    {}
func (n *TrackerSelectorLit) isExpr() {}
func (n *WildcardLit) isExpr()        {}
