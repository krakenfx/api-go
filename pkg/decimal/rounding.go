package decimal

import "math/big"

type RoundingFunction func(value *big.Int, scale *big.Int) *big.Int

func BankersRound(value *big.Int, scale *big.Int) *big.Int {
	quotient, remainder := new(big.Int).QuoRem(value, scale, new(big.Int))
	half := new(big.Int).Div(scale, big.NewInt(2))
	cmp := remainder.Cmp(half)
	switch {
	case cmp > 0:
		// remainder > half: round up
		quotient.Add(quotient, big.NewInt(1))
	case cmp < 0:
		// remainder < half: round down (already done)
	default:
		// remainder == half: round to even
		if quotient.Bit(0) == 1 { // odd
			quotient.Add(quotient, big.NewInt(1))
		}
	}
	return quotient
}
