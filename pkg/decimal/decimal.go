package decimal

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Decimal implements fixed-point arithmetic.
type Decimal struct {
	//  Unscaled integer representation.
	integer *big.Int
	// Smallest allowable unit for the decimal value.
	increment int64
	// Number of digits to the right of the decimal point.
	scale int64
	// Rounding function
	rounding RoundingFunction
}

// Default increment for integer constructors.
const DefaultIncrement = 1

// Default decimal points set for integer constructors.
const DefaultScale = 12

// NewFromString creates a new [Decimal] object from a string.
func NewFromString(s string) (*Decimal, error) {
	var useBigFloat bool
	for _, l := range s {
		if (l < '0' || l > '9') && l != '.' {
			useBigFloat = true
			break
		}
	}
	if !useBigFloat {
		d := new(Decimal)
		d.increment = DefaultIncrement
		d.rounding = BankersRound
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
		d.integer = integer
		d.scale = int64(len(decimals))
		return d, nil
	} else {
		f, success := new(big.Float).SetPrec(256).SetString(s)
		if !success {
			return nil, fmt.Errorf("invalid number \"%s\"", s)
		}
		d := NewFromBigFloat(f)
		if _, decimals, found := strings.Cut(s, "."); found {
			d = d.SetScale(int64(len(decimals)))
		}
		return d, nil
	}
}

// NewFromBigInt creates a new [Decimal] object from a [big.Int].
func NewFromBigInt(bi *big.Int) *Decimal {
	d := new(Decimal)
	d.increment = DefaultIncrement
	d.rounding = BankersRound
	d.integer = new(big.Int).Set(bi)
	return d.SetScale(DefaultScale)
}

// NewFromInt64 creates a new [Decimal] object from an int64.
func NewFromInt64(i int64) *Decimal {
	d := new(Decimal)
	d.increment = DefaultIncrement
	d.rounding = BankersRound
	d.integer = new(big.Int).SetInt64(i)
	return d.SetScale(DefaultScale)
}

// NewFromBigFloat creates a new [Decimal] object from a [big.Float].
func NewFromBigFloat(f *big.Float) *Decimal {
	var numDecimals int
	if _, decimals, found := strings.Cut(f.Text('f', -1), "."); found {
		numDecimals = len(decimals)
	}
	multiplicand := new(big.Float).SetFloat64(math.Pow10(numDecimals))
	d := new(Decimal)
	d.increment = DefaultIncrement
	d.rounding = BankersRound
	d.integer, _ = new(big.Float).Mul(f, multiplicand).Int(nil)
	d.scale = int64(numDecimals)
	return d
}

// NewFromFloat64 creates a new [Decimal] object from a float64.
func NewFromFloat64(f float64) *Decimal {
	return NewFromBigFloat(new(big.Float).SetFloat64(f))
}

// SetScale returns m with adjusted decimal places.
func (d *Decimal) SetScale(scale int64) *Decimal {
	result := d.Copy()
	if scale == d.scale {
		return result
	}
	diff := scale - result.scale
	result.scale = scale
	if result.Sign() == 0 {
		return result
	}
	absoluteDiff := int64(diff)
	if absoluteDiff < 0 {
		absoluteDiff = -absoluteDiff
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(absoluteDiff), nil)
	if diff > 0 {
		result.integer.Mul(d.integer, factor)
	} else {
		result.integer = result.rounding(result.integer, factor)
	}
	result.roundToGranularity()
	return result
}

// GetScale returns the number of decimal points.
func (d *Decimal) GetScale() int64 {
	return d.scale
}

// Rat returns the rational number of m.
func (d *Decimal) Rat() *big.Rat {
	var integer *big.Rat
	if d.integer != nil {
		integer = new(big.Rat).SetInt(d.integer)
	} else {
		integer = new(big.Rat)
	}
	scale := d.ScalingFactor()
	scaleRat := new(big.Rat).SetInt(scale)
	value := new(big.Rat).Quo(integer, scaleRat)
	return value
}

// Float64 returns the floating point representation of m with potential loss of precision.
func (d *Decimal) Float64() float64 {
	value, _ := d.Rat().Float64()
	return value
}

// Int64 returns the integer part of m with truncated decimals.
func (d *Decimal) Int64() int64 {
	if d.integer == nil || d.Sign() == 0 {
		return 0
	}
	if d.scale == 0 {
		return d.integer.Int64()
	}
	scale := d.ScalingFactor()
	return new(big.Int).
		Quo(d.integer, scale).
		Int64()
}

// String returns the literal representation of m.
func (d *Decimal) String() string {
	return d.Rat().FloatString(int(d.scale))
}

// Copy creates a copy of m.
func (d *Decimal) Copy() *Decimal {
	return &Decimal{
		integer:   new(big.Int).Set(d.integer),
		scale:     d.scale,
		increment: d.increment,
		rounding:  d.rounding,
	}
}

// ScalingFactor returns 10 ^ decimals in [big.Int].
func (d *Decimal) ScalingFactor() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(d.scale), nil)
}

// Add returns the result of x + y.
func (x *Decimal) Add(y *Decimal) *Decimal {
	result := x.Copy()
	if y.Sign() == 0 {
		return result
	}
	result.integer.Add(x.integer, y.SetScale(x.scale).integer)
	result.roundToGranularity()
	return result
}

// Sub returns the result of x - y
func (x *Decimal) Sub(y *Decimal) *Decimal {
	result := x.Copy()
	if y.Sign() == 0 {
		return result
	}
	result.integer.Sub(x.integer, y.SetScale(x.scale).integer)
	result.roundToGranularity()
	return result
}

// Mul returns the result of x * y
func (x *Decimal) Mul(y *Decimal) *Decimal {
	result := x.Copy()
	result.integer.Mul(x.integer, y.SetScale(x.scale).integer)
	scale := x.ScalingFactor()
	result.integer = result.rounding(result.integer, scale)
	result.roundToGranularity()
	return result
}

// Div returns the result of x / y.
func (x *Decimal) Div(y *Decimal) *Decimal {
	result := x.Copy()
	if y.Sign() == 0 {
		panic("division by zero")
	}
	scale := x.ScalingFactor()
	result.integer.Mul(result.integer, scale)
	result.integer = result.rounding(result.integer, y.SetScale(x.scale).integer)
	result.roundToGranularity()
	return result
}

// Pow returns x ** y.
func (x *Decimal) Pow(y *Decimal) *Decimal {
	return NewFromFloat64(math.Pow(x.Float64(), y.Float64())).SetScale(x.scale)
}

// Sign returns -1 if m < 0, 0 if m == 0, and +1 if m > 0.
func (d *Decimal) Sign() int {
	return d.integer.Sign()
}

// Cmp compares x to y (-1, 0, 1 for <, =, >).
func (x *Decimal) Cmp(y *Decimal) int {
	if x.scale == y.scale {
		return x.integer.Cmp(y.integer)
	}
	return x.Rat().Cmp(y.Rat())
}

// GetSmallestIncrement returns the smallest possible increment of m.
func (d *Decimal) GetSmallestIncrement() *Decimal {
	amount := big.NewInt(d.increment)
	smallest := new(Decimal)
	smallest.integer = amount
	smallest.scale = d.scale
	smallest.increment = d.increment
	smallest.rounding = d.rounding
	return smallest
}

// GetIncrement returns the smallest allowable unit for the decimal value.
func (d *Decimal) GetIncrement() int64 {
	return d.increment
}

// SetIncrement sets the smallest allowable unit of d.
func (d *Decimal) SetIncrement(increment int64) *Decimal {
	result := d.Copy()
	result.increment = increment
	result.roundToGranularity()
	return result
}

func (d *Decimal) SetRounding(rounding RoundingFunction) *Decimal {
	result := d.Copy()
	result.rounding = rounding
	return result
}

// SetSize ensures the value of d is always a multiple of the specified.
func (d *Decimal) SetSize(size *Decimal) *Decimal {
	return d.
		SetScale(size.scale).
		SetIncrement(size.integer.Int64())
}

// roundToGranularity returns the rounding of m to the granularity constraint.
func (d *Decimal) roundToGranularity() {
	if d.increment <= 1 {
		return
	}
	tick := big.NewInt(d.increment)
	remainder := new(big.Int).Mod(d.integer, tick)
	half := new(big.Int).Div(tick, big.NewInt(2))
	remainderCmpHalf := remainder.Cmp(half)
	if remainderCmpHalf > 0 {
		d.integer.
			Sub(d.integer, remainder).
			Add(d.integer, tick)
	} else if remainderCmpHalf < 0 {
		d.integer.Sub(d.integer, remainder)
	} else {
		roundedDown := new(big.Int).Sub(d.integer, remainder)
		quotient := new(big.Int).Div(roundedDown, tick)
		if quotient.Bit(0) == 0 {
			d.integer = roundedDown
		} else {
			d.integer = roundedDown.Add(roundedDown, tick)
		}
	}
}

// OffsetTicks returns the adjustment of m by an increment proportional to o.
func (d *Decimal) OffsetTicks(o *Decimal) *Decimal {
	return d.Add(
		d.GetSmallestIncrement().
			Mul(o),
	)
}

// OffsetPercent returns the adjustment of m by %o.
// Formula: m * (1 + o).
// The multiplicand decimals are set to the max precision of both m and o.
func (d *Decimal) OffsetPercent(o *Decimal) *Decimal {
	multiplicand := NewFromInt64(1).Add(o)
	originalDecimals := d.scale
	return d.SetScale(int64(math.Max(float64(multiplicand.scale), float64(originalDecimals)))).
		Mul(multiplicand).
		SetScale(originalDecimals)
}

// Abs returns the absolute value of m.
func (d *Decimal) Abs() *Decimal {
	result := d.Copy()
	result.integer = new(big.Int).Abs(d.integer)
	return result
}

// RawBigInt returns the raw integer representation in the form of [big.Int].
func (d *Decimal) RawBigInt() *big.Int {
	return new(big.Int).Set(d.integer)
}

// MarshalJSON implements the [json.Marshaler] interface.
func (d *Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	parsed, err := NewFromString(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*d = *parsed
	return nil
}
