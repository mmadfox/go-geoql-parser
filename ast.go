package geoqlparser

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type Statement interface {
	isStatement()
}

type Expr interface {
	Pos() Pos
	End() Pos

	format(b *bytes.Buffer, padding string, inline bool)
	isExpr()
}

func (t *TriggerStmt) isStatement() {}

// TriggerStmt represents a TRIGGER statement.
type TriggerStmt struct {
	Vars           []*AssignStmt
	When           Expr
	RepeatCount    Expr
	RepeatInterval Expr
	ResetAfter     Expr
	lpos           Pos
	rpos           Pos
}

func (t *TriggerStmt) SetVar(v *AssignStmt) error {
	if t.isAssigned(v.Left.Val) {
		return fmt.Errorf("variable %s already assigned", v.Left.Val)
	}
	t.Vars = append(t.Vars, v)
	return nil
}

func (t *TriggerStmt) isAssigned(varname string) bool {
	for i := 0; i < len(t.Vars); i++ {
		if t.Vars[i].Left.Val == varname {
			return true
		}
	}
	return false
}

func (t *TriggerStmt) initVars() {
	if t.Vars != nil {
		return
	}
	t.Vars = make([]*AssignStmt, 0)
}

func (t *TriggerStmt) findVar(varname string) (*AssignStmt, error) {
	for i := 0; i < len(t.Vars); i++ {
		if t.Vars[i].Left.Val == varname {
			return t.Vars[i], nil
		}
	}
	return nil, fmt.Errorf("variable %s not found", varname)
}

type AssignStmt struct {
	Left   *IdentLit
	Right  Expr
	TokPos Pos
}

type ArrayExpr struct {
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

type WildcardLit struct {
	lpos Pos
}

func (e *WildcardLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('*')
}

type CalendarLit struct {
	Kind                                    Token
	Year, Day, Hours, Minutes, Seconds, Val int
	Month                                   time.Month
	U                                       Unit
	lpos                                    Pos
	rpos                                    Pos
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
}

var shortMonthNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

func (e *CalendarLit) format(b *bytes.Buffer, _ string, _ bool) {
	switch e.Kind {
	case WEEKDAY:
		b.WriteString(shortDayNames[e.Val])
	case MONTH:
		b.WriteString(shortMonthNames[e.Val])
	}
}

type GeometryPointExpr struct {
	Val    [2]float64
	Radius *DistanceLit
	lpos   Pos
	rpos   Pos
}

type GeometryMultiObject struct {
	Kind Token
	Val  []Expr
	lpos Pos
	rpos Pos
}

type GeometryLineExpr struct {
	Val    [][2]float64
	Margin *DistanceLit
	lpos   Pos
	rpos   Pos
}

type GeometryPolygonExpr struct {
	Val  [][][2]float64
	lpos Pos
	rpos Pos
}

func (e *GeometryPolygonExpr) HasHoles() bool {
	return len(e.Val) > 1
}

type GeometryCollectionExpr struct {
	Objects []Expr
	lpos    Pos
	rpos    Pos
}

type IntLit struct {
	Val  int
	lpos Pos
	rpos Pos
}

func (e *IntLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(strconv.Itoa(e.Val))
}

type RangeExpr struct {
	Low  Expr
	High Expr
	lpos Pos
	rpos Pos
}

type PercentLit struct {
	Val  float64
	lpos Pos
	rpos Pos
}

type StringLit struct {
	Val  string
	lpos Pos
	rpos Pos
}

type FloatLit struct {
	Val  float64
	lpos Pos
	rpos Pos
}

type DurationLit struct {
	Val  time.Duration
	lpos Pos
	rpos Pos
}

type TemperatureLit struct {
	Val  float64
	U    Unit
	Vec  Sign
	lpos Pos
	rpos Pos
}

type PressureLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type DistanceLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type SpeedLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

type IdentLit struct {
	Val  string
	lpos Pos
	rpos Pos
}

func dt2str(v int) string {
	s := strconv.Itoa(v)
	if v < 10 {
		return "0" + s
	}
	return s
}

//func (e *DateLit) format(b *bytes.Buffer, _ string, _ bool) {
//	b.WriteString(strconv.Itoa(e.Year))
//	b.WriteRune('-')
//	b.WriteString(dt2str(int(e.Month)))
//	b.WriteRune('-')
//	b.WriteString(dt2str(e.Day))
//}

//type TimeLit struct {
//	Hour, Minute, Seconds int
//	U                     Unit
//	lpos                  Pos
//	rpos                  Pos
//}

//func (e *TimeLit) format(b *bytes.Buffer, _ string, _ bool) {
//	switch e.U {
//	case AM, PM:
//		b.WriteString(strconv.Itoa(e.Hour))
//		b.WriteRune(':')
//		b.WriteString(dt2str(e.Minute))
//		if e.Seconds > 0 {
//			b.WriteRune(':')
//			b.WriteString(dt2str(e.Seconds))
//		}
//		b.WriteString(e.U.String())
//	default:
//		b.WriteString(dt2str(e.Hour))
//		b.WriteRune(':')
//		b.WriteString(dt2str(e.Minute))
//		if e.Seconds > 0 {
//			b.WriteRune(':')
//			b.WriteString(dt2str(e.Seconds))
//		}
//	}
//}

type RefLit struct {
	ID   string
	lpos Pos
	rpos Pos
}

func (e *RefLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('@')
	b.WriteString(e.ID)
}

//type DateTimeLit struct {
//	Year, Day, Hours, Minutes, Seconds int
//	Month                              time.Month
//	U                                  Unit
//	lpos                               Pos
//	rpos                               Pos
//}

//func (e *DateTimeLit) format(b *bytes.Buffer, _ string, _ bool) {
//	b.WriteString(strconv.Itoa(e.Year))
//	b.WriteRune('-')
//	b.WriteString(dt2str(int(e.Month)))
//	b.WriteRune('-')
//	b.WriteString(dt2str(e.Day))
//	b.WriteRune('T')
//	b.WriteString(dt2str(e.Hours))
//	b.WriteRune(':')
//	b.WriteString(dt2str(e.Minutes))
//	if e.Seconds > 0 {
//		b.WriteRune(':')
//		b.WriteString(dt2str(e.Seconds))
//	}
//}

type SelectorExpr struct {
	Ident    string              // selector name
	Args     map[string]struct{} // device ids
	Wildcard bool                // indicates the current device
	Props    []Expr              // some props
	lpos     Pos
	rpos     Pos
}

func (e *SelectorExpr) calculateEnd(p Pos) {
	if len(e.Props) > 0 {
		e.rpos = e.Props[len(e.Props)-1].End()
	} else {
		if p > 0 {
			p -= 1
		}
		e.rpos = p
	}
}

func (e *SelectorExpr) needExpand() (ok bool) {
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

type BooleanLit struct {
	Val  bool
	lpos Pos
	rpos Pos
}

func (e *BooleanLit) format(b *bytes.Buffer, _ string, _ bool) {
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
func (e *SelectorExpr) isExpr()           {}
func (e *WildcardLit) isExpr()            {}
func (e *BooleanLit) isExpr()             {}
func (e *SpeedLit) isExpr()               {}
func (e *IntLit) isExpr()                 {}
func (e *FloatLit) isExpr()               {}
func (e *DurationLit) isExpr()            {}
func (e *DistanceLit) isExpr()            {}
func (e *TemperatureLit) isExpr()         {}
func (e *PressureLit) isExpr()            {}
func (e *GeometryPointExpr) isExpr()      {}
func (e *GeometryLineExpr) isExpr()       {}
func (e *GeometryPolygonExpr) isExpr()    {}
func (e *GeometryMultiObject) isExpr()    {}
func (e *GeometryCollectionExpr) isExpr() {}
func (e *ArrayExpr) isExpr()              {}
func (e *StringLit) isExpr()              {}
func (e *PercentLit) isExpr()             {}
func (e *RefLit) isExpr()                 {}
func (e *RangeExpr) isExpr()              {}
func (e *CalendarLit) isExpr()            {}
func (t *TriggerStmt) isExpr()            {}
func (e *IdentLit) isExpr()               {}
func (e *AssignStmt) isExpr()             {}

func (e *BinaryExpr) Pos() Pos             { return e.Left.Pos() }
func (e *BinaryExpr) End() Pos             { return e.Right.Pos() }
func (e *ParenExpr) Pos() Pos              { return e.lpos }
func (e *ParenExpr) End() Pos              { return e.Expr.End() }
func (e *SelectorExpr) Pos() Pos           { return e.lpos }
func (e *SelectorExpr) End() Pos           { return e.rpos }
func (e *WildcardLit) Pos() Pos            { return e.lpos }
func (e *WildcardLit) End() Pos            { return e.lpos + 1 }
func (e *BooleanLit) Pos() Pos             { return e.lpos }
func (e *BooleanLit) End() Pos             { return e.rpos }
func (e *SpeedLit) Pos() Pos               { return e.lpos }
func (e *SpeedLit) End() Pos               { return e.rpos }
func (e *IntLit) Pos() Pos                 { return e.lpos }
func (e *IntLit) End() Pos                 { return e.rpos }
func (e *FloatLit) Pos() Pos               { return e.lpos }
func (e *FloatLit) End() Pos               { return e.rpos }
func (e *DurationLit) Pos() Pos            { return e.lpos }
func (e *DurationLit) End() Pos            { return e.rpos }
func (e *DistanceLit) Pos() Pos            { return e.lpos }
func (e *DistanceLit) End() Pos            { return e.rpos }
func (e *TemperatureLit) Pos() Pos         { return e.lpos }
func (e *TemperatureLit) End() Pos         { return e.rpos }
func (e *PressureLit) Pos() Pos            { return e.lpos }
func (e *PressureLit) End() Pos            { return e.rpos }
func (e *GeometryPointExpr) Pos() Pos      { return e.lpos }
func (e *GeometryPointExpr) End() Pos      { return e.rpos }
func (e *GeometryLineExpr) Pos() Pos       { return e.lpos }
func (e *GeometryLineExpr) End() Pos       { return e.rpos }
func (e *GeometryPolygonExpr) Pos() Pos    { return e.lpos }
func (e *GeometryPolygonExpr) End() Pos    { return e.rpos }
func (e *GeometryMultiObject) Pos() Pos    { return e.lpos }
func (e *GeometryMultiObject) End() Pos    { return e.rpos }
func (e *GeometryCollectionExpr) Pos() Pos { return e.lpos }
func (e *GeometryCollectionExpr) End() Pos { return e.rpos }
func (e *ArrayExpr) Pos() Pos              { return e.lpos }
func (e *ArrayExpr) End() Pos              { return e.rpos }
func (e *StringLit) Pos() Pos              { return e.lpos }
func (e *StringLit) End() Pos              { return e.rpos }
func (e *PercentLit) Pos() Pos             { return e.lpos }
func (e *PercentLit) End() Pos             { return e.rpos }
func (e *RefLit) Pos() Pos                 { return e.lpos }
func (e *RefLit) End() Pos                 { return e.rpos }
func (e *RangeExpr) Pos() Pos              { return e.lpos }
func (e *RangeExpr) End() Pos              { return e.rpos }
func (e *CalendarLit) Pos() Pos            { return e.lpos }
func (e *CalendarLit) End() Pos            { return e.rpos }
func (t *TriggerStmt) Pos() Pos            { return t.lpos }
func (t *TriggerStmt) End() Pos            { return t.rpos }
func (e *IdentLit) Pos() Pos               { return e.lpos }
func (e *IdentLit) End() Pos               { return e.rpos }
func (e *AssignStmt) Pos() Pos             { return e.Left.Pos() }
func (e *AssignStmt) End() Pos             { return e.Right.End() }
