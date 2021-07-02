package util

import (
	"errors"
	"math/big"

	"github.com/figment-networks/near-indexer/near"
)

// Percentage returns a percentage value for given inputs
func Percentage(max int, cur int) float64 {
	if max < 0 && cur < 0 {
		return 0
	}
	if max == 0 {
		return 0
	}

	return (float64(cur) * 100.0) / float64(max)
}

// CalculateTransactionFee calculates transaction fee
func CalculateTransactionFee(trx near.TransactionDetails) (string, error) {
	fee := new(big.Int)
	o := new(big.Int)
	r := new(big.Int)
	_, ok := o.SetString(trx.TransactionOutcome.Outcome.TokensBurnt, 10)
	if !ok {
		return "", errors.New("error with tokens burnt field")
	}
	fee.Add(fee, o)
	for _, o := range trx.ReceiptsOutcome {
		_, ok := r.SetString(o.Outcome.TokensBurnt, 10)
		if !ok {
			return "", errors.New("error with tokens burnt field")
		}
		fee.Add(fee, r)
	}
	return fee.String(), nil
}
