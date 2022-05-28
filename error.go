package geoqlparser

import (
	"errors"
	"fmt"
)

var errSkipExpr = errors.New("skip expression")

type Error struct {
	Line   int
	Column int
	Offset int
	Msg    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("syntax error at position line=%d, column=%d, offset=%d %s",
		e.Line, e.Column, e.Offset, e.Msg)
}

func newError(t *Tokenizer, msg string) *Error {
	return &Error{
		Line:   t.scanner.Position.Line,
		Column: t.scanner.Position.Column,
		Offset: t.scanner.Position.Offset,
		Msg:    msg,
	}
}
