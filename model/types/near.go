package types

import (
	"math/big"
)

const (
	// exponent for calculating how many indivisible units are there in one NEAR
	nearNominationExp = 24

	// formatting string for the fraction part of the near amount
	nearFractionFormat = "%024s"

	// default number of fraction digits for the formatter
	nearFractionDigits = 5
)

var (
	// zero value
	nearZero = bigint("0")

	// number of indivisible units in one NEAR
	nearNomination = bigint("10").Exp(bigint("10"), bigint("10"), nil)

	// pre-calculated offests used for rounding to different number of digits
	nearRoundingOffsets []big.Int
)

func init() {
	nearRoundingOffsets = make([]big.Int, nearNominationExp)
	bn10 := bigint("10")
	offset := bigint("5")

	for i := 0; i < nearNominationExp; i++ {
		nearRoundingOffsets[i] = *offset
		offset.Mul(offset, bn10)
	}
}

func bigint(src string) *big.Int {
	var v big.Int
	v.SetString(src, 10)
	return &v
}
