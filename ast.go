package geoqlparser

import (
	"fmt"
	"io"
	"time"
)

type Statement interface {
	isStatement()
}

type Expr interface {
	Pos() Pos
	End() Pos

	format(b io.StringWriter, padding string, inline bool)
	isExpr()
}

func (t *Trigger) isStatement() {}

// Trigger represents a TRIGGER statement.
type Trigger struct {
	Vars           []*Assign
	When           Expr
	RepeatCount    Expr
	RepeatInterval Expr
	ResetAfter     Expr
	lpos           Pos
	rpos           Pos
}

func (t *Trigger) SetVar(v *Assign) error {
	if t.isAssigned(v.Left.Val) {
		return fmt.Errorf("variable %s already assigned", v.Left.Val)
	}
	t.Vars = append(t.Vars, v)
	return nil
}

func (t *Trigger) isAssigned(varname string) bool {
	for i := 0; i < len(t.Vars); i++ {
		if t.Vars[i].Left.Val == varname {
			return true
		}
	}
	return false
}

func (t *Trigger) initVars() {
	if t.Vars != nil {
		return
	}
	t.Vars = make([]*Assign, 0)
}

func (t *Trigger) findVar(varname string) (*Assign, error) {
	// TODO: if the vars has more than 16 elements, then binary search
	for i := 0; i < len(t.Vars); i++ {
		if t.Vars[i].Left.Val == varname {
			return t.Vars[i], nil
		}
	}
	return nil, fmt.Errorf("variable %s not found", varname)
}

type Assign struct {
	Left   *Ident
	Right  Expr
	TokPos Pos
}

type ArrayTyp struct {
	Kind Token
	List []Expr
	lpos Pos
	rpos Pos
}

type BinaryExpr struct {
	Op    Token
	Left  Expr
	Right Expr
	OpPos Pos
}

type ParenExpr struct {
	Expr Expr
	lpos Pos
	rpos Pos
}

type WildcardTyp struct {
	lpos Pos
}

type TimeTyp struct {
	Hours, Minutes, Seconds int
	U                       Unit
	lpos                    Pos
	rpos                    Pos
}

type DateTyp struct {
	Year, Day, Month int
	lpos             Pos
	rpos             Pos
}

type WeekdayTyp struct {
	Val  int
	lpos Pos
	rpos Pos
}

type MonthTyp struct {
	Val  int
	lpos Pos
	rpos Pos
}

type GeometryPointTyp struct {
	Val    [2]float64
	Radius *DistanceTyp
	lpos   Pos
	rpos   Pos
}

type GeometryMultiObjectTyp struct {
	Kind Token
	Val  []Expr
	lpos Pos
	rpos Pos
}

type GeometryLineTyp struct {
	Val    [][2]float64
	Margin *DistanceTyp
	lpos   Pos
	rpos   Pos
}

func (e *GeometryLineTyp) needExpand() bool {
	return len(e.Val) > 4
}

type GeometryPolygonTyp struct {
	Val  [][][2]float64
	lpos Pos
	rpos Pos
}

func (e *GeometryPolygonTyp) needExpand() (ok bool) {
	for i := 0; i < len(e.Val); i++ {
		if len(e.Val[i]) > 4 {
			ok = true
			break
		}
	}
	return
}

func (e *GeometryPolygonTyp) HasHoles() bool {
	return len(e.Val) > 1
}

type GeometryCollectionTyp struct {
	Objects []Expr
	lpos    Pos
	rpos    Pos
}

type IntTyp struct {
	Val  int
	lpos Pos
	rpos Pos
}

type Range struct {
	Low  Expr
	High Expr
	lpos Pos
	rpos Pos
}

type PercentTyp struct {
	Val  float64
	lpos Pos
	rpos Pos
}

type StringTyp struct {
	Val  string
	lpos Pos
	rpos Pos
}

type FloatTyp struct {
	Val  float64
	lpos Pos
	rpos Pos
}

type DurationTyp struct {
	Val  time.Duration
	lpos Pos
	rpos Pos
}

type TemperatureTyp struct {
	Val  float64
	U    Unit
	Vec  Sign
	lpos Pos
	rpos Pos
}

type PressureTyp struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type DistanceTyp struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type SpeedTyp struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type Ident struct {
	Val  string
	lpos Pos
	rpos Pos
}

type Ref struct {
	ID   string
	lpos Pos
	rpos Pos
}

type Selector struct {
	Ident    string              // selector name
	Args     map[string]struct{} // device ids
	Wildcard bool                // indicates the current device
	Props    []Expr              // some props
	lpos     Pos
	rpos     Pos
}

func (e *Selector) calculateEnd(p Pos) {
	if len(e.Props) > 0 {
		e.rpos = e.Props[len(e.Props)-1].End()
	} else {
		if p > 0 {
			p -= 1
		}
		e.rpos = p
	}
}

func (e *Selector) needExpand() (ok bool) {
	var n int
	var i int
	for k := range e.Args {
		n += len(k)
		if n > 64 {
			ok = true
			break
		}
		if i > 3 {
			break
		}
		i++
	}
	return
}

type BooleanTyp struct {
	Val  bool
	lpos Pos
	rpos Pos
}

func (e *BooleanTyp) format(b io.StringWriter, padding string, inline bool) {
	switch e.Val {
	case true:
		b.WriteString("true")
	case false:
		b.WriteString("false")
	}
	return
}

func (e *BinaryExpr) isExpr()             {}
func (e *ParenExpr) isExpr()              {}
func (e *Selector) isExpr()               {}
func (e *WildcardTyp) isExpr()            {}
func (e *BooleanTyp) isExpr()             {}
func (e *SpeedTyp) isExpr()               {}
func (e *IntTyp) isExpr()                 {}
func (e *FloatTyp) isExpr()               {}
func (e *DurationTyp) isExpr()            {}
func (e *DistanceTyp) isExpr()            {}
func (e *TemperatureTyp) isExpr()         {}
func (e *PressureTyp) isExpr()            {}
func (e *GeometryPointTyp) isExpr()       {}
func (e *GeometryLineTyp) isExpr()        {}
func (e *GeometryPolygonTyp) isExpr()     {}
func (e *GeometryMultiObjectTyp) isExpr() {}
func (e *GeometryCollectionTyp) isExpr()  {}
func (e *ArrayTyp) isExpr()               {}
func (e *StringTyp) isExpr()              {}
func (e *PercentTyp) isExpr()             {}
func (e *Ref) isExpr()                    {}
func (e *Range) isExpr()                  {}
func (t *Trigger) isExpr()                {}
func (e *Ident) isExpr()                  {}
func (e *Assign) isExpr()                 {}
func (e *DateTyp) isExpr()                {}
func (e *TimeTyp) isExpr()                {}
func (e *WeekdayTyp) isExpr()             {}
func (e *MonthTyp) isExpr()               {}

func (e *BinaryExpr) Pos() Pos             { return e.Left.Pos() }
func (e *BinaryExpr) End() Pos             { return e.Right.Pos() }
func (e *ParenExpr) Pos() Pos              { return e.lpos }
func (e *ParenExpr) End() Pos              { return e.Expr.End() }
func (e *Selector) Pos() Pos               { return e.lpos }
func (e *Selector) End() Pos               { return e.rpos }
func (e *WildcardTyp) Pos() Pos            { return e.lpos }
func (e *WildcardTyp) End() Pos            { return e.lpos + 1 }
func (e *BooleanTyp) Pos() Pos             { return e.lpos }
func (e *BooleanTyp) End() Pos             { return e.rpos }
func (e *SpeedTyp) Pos() Pos               { return e.lpos }
func (e *SpeedTyp) End() Pos               { return e.rpos }
func (e *IntTyp) Pos() Pos                 { return e.lpos }
func (e *IntTyp) End() Pos                 { return e.rpos }
func (e *FloatTyp) Pos() Pos               { return e.lpos }
func (e *FloatTyp) End() Pos               { return e.rpos }
func (e *DurationTyp) Pos() Pos            { return e.lpos }
func (e *DurationTyp) End() Pos            { return e.rpos }
func (e *DistanceTyp) Pos() Pos            { return e.lpos }
func (e *DistanceTyp) End() Pos            { return e.rpos }
func (e *TemperatureTyp) Pos() Pos         { return e.lpos }
func (e *TemperatureTyp) End() Pos         { return e.rpos }
func (e *PressureTyp) Pos() Pos            { return e.lpos }
func (e *PressureTyp) End() Pos            { return e.rpos }
func (e *GeometryPointTyp) Pos() Pos       { return e.lpos }
func (e *GeometryPointTyp) End() Pos       { return e.rpos }
func (e *GeometryLineTyp) Pos() Pos        { return e.lpos }
func (e *GeometryLineTyp) End() Pos        { return e.rpos }
func (e *GeometryPolygonTyp) Pos() Pos     { return e.lpos }
func (e *GeometryPolygonTyp) End() Pos     { return e.rpos }
func (e *GeometryMultiObjectTyp) Pos() Pos { return e.lpos }
func (e *GeometryMultiObjectTyp) End() Pos { return e.rpos }
func (e *GeometryCollectionTyp) Pos() Pos  { return e.lpos }
func (e *GeometryCollectionTyp) End() Pos  { return e.rpos }
func (e *ArrayTyp) Pos() Pos               { return e.lpos }
func (e *ArrayTyp) End() Pos               { return e.rpos }
func (e *StringTyp) Pos() Pos              { return e.lpos }
func (e *StringTyp) End() Pos              { return e.rpos }
func (e *PercentTyp) Pos() Pos             { return e.lpos }
func (e *PercentTyp) End() Pos             { return e.rpos }
func (e *Ref) Pos() Pos                    { return e.lpos }
func (e *Ref) End() Pos                    { return e.rpos }
func (e *Range) Pos() Pos                  { return e.lpos }
func (e *Range) End() Pos                  { return e.rpos }
func (t *Trigger) Pos() Pos                { return t.lpos }
func (t *Trigger) End() Pos                { return t.rpos }
func (e *Ident) Pos() Pos                  { return e.lpos }
func (e *Ident) End() Pos                  { return e.rpos }
func (e *Assign) Pos() Pos                 { return e.Left.Pos() }
func (e *Assign) End() Pos                 { return e.Right.End() }
func (e *DateTyp) Pos() Pos                { return e.lpos }
func (e *DateTyp) End() Pos                { return e.rpos }
func (e *TimeTyp) Pos() Pos                { return e.lpos }
func (e *TimeTyp) End() Pos                { return e.rpos }
func (e *WeekdayTyp) Pos() Pos             { return e.lpos }
func (e *WeekdayTyp) End() Pos             { return e.rpos }
func (e *MonthTyp) Pos() Pos               { return e.lpos }
func (e *MonthTyp) End() Pos               { return e.rpos }
