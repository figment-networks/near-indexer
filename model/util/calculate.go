package util

import (
	"errors"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
	"math/big"
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

// Divide divides for given inputs x/y
func Divide(x int, y int) (*big.Int, error) {
	if y == 0 {
		return nil, errors.New("can not be divided")
	}
	n := big.NewInt(int64(x))
	d := big.NewInt(int64(y))
	return n.Div(n, d), nil
}

// CalculateReward calculates reward of the given validator
func CalculateReward(validator *model.Validator, rewardFeeFraction near.RewardFee) (types.Amount, error) {
	rew := new(big.Int)
	reward, _ := rew.SetString(validator.Stake.String(), 10)
	reward.Mul(reward, big.NewInt(int64(rewardFeeFraction.Numerator)))
	reward.Quo(reward, big.NewInt(int64(rewardFeeFraction.Denominator)))
	return types.NewAmount(reward.String()), nil
}
