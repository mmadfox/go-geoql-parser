package geoqlparser

import (
	"fmt"
	"io"
	"strconv"
)

const (
	nl      = "\n"
	padding = "\t"
)

func Format(w io.StringWriter, stmt Statement) error {
	switch typ := stmt.(type) {
	default:
		return fmt.Errorf("todo")
	case *TriggerStmt:
		return formatTriggerStmt(typ, w)
	}
}

func checkError(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func writeNewLine(w io.StringWriter) {
	checkError(w.WriteString(nl))
}

func formatTriggerStmt(t *TriggerStmt, w io.StringWriter) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = er.(error)
		}
	}()

	checkError(w.WriteString("TRIGGER"))
	writeNewLine(w)

	if len(t.Vars) > 0 {
		checkError(w.WriteString("SET"))
		writeNewLine(w)
		for i := 0; i < len(t.Vars); i++ {
			t.Vars[i].format(w, padding, false)
			writeNewLine(w)
		}
	}

	checkError(w.WriteString("WHEN\n" + padding))
	t.When.format(w, padding, true)
	writeNewLine(w)

	if t.RepeatCount != nil || t.RepeatInterval != nil {
		checkError(w.WriteString("REPEAT "))
		if t.RepeatCount != nil {
			t.RepeatCount.format(w, padding, true)
		}
		if t.RepeatInterval != nil {
			checkError(w.WriteString(" every "))
			t.RepeatInterval.format(w, padding, true)
		}
		writeNewLine(w)
	}

	if t.ResetAfter != nil {
		checkError(w.WriteString("RESET after "))
		t.ResetAfter.format(w, padding, true)
	}
	return
}

func (t *TriggerStmt) format(_ io.StringWriter, _ string, _ bool) {}

func formatFloat(w io.StringWriter, v float64) {
	checkError(w.WriteString(strconv.FormatFloat(v, 'f', -1, 64)))
}

func (e *GeometryMultiObject) format(w io.StringWriter, padding string, inline bool) {
	switch e.Kind {
	case GEOMETRY_MULTIPOINT:
		checkError(w.WriteString("multipoint["))
	case GEOMETRY_MULTILINE:
		checkError(w.WriteString("multiline["))
	case GEOMETRY_MULTIPOLYGON:
		checkError(w.WriteString("multipolygon["))
	}
	for i := 0; i < len(e.Val); i++ {
		if !inline {
			writeNewLine(w)
			checkError(w.WriteString(padding))
		}
		e.Val[i].format(w, padding, inline)
		if i+1 < len(e.Val) {
			checkError(w.WriteString(", "))
		}
	}
	checkError(w.WriteString("]"))
}

func (e *GeometryLineExpr) format(w io.StringWriter, padding string, inline bool) {
	needExpand := e.needExpand()
	checkError(w.WriteString("line["))
	pad2 := padding + padding
	if needExpand && !inline {
		writeNewLine(w)
		checkError(w.WriteString(pad2))
	}
	var block int
	for i := 0; i < len(e.Val); i++ {
		checkError(w.WriteString("["))
		formatFloat(w, e.Val[i][0])
		checkError(w.WriteString(", "))
		formatFloat(w, e.Val[i][1])
		checkError(w.WriteString("]"))
		if i+1 < len(e.Val) {
			checkError(w.WriteString(", "))
		}
		if block > 3 {
			block = 0
			writeNewLine(w)
			if i+2 < len(e.Val) {
				checkError(w.WriteString(pad2))
			}
			continue
		}
		block++
	}
	if needExpand && !inline && block > 0 {
		writeNewLine(w)
	}
	if needExpand && !inline {
		checkError(w.WriteString(padding))
	}
	checkError(w.WriteString("]"))
	if e.Margin != nil {
		checkError(w.WriteString(":"))
		e.Margin.format(w, padding, inline)
	}
}

func (e *GeometryPolygonExpr) format(w io.StringWriter, padding string, inline bool) {
	needExpand := e.needExpand()
	checkError(w.WriteString("polygon["))
	pad2 := padding + padding
	pad3 := padding + padding + padding
	var block int
	for i := 0; i < len(e.Val); i++ {
		if needExpand && !inline {
			writeNewLine(w)
			checkError(w.WriteString(pad2))
		}
		checkError(w.WriteString("["))
		if needExpand && !inline {
			writeNewLine(w)
			checkError(w.WriteString(pad3))
		}
		for j := 0; j < len(e.Val[i]); j++ {
			checkError(w.WriteString("["))
			formatFloat(w, e.Val[i][j][0])
			checkError(w.WriteString(", "))
			formatFloat(w, e.Val[i][j][1])
			checkError(w.WriteString("]"))
			if j+1 < len(e.Val[i]) {
				checkError(w.WriteString(", "))
			}
			if block > 3 {
				block = 0
				writeNewLine(w)
				checkError(w.WriteString(pad3))
				continue
			}
			block++
		}
		if needExpand && !inline {
			writeNewLine(w)
			checkError(w.WriteString(pad2))
		}
		checkError(w.WriteString("]"))
		if i+1 < len(e.Val) {
			checkError(w.WriteString(", "))
		}
	}
	if needExpand && !inline {
		writeNewLine(w)
		checkError(w.WriteString(padding))
	}
	checkError(w.WriteString("]"))
}

func (e *GeometryCollectionExpr) format(w io.StringWriter, padding string, inline bool) {
	checkError(w.WriteString("collection["))
	for i := 0; i < len(e.Objects); i++ {
		if !inline {
			writeNewLine(w)
			checkError(w.WriteString(padding))
		}
		e.Objects[i].format(w, padding, inline)
		if i+1 < len(e.Objects) {
			checkError(w.WriteString(", "))
		}
	}
	checkError(w.WriteString("]"))
}

func (e *RangeExpr) format(w io.StringWriter, padding string, inline bool) {
	e.Low.format(w, padding, inline)
	checkError(w.WriteString(" ... "))
	e.High.format(w, padding, inline)
}

func (e *PercentLit) format(w io.StringWriter, _ string, _ bool) {
	formatFloat(w, e.Val)
	checkError(w.WriteString("%"))
}

func (e *StringLit) format(w io.StringWriter, _ string, _ bool) {
	checkError(w.WriteString(`"`))
	checkError(w.WriteString(e.Val))
	checkError(w.WriteString(`"`))
}

func (e *FloatLit) format(w io.StringWriter, _ string, _ bool) {
	formatFloat(w, e.Val)
}

func (e *DurationLit) format(b io.StringWriter, _ string, _ bool) {
	checkError(b.WriteString(e.Val.String()))
}

func (e *TemperatureLit) format(w io.StringWriter, _ string, _ bool) {
	switch e.Vec {
	case Plus:
		checkError(w.WriteString("+"))
	case Minus:
		checkError(w.WriteString("-"))
	}
	formatFloat(w, e.Val)
	checkError(w.WriteString(e.U.String()))
}

func (e *PressureLit) format(w io.StringWriter, _ string, _ bool) {
	formatFloat(w, e.Val)
	checkError(w.WriteString(e.U.String()))
}

func (e *DistanceLit) format(w io.StringWriter, _ string, _ bool) {
	formatFloat(w, e.Val)
	checkError(w.WriteString(e.U.String()))
}

func (e *SpeedLit) format(w io.StringWriter, _ string, _ bool) {
	formatFloat(w, e.Val)
	checkError(w.WriteString(e.U.String()))
}

func (e *SelectorExpr) format(b io.StringWriter, padding string, inline bool) {
	checkError(b.WriteString(e.Ident))
	if len(e.Args) > 0 {
		checkError(b.WriteString("{"))
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
			checkError(b.WriteString("*"))
			if len(e.Args) > 0 {
				checkError(b.WriteString(", "))
			}
		}
		for k := range e.Args {
			if !inline && expand {
				checkError(b.WriteString("\n" + pad2))
			}
			checkError(b.WriteString(`"` + k + `"`))
			if i+1 < len(e.Args) {
				checkError(b.WriteString(", "))
			}
			i++
		}
		if !inline && expand {
			checkError(b.WriteString("\n" + padding))
		}
		checkError(b.WriteString("}"))
	}
	if len(e.Props) > 0 {
		checkError(b.WriteString(":"))
		for i := 0; i < len(e.Props); i++ {
			e.Props[i].format(b, padding, inline)
			if i+1 < len(e.Props) {
				checkError(b.WriteString(", "))
			}
		}
	}
}

func (e *GeometryPointExpr) format(w io.StringWriter, padding string, inline bool) {
	checkError(w.WriteString("point["))
	formatFloat(w, e.Val[0])
	checkError(w.WriteString(", "))
	formatFloat(w, e.Val[1])
	checkError(w.WriteString("]"))
	if e.Radius != nil {
		checkError(w.WriteString(":"))
		e.Radius.format(w, padding, inline)
	}
}

func (e *ArrayExpr) format(b io.StringWriter, padding string, inline bool) {
	checkError(b.WriteString("["))
	for i, expr := range e.List {
		expr.format(b, padding, inline)
		if i+1 < len(e.List) {
			checkError(b.WriteString(", "))
		}
	}
	checkError(b.WriteString("]"))
}

func (e *BinaryExpr) format(w io.StringWriter, padding string, inline bool) {
	e.Left.format(w, padding, inline)
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
		checkError(w.WriteString(" "))
	}
	if nl {
		checkError(w.WriteString("\n" + padding))
	}
	checkError(w.WriteString(op))
	if !nospace {
		checkError(w.WriteString(" "))
	}
	e.Right.format(w, padding, inline)
}

func (e *ParenExpr) format(w io.StringWriter, padding string, inline bool) {
	expand := true
	switch node := e.Expr.(type) {
	case *BinaryExpr:
		_, lok := node.Left.(*BinaryExpr)
		_, rok := node.Right.(*BinaryExpr)
		if !lok && !rok {
			expand = false
		}
	}
	pad2 := padding + padding
	checkError(w.WriteString("("))
	if expand {
		checkError(w.WriteString("\n" + pad2))
	}
	e.Expr.format(w, pad2, inline)
	if expand {
		checkError(w.WriteString("\n" + padding))
	}
	checkError(w.WriteString(")"))
}

func (e *AssignStmt) format(w io.StringWriter, padding string, inline bool) {
	checkError(w.WriteString(padding))
	e.Left.format(w, padding, inline)
	checkError(w.WriteString(" = "))
	e.Right.format(w, padding, inline)
	checkError(w.WriteString(";"))
}

func (e *IdentLit) format(w io.StringWriter, _ string, _ bool) {
	checkError(w.WriteString(e.Val))
}

func (e *WildcardLit) format(b io.StringWriter, _ string, _ bool) {
	checkError(b.WriteString("*"))
}

func (e *CalendarLit) format(b io.StringWriter, padding string, inline bool) {
	switch e.Kind {
	case DATETIME:
		checkError(b.WriteString(strconv.Itoa(e.Year)))
		checkError(b.WriteString("-"))
		checkError(b.WriteString(dt2str(int(e.Month))))
		checkError(b.WriteString("-"))
		checkError(b.WriteString(dt2str(e.Day)))
		checkError(b.WriteString("T"))
		switch e.U {
		case AM, PM:
			checkError(b.WriteString(strconv.Itoa(e.Hours) + ":" + dt2str(e.Minutes)))
			if e.Seconds > 0 {
				checkError(b.WriteString(":" + dt2str(e.Seconds)))
			}
			checkError(b.WriteString(e.U.String()))
		default:
			checkError(b.WriteString(dt2str(e.Hours) + ":" + dt2str(e.Minutes)))
			if e.Seconds > 0 {
				checkError(b.WriteString(":" + dt2str(e.Seconds)))
			}
		}
	case DATE:
		checkError(b.WriteString(strconv.Itoa(e.Year)))
		checkError(b.WriteString("-"))
		checkError(b.WriteString(dt2str(int(e.Month))))
		checkError(b.WriteString("-"))
		checkError(b.WriteString(dt2str(e.Day)))
	case TIME:
		switch e.U {
		case AM, PM:
			checkError(b.WriteString(strconv.Itoa(e.Hours) + ":" + dt2str(e.Minutes)))
			if e.Seconds > 0 {
				checkError(b.WriteString(":" + dt2str(e.Seconds)))
			}
			checkError(b.WriteString(e.U.String()))
		default:
			checkError(b.WriteString(dt2str(e.Hours) + ":" + dt2str(e.Minutes)))
			if e.Seconds > 0 {
				checkError(b.WriteString(":" + dt2str(e.Seconds)))
			}
		}
	case WEEKDAY:
		checkError(b.WriteString(shortDayNames[e.Val]))
	case MONTH:
		checkError(b.WriteString(shortMonthNames[e.Val]))
	}
}

func (e *RefLit) format(w io.StringWriter, _ string, _ bool) {
	checkError(w.WriteString("@" + e.ID))
}
