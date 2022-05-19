package geoqlparser

import (
	"strconv"
	"strings"
	"time"
)

var (
	DefaultResetVal  = DurVal{V: 24 * time.Hour}
	DefaultRepeatVal = Repeat{V: 1}
)

type StrVal struct {
	V string
}

type DurVal struct {
	V time.Duration
}

func toStringVal(lit string) (StrVal, error) {
	lit = strings.TrimLeft(lit, `"`)
	lit = strings.TrimRight(lit, `"`)
	return StrVal{V: lit}, nil
}

type IntVal struct {
	V int
}

func toIntVal(lit string) (IntVal, error) {
	val, err := strconv.Atoi(lit)
	if err != nil {
		return IntVal{}, err
	}
	return IntVal{V: val}, nil
}

type FloatVal struct {
	V float64
}

func toFloatVal(lit string) (FloatVal, error) {
	val, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		return FloatVal{}, err
	}
	return FloatVal{V: val}, nil
}

type ArrayFloatVal struct {
	V []float64
}

type ArrayIntVal struct {
	V []int
}

type ArrayStringVal struct {
	V []string
}

type ListFloatVal struct {
	V map[float64]struct{}
}

type ListIntVal struct {
	V map[int]struct{}
}

type ListStringVal struct {
	V map[string]struct{}
}

type Variable struct {
	Ident string
	Val   interface{}
}
