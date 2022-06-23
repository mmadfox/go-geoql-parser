package geoqlparser

import (
	"bytes"
	"strconv"
	"strings"
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
	Set            map[string]Expr
	When           Expr
	RepeatCount    Expr
	RepeatInterval Expr
	ResetAfter     Expr
	lpos           Pos
	rpos           Pos
}

func (t *TriggerStmt) String() string {
	if t == nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	formatTriggerStmt(t, buf)
	return buf.String()
}

func (t *TriggerStmt) format(b *bytes.Buffer, _ string, _ bool) {
	formatTriggerStmt(t, b)
}

func formatTriggerStmt(t *TriggerStmt, b *bytes.Buffer) {
	padding := "\t"
	b.WriteString("TRIGGER")
	b.WriteRune('\n')
	if len(t.Set) > 0 {
		format := func(expr Expr, vname string) {
			b.WriteString(padding)
			b.WriteString(vname)
			b.WriteRune(' ')
			b.WriteRune('=')
			b.WriteRune(' ')
			expr.format(b, padding, false)
			b.WriteRune('\n')
		}
		b.WriteString("SET")
		b.WriteRune('\n')
		for vname, expr := range t.Set {
			if IsGeometryExpr(expr) {
				continue
			}
			format(expr, vname)
		}
		for vname, expr := range t.Set {
			if !IsGeometryExpr(expr) {
				continue
			}
			format(expr, vname)
		}
	}
	b.WriteString("WHEN")
	b.WriteRune('\n')
	b.WriteString(padding)
	t.When.format(b, padding, true)
	b.WriteRune('\n')

	if t.RepeatCount != nil && t.RepeatInterval != nil {
		b.WriteString("REPEAT")
		b.WriteRune(' ')
		t.RepeatCount.format(b, padding, true)
		b.WriteRune(' ')
		b.WriteString("times")
		b.WriteRune(' ')
		t.RepeatInterval.format(b, padding, true)
		b.WriteRune(' ')
		b.WriteString("interval")
		b.WriteRune('\n')
	}

	if t.ResetAfter != nil {
		b.WriteString("RESET after")
		b.WriteRune(' ')
		t.ResetAfter.format(b, padding, true)
		b.WriteRune('\n')
	}
}

func formatFloat(b *bytes.Buffer, v float64) {
	b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
}

func (t *TriggerStmt) initVars() {
	if t.Set != nil {
		return
	}
	t.Set = make(map[string]Expr)
}

type ArrayExpr struct {
	Kind Token
	List []Expr
	lpos Pos
	rpos Pos
}

func (e *ArrayExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteRune('[')
	for i, expr := range e.List {
		expr.format(b, padding, inline)
		if i+1 < len(e.List) {
			b.WriteRune(',')
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')
}

type BinaryExpr struct {
	Op    Token
	Left  Expr
	Right Expr
	OpPos Pos
}

func (e *BinaryExpr) format(b *bytes.Buffer, padding string, inline bool) {
	e.Left.format(b, padding, inline)
	var nospace bool
	op := KeywordString(e.Op)
	var nl bool
	switch e.Op {
	case AND, OR:
		nl = true
	case ADD, SUB, MUL, QUO, REM:
		nospace = true
	}
	if !nospace {
		b.WriteRune(' ')
	}
	if nl {
		b.WriteRune('\n')
		b.WriteString(padding)
	}
	b.WriteString(op)
	if !nospace {
		b.WriteRune(' ')
	}
	e.Right.format(b, padding, inline)
}

type ParenExpr struct {
	Expr Expr
	lpos Pos
	rpos Pos
}

func (e *ParenExpr) format(b *bytes.Buffer, padding string, inline bool) {
	expand := true
	switch node := e.Expr.(type) {
	case *BinaryExpr:
		_, lok := node.Left.(*BinaryExpr)
		_, rok := node.Right.(*BinaryExpr)
		if !lok && !rok {
			expand = false
		}
	}
	pad2 := padding + "\t"
	b.WriteRune('(')
	if expand {
		b.WriteRune('\n')
		b.WriteString(pad2)
	}
	e.Expr.format(b, pad2, inline)
	if expand {
		b.WriteRune('\n')
		b.WriteString(padding)
	}
	b.WriteRune(')')
}

type WildcardLit struct {
	lpos Pos
}

func (e *WildcardLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('*')
}

type CalendarLit struct {
	Kind Token
	Val  int
	lpos Pos
	rpos Pos
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

func (e *GeometryPointExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("point")
	b.WriteRune('[')
	formatFloat(b, e.Val[0])
	b.WriteRune(',')
	b.WriteRune(' ')
	formatFloat(b, e.Val[1])
	b.WriteRune(']')
	if e.Radius != nil {
		b.WriteRune(':')
		e.Radius.format(b, padding, inline)
	}
}

type GeometryMultiObject struct {
	Kind Token
	Val  []Expr
	lpos Pos
	rpos Pos
}

func (e *GeometryMultiObject) format(b *bytes.Buffer, padding string, inline bool) {
	switch e.Kind {
	case GEOMETRY_MULTIPOINT:
		b.WriteString("multipoint")
	case GEOMETRY_MULTILINE:
		b.WriteString("multiline")
	case GEOMETRY_MULTIPOLYGON:
		b.WriteString("multipolygon")
	}
	b.WriteRune('[')
	pad2 := padding + "\t"
	for i := 0; i < len(e.Val); i++ {
		if !inline {
			b.WriteRune('\n')
			if e.Kind == GEOMETRY_MULTIPOINT {
				b.WriteString(pad2)
			} else {
				b.WriteString(padding)
			}
		}
		e.Val[i].format(b, padding, inline)
		if i+1 < len(e.Val) {
			b.WriteRune(',')
			b.WriteRune(' ')
		}
	}
	if !inline && e.Kind == GEOMETRY_MULTIPOINT {
		b.WriteRune('\n')
		b.WriteString(padding)
	}
	b.WriteRune(']')
}

type GeometryLineExpr struct {
	Val    [][2]float64
	Margin *DistanceLit
	lpos   Pos
	rpos   Pos
}

func (e *GeometryLineExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("line")
	b.WriteRune('[')
	var pad2 string
	if !inline {
		pad2 = padding + "\t"
	}
	for i := 0; i < len(e.Val); i++ {
		if !inline {
			b.WriteRune('\n')
			b.WriteString(pad2)
		}
		b.WriteRune('[')
		formatFloat(b, e.Val[i][0])
		b.WriteRune(',')
		b.WriteRune(' ')
		formatFloat(b, e.Val[i][1])
		b.WriteRune(']')
		if i+1 < len(e.Val) {
			b.WriteRune(',')
			b.WriteRune(' ')
		}
	}
	if !inline {
		b.WriteRune('\n')
		b.WriteString(padding)
	}
	b.WriteRune(']')
	if e.Margin != nil {
		b.WriteRune(':')
		e.Margin.format(b, padding, inline)
	}
}

type GeometryPolygonExpr struct {
	Val  [][][2]float64
	lpos Pos
	rpos Pos
}

func (e *GeometryPolygonExpr) HasHoles() bool {
	return len(e.Val) > 1
}

func (e *GeometryPolygonExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("polygon")
	b.WriteRune('[')
	var pad2, pad3 string
	if !inline {
		pad2 = padding + "\t"
		pad3 = padding + "\t\t"
		b.WriteRune('\n')
		b.WriteString(pad2)
	}
	for i := 0; i < len(e.Val); i++ {
		b.WriteRune('[')
		for j := 0; j < len(e.Val[i]); j++ {
			if !inline {
				b.WriteRune('\n')
				b.WriteString(pad3)
			}
			b.WriteRune('[')
			formatFloat(b, e.Val[i][j][0])
			b.WriteRune(',')
			b.WriteRune(' ')
			formatFloat(b, e.Val[i][j][1])
			b.WriteRune(']')
			if j+1 < len(e.Val[i]) {
				b.WriteRune(',')
				b.WriteRune(' ')
			}
		}
		if !inline {
			b.WriteRune('\n')
			b.WriteString(pad2)
		}
		b.WriteRune(']')
		if i+1 < len(e.Val) {
			b.WriteRune(',')
			b.WriteRune(' ')
		}
	}
	if !inline {
		b.WriteRune('\n')
		b.WriteString(padding)
	}
	b.WriteRune(']')
}

type GeometryCollectionExpr struct {
	Objects []Expr
	lpos    Pos
	rpos    Pos
}

func (e *GeometryCollectionExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("collection")
	b.WriteRune('[')
	if !inline {
		b.WriteRune('\n')
	}
	newpad := strings.Repeat(padding, 2)
	for i := 0; i < len(e.Objects); i++ {
		if !inline {
			b.WriteString(newpad)
		}
		e.Objects[i].format(b, newpad, inline)
		if i+1 < len(e.Objects) {
			b.WriteRune(',')
			b.WriteRune(' ')
		}
		if !inline {
			b.WriteRune('\n')
		}
	}
	if !inline {
		b.WriteString(padding)
	}
	b.WriteRune(']')
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

func (e *RangeExpr) format(b *bytes.Buffer, padding string, inline bool) {
	e.Low.format(b, padding, inline)
	b.WriteRune(' ')
	b.WriteRune('.')
	b.WriteRune('.')
	b.WriteRune(' ')
	e.High.format(b, padding, inline)
}

type PercentLit struct {
	Val  float64
	lpos Pos
	rpos Pos
}

func (e *PercentLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString("%")
}

type StringLit struct {
	Val  string
	lpos Pos
	rpos Pos
}

func (e *StringLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('"')
	b.WriteString(e.Val)
	b.WriteRune('"')
}

type FloatLit struct {
	Val  float64
	lpos Pos
	rpos Pos
}

func (e *FloatLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
}

type DurationLit struct {
	Val  time.Duration
	lpos Pos
	rpos Pos
}

func (e *DurationLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(e.Val.String())
}

type TemperatureLit struct {
	Val  float64
	U    Unit
	Vec  Vector
	lpos Pos
	rpos Pos
}

func (e *TemperatureLit) format(b *bytes.Buffer, _ string, _ bool) {
	if e.Val != 0 {
		b.WriteString(e.Vec.String())
	}
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type PressureLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

func (e *PressureLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type DistanceLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

func (e *DistanceLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type SpeedLit struct {
	Val  float64
	U    Unit
	lpos Pos
	rpos Pos
}

func (e *SpeedLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type DateLit struct {
	Year, Day int
	Month     time.Month
	lpos      Pos
	rpos      Pos
}

func dt2str(v int) string {
	s := strconv.Itoa(v)
	if v < 10 {
		return "0" + s
	}
	return s
}

func (e *DateLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(strconv.Itoa(e.Year))
	b.WriteRune('-')
	b.WriteString(dt2str(int(e.Month)))
	b.WriteRune('-')
	b.WriteString(dt2str(e.Day))
}

type TimeLit struct {
	Hour, Minute, Seconds int
	U                     Unit
	lpos                  Pos
	rpos                  Pos
}

func (e *TimeLit) format(b *bytes.Buffer, _ string, _ bool) {
	switch e.U {
	case AM, PM:
		b.WriteString(strconv.Itoa(e.Hour))
		b.WriteRune(':')
		b.WriteString(dt2str(e.Minute))
		if e.Seconds > 0 {
			b.WriteRune(':')
			b.WriteString(dt2str(e.Seconds))
		}
		b.WriteString(e.U.String())
	default:
		b.WriteString(dt2str(e.Hour))
		b.WriteRune(':')
		b.WriteString(dt2str(e.Minute))
		if e.Seconds > 0 {
			b.WriteRune(':')
			b.WriteString(dt2str(e.Seconds))
		}
	}
}

type VarLit struct {
	ID   string
	lpos Pos
	rpos Pos
}

func (e *VarLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('@')
	b.WriteString(e.ID)
}

type DateTimeLit struct {
	Year, Day, Hours, Minutes, Seconds int
	Month                              time.Month
	U                                  Unit
	lpos                               Pos
	rpos                               Pos
}

func (e *DateTimeLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(strconv.Itoa(e.Year))
	b.WriteRune('-')
	b.WriteString(dt2str(int(e.Month)))
	b.WriteRune('-')
	b.WriteString(dt2str(e.Day))
	b.WriteRune('T')
	b.WriteString(dt2str(e.Hours))
	b.WriteRune(':')
	b.WriteString(dt2str(e.Minutes))
	if e.Seconds > 0 {
		b.WriteRune(':')
		b.WriteString(dt2str(e.Seconds))
	}
}

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
		e.rpos = e.Props[len(e.Props)-1].End() + 1
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

func (e *SelectorExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString(e.Ident)
	if len(e.Args) > 0 {
		b.WriteRune('{')
		var i int
		var expand bool
		var pad2 string
		if !inline {
			expand = e.needExpand()
			if expand {
				pad2 = padding + "\t"
			}
		}
		if e.Wildcard {
			b.WriteRune('*')
			if len(e.Args) > 0 {
				b.WriteRune(',')
				b.WriteRune(' ')
			}
		}
		for k := range e.Args {
			if !inline && expand {
				b.WriteRune('\n')
				b.WriteString(pad2)
			}
			b.WriteRune('"')
			b.WriteString(k)
			b.WriteRune('"')
			if i+1 < len(e.Args) {
				b.WriteRune(',')
				b.WriteRune(' ')
			}
			i++
		}
		if !inline && expand {
			b.WriteRune('\n')
			b.WriteString(padding)
		}
		b.WriteRune('}')
	}
	if len(e.Props) > 0 {
		b.WriteRune(':')
		for i := 0; i < len(e.Props); i++ {
			e.Props[i].format(b, padding, inline)
			if i+1 < len(e.Props) {
				b.WriteRune(',')
			}
		}
	}
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

func IsGeometryExpr(expr Expr) (ok bool) {
	switch expr.(type) {
	case *GeometryPointExpr, *GeometryCollectionExpr,
		*GeometryMultiObject, *GeometryPolygonExpr, *GeometryLineExpr:
		ok = true
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
func (e *DateLit) isExpr()                {}
func (e *TimeLit) isExpr()                {}
func (e *DateTimeLit) isExpr()            {}
func (e *ArrayExpr) isExpr()              {}
func (e *StringLit) isExpr()              {}
func (e *PercentLit) isExpr()             {}
func (e *VarLit) isExpr()                 {}
func (e *RangeExpr) isExpr()              {}
func (e *CalendarLit) isExpr()            {}
func (t *TriggerStmt) isExpr()            {}

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
func (e *DateLit) Pos() Pos                { return e.lpos }
func (e *DateLit) End() Pos                { return e.rpos }
func (e *TimeLit) Pos() Pos                { return e.lpos }
func (e *TimeLit) End() Pos                { return e.rpos }
func (e *DateTimeLit) Pos() Pos            { return e.lpos }
func (e *DateTimeLit) End() Pos            { return e.rpos }
func (e *ArrayExpr) Pos() Pos              { return e.lpos }
func (e *ArrayExpr) End() Pos              { return e.rpos }
func (e *StringLit) Pos() Pos              { return e.lpos }
func (e *StringLit) End() Pos              { return e.rpos }
func (e *PercentLit) Pos() Pos             { return e.lpos }
func (e *PercentLit) End() Pos             { return e.rpos }
func (e *VarLit) Pos() Pos                 { return e.lpos }
func (e *VarLit) End() Pos                 { return e.rpos }
func (e *RangeExpr) Pos() Pos              { return e.lpos }
func (e *RangeExpr) End() Pos              { return e.rpos }
func (e *CalendarLit) Pos() Pos            { return e.lpos }
func (e *CalendarLit) End() Pos            { return e.rpos }
func (t *TriggerStmt) Pos() Pos            { return t.lpos }
func (t *TriggerStmt) End() Pos            { return t.rpos }
