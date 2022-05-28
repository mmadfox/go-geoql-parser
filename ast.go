package geoqlparser

import "time"

type Statement interface {
	isStatement()
}

func (t *Trigger) isStatement() {}

// Trigger represents a TRIGGER statement.
type Trigger struct {
	Vars     map[string]interface{}
	When     Expr
	WhenFlat []Expr
	Repeat   *RepeatExpr
	Reset    *ResetExpr
	Pos      Pos
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
	Pos   Pos
}

type RepeatExpr struct {
	V        int
	Interval time.Duration
	Pos      Pos
}

type ResetExpr struct {
	Dur DurVal
	Pos Pos
}

type ParenExpr struct {
	Expr     Expr
	LeftPos  Pos
	RightPos Pos
}

type WildcardLit struct {
	Pos Pos
}

type BasicLit struct {
	V    interface{}
	Kind Token
	Pos  Pos
}

type BaseSelectorExpr struct {
	Ident     Token
	Args      map[string]struct{}
	Vars      map[string]struct{}
	Qualifier Qualifier
	Wildcard  bool
	LeftPos   Pos
	RightPos  Pos
}

type TrackerSelectorExpr struct {
	Ident    Token
	Args     map[string]struct{}
	Vars     map[string]struct{}
	Wildcard bool
	Radius   RadiusVal
	LeftPos  Pos
	RightPos Pos
}

type BooleanExpr struct {
	V bool
}

type OpExpr struct {
	Op  Token
	Pos Pos
}

type ExprList []Expr

func (n *BinaryExpr) isExpr()          {}
func (n *ParenExpr) isExpr()           {}
func (n *BaseSelectorExpr) isExpr()    {}
func (n *TrackerSelectorExpr) isExpr() {}
func (n *WildcardLit) isExpr()         {}
func (n *BasicLit) isExpr()            {}
func (n *RepeatExpr) isExpr()          {}
func (n *ResetExpr) isExpr()           {}
func (n *BooleanExpr) isExpr()         {}
func (e *OpExpr) isExpr()              {}
