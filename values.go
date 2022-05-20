package geoqlparser

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultResetVal  = DurVal{V: 24 * time.Hour}
	DefaultRepeatVal = Repeat{V: 1}
	DefaultRadiusVal = RadiusVal{V: 500, U: Meter}
)

type StrVal struct {
	V string
}

type DurVal struct {
	V time.Duration
}

type RadiusUnit uint

const (
	Kilometer RadiusUnit = iota
	Meter
)

type Qualifier int

const (
	All Qualifier = iota + 1
	Any
)

type RadiusVal struct {
	V uint
	U RadiusUnit
}

var errInvalidRadius = errors.New("invalid radius")

func toRadiusVal(lit string) (rv RadiusVal, err error) {
	if len(lit) == 0 {
		return rv, errInvalidRadius
	}
	if len(lit) > 12 {
		return rv, errInvalidRadius
	}
	if lit == "0" {
		return rv, nil
	}
	c := lit[0]
	if c != 'r' && c != 'R' {
		return rv, errInvalidRadius
	}
	lit = lit[1:]
	n, lit, err := scanInt(lit)
	if err != nil {
		return rv, err
	}
	switch lit {
	case "km", "KM":
		rv.U = Kilometer
	case "m", "M":
		rv.U = Meter
	default:
		return rv, errInvalidRadius
	}
	rv.V = uint(n)
	return rv, nil
}

var errBadIntVal = errors.New("invalid INT value")

func scanInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			return 0, "", errBadIntVal
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			return 0, "", errBadIntVal
		}
	}
	return x, s[i:], nil
}

func toStringVal(lit string) (StrVal, error) {
	return StrVal{V: trim(lit)}, nil
}

func trim(lit string) string {
	lit = strings.TrimLeft(lit, `"`)
	lit = strings.TrimRight(lit, `"`)
	return lit
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
