package decimal

import (
	"math/big"
	"testing"
)

func TestNewFromString(t *testing.T) {
	if d, err := NewFromString("1.015"); err != nil {
		t.Error(err)
	} else if d.GetIncrement() != 1 {
		t.Errorf("d.GetIncrement() != 1, got %d", d.GetIncrement())
	} else if d.RawBigInt().Cmp(new(big.Int).SetInt64(1015)) != 0 {
		t.Errorf("d.RawBigInt() != 1015, got %s", d.integer)
	} else if d.GetScale() != 3 {
		t.Errorf("d.GetScale() != 3, got %d", d.GetScale())
	} else if d.String() != "1.015" {
		t.Errorf("d.String() != 1.015, got %s", d.String())
	}
}

func TestMath(t *testing.T) {
	if d, err := NewFromString("1.015"); err != nil {
		t.Error(err)
	} else if d = d.Add(NewFromInt64(1)); d.String() != "2.015" {
		t.Errorf("Add(1) != 2.015, got %s", d)
	} else if d = d.Sub(NewFromInt64(1)); d.String() != "1.015" {
		t.Errorf("Sub(1) != 1.015, got %s", d)
	} else if d = d.Mul(NewFromInt64(2)); d.String() != "2.030" {
		t.Errorf("Mul(2) != 2.030, got %s", d)
	} else if d = d.Div(NewFromInt64(2)); d.String() != "1.015" {
		t.Errorf("Div(2) != d.String(), got %s", d)
	} else if d = d.Pow(NewFromInt64(2)); d.String() != "1.030" {
		t.Errorf("Pow(2) != d.String(), got %s", d)
	}
}

func TestRounding(t *testing.T) {
	if d, err := NewFromString("1.002"); err != nil {
		t.Error(err)
	} else if d = d.SetIncrement(5); d.String() != "1.000" {
		t.Errorf("SetIncrement(5) != 1.000, got %s", d)
	} else if d = d.OffsetTicks(NewFromInt64(1)); d.String() != "1.005" {
		t.Errorf("OffsetTicks(1) != 1.005, got %s", d)
	} else if d = d.SetScale(2); d.String() != "1.00" {
		t.Errorf("SetScale(2) != 1.00, got %s", d)
	}
}
