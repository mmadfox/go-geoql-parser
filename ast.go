package geoqlparser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Statement interface {
	isStatement()
}

type Expr interface {
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
	Pos            Pos
}

func (t *TriggerStmt) String() string {
	if t == nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	formatTriggerStmt(t, buf)
	return buf.String()
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
	Kind     Token
	List     []Expr
	StartPos Pos
	EndPos   Pos
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
	Pos   Pos
}

func (e *BinaryExpr) format(b *bytes.Buffer, padding string, inline bool) {
	e.Left.format(b, padding, inline)
	var nospace bool
	op := KeywordString(e.Op)
	var nl bool
	switch e.Op {
	case AND, OR:
		nl = true
	case ADD, SUB, MUL, QUO:
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
	Expr     Expr
	StartPos Pos
	EndPos   Pos
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
	Pos Pos
}

func (e *WildcardLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('*')
}

type GeometryPointExpr struct {
	Val      [2]float64
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryPointExpr) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString("point")
	b.WriteRune('[')
	formatFloat(b, e.Val[0])
	b.WriteRune(',')
	b.WriteRune(' ')
	formatFloat(b, e.Val[1])
	b.WriteRune(']')
}

type GeometryMultiPointExpr struct {
	Val      [][2]float64
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryMultiPointExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("multipoint")
	b.WriteRune('[')
	var pad2 string
	for i := 0; i < len(e.Val); i++ {
		if !inline {
			pad2 = padding + "\t"
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
}

type GeometryLineExpr struct {
	Val      [][2]float64
	StartPos Pos
	EndPos   Pos
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
}

type GeometryMultiLineExpr struct {
	Val      [][][2]float64
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryMultiLineExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("multiline")
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

type GeometryPolygonExpr struct {
	Val      [][2]float64
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryPolygonExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("polygon")
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
}

type GeometryCircleExpr struct {
	Val      [2]float64
	Radius   *DistanceLit
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryCircleExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("circle")
	b.WriteRune('[')
	formatFloat(b, e.Val[0])
	b.WriteRune(',')
	b.WriteRune(' ')
	formatFloat(b, e.Val[1])
	b.WriteRune(']')
	b.WriteRune(':')
	e.Radius.format(b, padding, inline)
}

type GeometryMultiPolygonExpr struct {
	Val      [][][2]float64
	StartPos Pos
	EndPos   Pos
}

func (e *GeometryMultiPolygonExpr) format(b *bytes.Buffer, padding string, inline bool) {
	b.WriteString("multipolygon")
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
	Objects  []Expr
	StartPos Pos
	EndPos   Pos
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
	Val int
	Pos Pos
}

func (e *IntLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(strconv.Itoa(e.Val))
}

type RangeExpr struct {
	Low      Expr
	High     Expr
	StartPos Pos
	EndPos   Pos
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
	Val float64
	Pos Pos
}

func (e *PercentLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(fmt.Sprintf("%.2f", e.Val))
	b.WriteString("%")
}

type StringLit struct {
	Val string
	Pos Pos
}

func (e *StringLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('"')
	b.WriteString(e.Val)
	b.WriteRune('"')
}

type FloatLit struct {
	Val float64
	Pos Pos
}

func (e *FloatLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
}

type DurationLit struct {
	Val time.Duration
	Pos Pos
}

func (e *DurationLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(e.Val.String())
}

type TemperatureLit struct {
	Val float64
	U   Unit
	Vec Vector
	Pos Pos
}

func (e *TemperatureLit) format(b *bytes.Buffer, _ string, _ bool) {
	if e.Val != 0 {
		b.WriteString(e.Vec.String())
	}
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type PressureLit struct {
	Val float64
	U   Unit
	Pos Pos
}

func (e *PressureLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type DistanceLit struct {
	Val float64
	U   Unit
	Pos Pos
}

func (e *DistanceLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type SpeedLit struct {
	Val float64
	U   Unit
	Pos Pos
}

func (e *SpeedLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

type DateLit struct {
	Year, Day int
	Month     time.Month
	Pos       Pos
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
	Pos                   Pos
}

func (e *TimeLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(dt2str(e.Hour))
	b.WriteRune(':')
	b.WriteString(dt2str(e.Minute))
	if e.Seconds > 0 {
		b.WriteRune(':')
		b.WriteString(dt2str(e.Seconds))
	}
}

type VarLit struct {
	ID  string
	Pos Pos
}

func (e *VarLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('@')
	b.WriteString(e.ID)
}

type DateTimeLit struct {
	Year, Day, Hours, Minutes, Seconds int
	Month                              time.Month
	Pos                                Pos
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
	Ident    string
	Args     map[string]struct{}
	Wildcard bool
	Props    []Expr
	StartPos Pos
	EndPos   Pos
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
	Val bool
	Pos Pos
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
		*GeometryCircleExpr, *GeometryPolygonExpr, *GeometryMultiPolygonExpr,
		*GeometryMultiLineExpr, *GeometryMultiPointExpr, *GeometryLineExpr:
		ok = true
	}
	return
}

func (e *BinaryExpr) isExpr()               {}
func (e *ParenExpr) isExpr()                {}
func (e *SelectorExpr) isExpr()             {}
func (e *WildcardLit) isExpr()              {}
func (e *BooleanLit) isExpr()               {}
func (e *SpeedLit) isExpr()                 {}
func (e *IntLit) isExpr()                   {}
func (e *FloatLit) isExpr()                 {}
func (e *DurationLit) isExpr()              {}
func (e *DistanceLit) isExpr()              {}
func (e *TemperatureLit) isExpr()           {}
func (e *PressureLit) isExpr()              {}
func (e *GeometryPointExpr) isExpr()        {}
func (e *GeometryMultiPointExpr) isExpr()   {}
func (e *GeometryLineExpr) isExpr()         {}
func (e *GeometryMultiLineExpr) isExpr()    {}
func (e *GeometryPolygonExpr) isExpr()      {}
func (e *GeometryMultiPolygonExpr) isExpr() {}
func (e *GeometryCircleExpr) isExpr()       {}
func (e *GeometryCollectionExpr) isExpr()   {}
func (e *DateLit) isExpr()                  {}
func (e *TimeLit) isExpr()                  {}
func (e *DateTimeLit) isExpr()              {}
func (e *ArrayExpr) isExpr()                {}
func (e *StringLit) isExpr()                {}
func (e *PercentLit) isExpr()               {}
func (e *VarLit) isExpr()                   {}
func (e *RangeExpr) isExpr()                {}
