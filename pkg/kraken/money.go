package kraken

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Default decimal places set for integer constructors.
const DefaultDecimals = 12

// Money implements fixed-point arithmetic.
type Money struct {
	//  Unscaled integer representation.
	Integer *big.Int
	// Granularity.
	Granularity int64
	// Amount of decimal places.
	Decimals int64
}

// NewMoneyFromString creates a new [Money] object from a string.
func NewMoneyFromString(s string) (*Money, error) {
	var useBigFloat bool
	for _, l := range s {
		if (l < '0' || l > '9') && l != '.' {
			useBigFloat = true
			break
		}
	}
	if !useBigFloat {
		m := new(Money)
		m.Granularity = 1
		if s == "" {
			s = "0"
		}
		s, decimals, found := strings.Cut(s, ".")
		if found {
			s += decimals
		} else {
			decimals = ""
		}
		integer, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number \"%s\"", s)
		}
		m.Integer = integer
		m.Decimals = int64(len(decimals))
		return m, nil
	} else {
		f, success := new(big.Float).SetPrec(256).SetString(s)
		if !success {
			return nil, fmt.Errorf("invalid number \"%s\"", s)
		}
		m := NewMoneyFromBigFloat(f)
		if _, decimals, found := strings.Cut(s, "."); found {
			m = m.SetDecimals(int64(len(decimals)))
		}
		return m, nil
	}
}

// NewMoneyFromBigInt creates a new [Money] object from a [big.Int].
func NewMoneyFromBigInt(bi *big.Int) *Money {
	m := new(Money)
	m.Granularity = 1
	m.Integer = new(big.Int).Set(bi)
	return m.SetDecimals(DefaultDecimals)
}

// NewMoneyFromInt64 creates a new [Money] object from an int64.
func NewMoneyFromInt64(i int64) *Money {
	m := new(Money)
	m.Granularity = 1
	m.Integer = new(big.Int).SetInt64(i)
	return m.SetDecimals(DefaultDecimals)
}

// NewMoneyFromBigFloat creates a new [Money] object from a [big.Float].
func NewMoneyFromBigFloat(f *big.Float) *Money {
	var numDecimals int
	if _, decimals, found := strings.Cut(f.Text('f', -1), "."); found {
		numDecimals = len(decimals)
	}
	multiplicand := new(big.Float).SetFloat64(math.Pow10(numDecimals))
	m := new(Money)
	m.Granularity = 1
	m.Integer, _ = new(big.Float).Mul(f, multiplicand).Int(nil)
	m.Decimals = int64(numDecimals)
	return m
}

// NewMoneyFromFloat64 creates a new [Money] object from a float64.
func NewMoneyFromFloat64(f float64) *Money {
	return NewMoneyFromBigFloat(new(big.Float).SetFloat64(f))
}

// SetDecimals returns m with adjusted decimal places.
func (m *Money) SetDecimals(d int64) *Money {
	result := m.Copy()
	if d == m.Decimals {
		return result
	}
	diff := d - result.Decimals
	result.Decimals = d
	if result.Sign() == 0 {
		return result
	}
	absoluteDiff := int64(diff)
	if absoluteDiff < 0 {
		absoluteDiff = -absoluteDiff
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(absoluteDiff), nil)
	if diff > 0 {
		result.Integer.Mul(m.Integer, factor)
	} else {
		half := new(big.Int).Div(factor, big.NewInt(2))
		if m.Sign() < 0 {
			result.Integer.Sub(m.Integer, half)
		} else {
			result.Integer.Add(m.Integer, half)
		}
		result.Integer.Div(result.Integer, factor)
	}
	return result.RoundToGranularity()
}

// Rat returns the rational number of m.
func (m *Money) Rat() *big.Rat {
	var integer *big.Rat
	if m.Integer != nil {
		integer = new(big.Rat).SetInt(m.Integer)
	} else {
		integer = new(big.Rat)
	}
	scale := m.ScalingFactor()
	scaleRat := new(big.Rat).SetInt(scale)
	value := new(big.Rat).Quo(integer, scaleRat)
	return value
}

// Float64 returns the floating point representation of m with potential loss of precision.
func (m *Money) Float64() float64 {
	value, _ := m.Rat().Float64()
	return value
}

// Int64 returns the integer part of m with truncated decimals.
func (m *Money) Int64() int64 {
	if m.Integer == nil || m.Sign() == 0 {
		return 0
	}
	if m.Decimals == 0 {
		return m.Integer.Int64()
	}
	scale := m.ScalingFactor()
	return new(big.Int).
		Quo(m.Integer, scale).
		Int64()
}

// String returns the literal representation of m.
func (m *Money) String() string {
	return m.Rat().FloatString(int(m.Decimals))
}

// Copy creates a copy of m.
func (m *Money) Copy() *Money {
	return &Money{
		Integer:     new(big.Int).Set(m.Integer),
		Decimals:    m.Decimals,
		Granularity: m.Granularity,
	}
}

// ScalingFactor returns 10 ^ decimals in [big.Int].
func (m *Money) ScalingFactor() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(m.Decimals), nil)
}

// Add returns the result of x + y.
func (x *Money) Add(y *Money) *Money {
	result := x.Copy()
	if y.Sign() == 0 {
		return result
	}
	result.Integer.Add(x.Integer, y.SetDecimals(x.Decimals).Integer)
	return result.RoundToGranularity()
}

// Sub returns the result of x - y
func (x *Money) Sub(y *Money) *Money {
	result := x.Copy()
	if y.Sign() == 0 {
		return result
	}
	result.Integer.Sub(x.Integer, y.SetDecimals(x.Decimals).Integer)
	return result.RoundToGranularity()
}

// Mul returns the result of x * y
func (x *Money) Mul(y *Money) *Money {
	result := x.Copy()
	result.Integer.Mul(x.Integer, y.SetDecimals(x.Decimals).Integer)
	scale := x.ScalingFactor()
	half := new(big.Int).Div(scale, big.NewInt(2))
	integerSign := result.Integer.Sign()
	if integerSign < 0 {
		result.Integer.Sub(result.Integer, half)
	} else if integerSign > 0 {
		result.Integer.Add(result.Integer, half)
	}
	result.Integer.Div(result.Integer, scale)
	return result.RoundToGranularity()
}

// Div returns the result of x / y.
func (x *Money) Div(y *Money) *Money {
	result := x.Copy()
	if y.Sign() == 0 {
		panic("division by zero")
	}
	scale := x.ScalingFactor()
	result.Integer = result.Integer.Mul(result.Integer, scale)
	half := new(big.Int).Div(scale, big.NewInt(2))
	if result.Sign() < 0 {
		result.Integer.Sub(result.Integer, half)
	} else {
		result.Integer.Add(result.Integer, half)
	}
	result.Integer.Div(result.Integer, y.SetDecimals(x.Decimals).Integer)
	return result.RoundToGranularity()
}

func (x *Money) Exp(y *Money) *Money {
	result := x.Copy()
	result.Integer.Exp(x.Integer, y.SetDecimals(x.Decimals).Integer, nil)
	return result
}

// Sign returns -1 if m < 0, 0 if m == 0, and +1 if m > 0.
func (m *Money) Sign() int {
	return m.Integer.Sign()
}

// Cmp compares x to y (-1, 0, 1 for <, =, >).
func (x *Money) Cmp(y *Money) int {
	if x.Decimals == y.Decimals {
		return x.Integer.Cmp(y.Integer)
	}
	return x.Rat().Cmp(y.Rat())
}

// SmallestIncrement returns the smallest possible increment of m.
func (m *Money) SmallestIncrement() *Money {
	var amount *big.Int
	if m.Granularity == 0 {
		amount = big.NewInt(1)
	} else {
		amount = big.NewInt(m.Granularity)
	}
	smallest := new(Money)
	smallest.Integer = amount
	smallest.Decimals = m.Decimals
	smallest.Granularity = m.Granularity
	return smallest
}

// SetGranularity configures the raw integer of m to always be a multiple of t.
func (m *Money) SetGranularity(t int64) *Money {
	result := m.Copy()
	result.Granularity = t
	return result.RoundToGranularity()
}

// SetSize configures the value of m to always be a multiple of s.
func (m *Money) SetSize(s *Money) *Money {
	return m.
		SetDecimals(s.Decimals).
		SetGranularity(s.Integer.Int64())
}

// RoundToGranularity returns the rounding of m to the granularity constraint.
func (m *Money) RoundToGranularity() *Money {
	result := m.Copy()
	if m.Granularity <= 1 {
		return result
	}
	tick := big.NewInt(result.Granularity)
	remainder := new(big.Int).Mod(result.Integer, tick)
	half := new(big.Int).Div(tick, big.NewInt(2))
	remainderCmpHalf := remainder.Cmp(half)
	if remainderCmpHalf >= 0 {
		result.Integer.
			Sub(result.Integer, remainder).
			Add(result.Integer, tick)
	} else {
		result.Integer.Sub(result.Integer, remainder)
	}
	return result
}

// OffsetTicks returns the adjustment of m by an increment proportional to o.
func (m *Money) OffsetTicks(o *Money) *Money {
	return m.Add(
		m.SmallestIncrement().
			Mul(o),
	)
}

// OffsetPercent returns the adjustment of m by %o.
// Formula: m * (1 + o).
// The multiplicand decimals are set to the max precision of both m and o.
func (m *Money) OffsetPercent(o *Money) *Money {
	multiplicand := NewMoneyFromInt64(1).Add(o)
	originalDecimals := m.Decimals
	return m.SetDecimals(int64(math.Max(float64(multiplicand.Decimals), float64(originalDecimals)))).
		Mul(multiplicand).
		SetDecimals(originalDecimals)
}

// Abs returns the absolute value of m.
func (m *Money) Abs() *Money {
	result := m.Copy()
	result.Integer = new(big.Int).Abs(m.Integer)
	return result
}

// MarshalJSON implements the [json.Marshaler] interface.
func (m *Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (m *Money) UnmarshalJSON(data []byte) error {
	parsed, err := NewMoneyFromString(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*m = *parsed
	return nil
}
