package geoqlparser

import (
	"testing"
)

func TestFlatten(t *testing.T) {
	q := "trigger when tracker > 1 and tracker < 300 or tracker > 1"
	expr, err := Parse(q)
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := ToFlat(expr)
	if err != nil {
		t.Fatal(err)
	}
	// TODO:
	_ = stmt
}
