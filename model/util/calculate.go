package util

import (
	"errors"
	"math/big"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
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

// Divide divides for given inputs x/y
func Divide(x int, y int) (*big.Int, error) {
	if y == 0 {
		return nil, errors.New("can not be divided")
	}
	n := big.NewInt(int64(x))
	d := big.NewInt(int64(y))
	return n.Div(n, d), nil
}

// CalculateValidatorReward calculates reward of the given validator
func CalculateValidatorReward(validator *model.Validator, rewardFeeFraction near.RewardFee) (types.Amount, error) {
	reward, _ := new(big.Int).SetString(validator.Stake.String(), 10)
	reward.Mul(reward, big.NewInt(int64(rewardFeeFraction.Numerator)))
	reward.Div(reward, big.NewInt(int64(rewardFeeFraction.Denominator)))
	return types.NewAmount(reward.String()), nil
}

// CalculateDelegatorReward calculates reward of the given delegation of the given validator
func CalculateDelegatorReward(delegation near.Delegation, validator *model.Validator, remainingRewards types.Amount) (types.Amount, error) {
	reward, _ := new(big.Int).SetString(remainingRewards.String(), 10)
	totalStakedBalance, _ := new(big.Int).SetString(validator.Stake.String(), 10)
	stakedBalance, _ := new(big.Int).SetString(delegation.StakedBalance, 10)
	reward.Mul(reward, stakedBalance)
	reward.Div(reward, totalStakedBalance)
	return types.NewAmount(reward.String()), nil
}
