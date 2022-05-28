package geoqlparser

import (
	"log"
	"testing"
)

func TestWalkBinaryExpr(t *testing.T) {
	q := "trigger when tracker > 1 and tracker < 300 or tracker > 1"
	expr, err := Parse(q)
	if err != nil {
		t.Fatal(err)
	}
	expr, err = ToFlat(expr)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := ApplyBinaryExpr(expr, func(left Expr, right Expr, op Token) (bool, error) {
		log.Printf("%T %T %s", left, right, KeywordString(op))
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// TODO:
	_ = ok
}

func BenchmarkApplyBinaryExpr(b *testing.B) {
	q := "trigger when tracker > 1 and tracker < 300 or tracker > 1"
	expr, err := Parse(q)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := ApplyBinaryExpr(expr, func(left Expr, right Expr, op Token) (bool, error) {
			return true, nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkApplyFlatBinaryExpr(b *testing.B) {
	q := "trigger when tracker > 1 and tracker < 300 or tracker > 1"
	expr, err := Parse(q)
	if err != nil {
		b.Fatal(err)
	}
	expr, err = ToFlat(expr)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := ApplyBinaryExpr(expr, func(left Expr, right Expr, op Token) (bool, error) {
			return true, nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
