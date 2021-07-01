package util

import (
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
func CalculateTransactionFee(trx near.TransactionDetails) string {
	fee := new(big.Int)
	o := new(big.Int)
	r := new(big.Int)
	o.SetString(trx.TransactionOutcome.Outcome.TokensBurnt, 10)
	fee.Add(fee, o)
	for _, o := range trx.ReceiptsOutcome {
		r.SetString(o.Outcome.TokensBurnt, 10)
		fee.Add(fee, r)
	}
	return fee.String()
}
