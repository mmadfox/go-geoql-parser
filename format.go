package geoqlparser

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

const (
	nl      = "\n"
	padding = "\t"
)

func Format(stmt Statement, w io.StringWriter) {
	switch typ := stmt.(type) {
	case *TriggerStmt:
		formatTriggerStmt(typ, w)
	}
}

func formatTriggerStmt(t *TriggerStmt, w io.StringWriter) {
	w.WriteString("TRIGGER")
	w.WriteString(nl)

	//if len(t.Vars) > 0 {
	//	format := func(expr Expr, vname string) {
	//		b.WriteString(padding)
	//		b.WriteString(vname)
	//		b.WriteRune(' ')
	//		b.WriteRune('=')
	//		b.WriteRune(' ')
	//		expr.format(b, padding, false)
	//		b.WriteRune('\n')
	//	}
	//	b.WriteString("SET")
	//	b.WriteRune('\n')
	//for vname, expr := range t.Vars {
	//	if IsGeometryExpr(expr) {
	//		continue
	//	}
	//	format(expr, vname)
	//}
	//for vname, expr := range t.Vars {
	//	if !IsGeometryExpr(expr) {
	//		continue
	//	}
	//	format(expr, vname)
	//}
	//}
	//b.WriteString("WHEN")
	//b.WriteRune('\n')
	//b.WriteString(padding)
	//t.When.format(b, padding, true)
	//b.WriteRune('\n')
	//
	//if t.RepeatCount != nil && t.RepeatInterval != nil {
	//	b.WriteString("REPEAT")
	//	b.WriteRune(' ')
	//	t.RepeatCount.format(b, padding, true)
	//	b.WriteRune(' ')
	//	b.WriteString("times")
	//	b.WriteRune(' ')
	//	t.RepeatInterval.format(b, padding, true)
	//	b.WriteRune(' ')
	//	b.WriteString("interval")
	//	b.WriteRune('\n')
	//}
	//
	//if t.ResetAfter != nil {
	//	b.WriteString("RESET after")
	//	b.WriteRune(' ')
	//	t.ResetAfter.format(b, padding, true)
	//	b.WriteRune('\n')
	//}
}

func (t *TriggerStmt) format(b *bytes.Buffer, _ string, _ bool) {
	formatTriggerStmt(t, b)
}

func formatFloat(b *bytes.Buffer, v float64) {
	b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
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

func (e *RangeExpr) format(b *bytes.Buffer, padding string, inline bool) {
	e.Low.format(b, padding, inline)
	b.WriteRune(' ')
	b.WriteRune('.')
	b.WriteRune('.')
	b.WriteRune(' ')
	e.High.format(b, padding, inline)
}

func (e *PercentLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString("%")
}

func (e *StringLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteRune('"')
	b.WriteString(e.Val)
	b.WriteRune('"')
}

func (e *FloatLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
}

func (e *DurationLit) format(b *bytes.Buffer, _ string, _ bool) {
	b.WriteString(e.Val.String())
}

func (e *TemperatureLit) format(b *bytes.Buffer, _ string, _ bool) {
	if e.Val != 0 {
		b.WriteString(e.Vec.String())
	}
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

func (e *PressureLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

func (e *DistanceLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
}

func (e *SpeedLit) format(b *bytes.Buffer, _ string, _ bool) {
	formatFloat(b, e.Val)
	b.WriteString(e.U.String())
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

func (e *AssignStmt) format(b *bytes.Buffer, padding string, inline bool) {}

func (e *IdentLit) format(b *bytes.Buffer, padding string, inline bool) {}
