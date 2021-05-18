package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
)

func TestCalculateValidatorReward(t *testing.T) {
	type args struct {
		validator         *model.Validator
		rewardFeeFraction near.RewardFee
	}
	tests := []struct {
		name    string
		args    args
		result  types.Amount
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				validator: &model.Validator{
					Stake: types.NewInt64Amount(10000),
				},
				rewardFeeFraction: near.RewardFee{
					Numerator:   10,
					Denominator: 100,
				},
			},
			result: types.NewInt64Amount(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := CalculateValidatorReward(tt.args.validator, tt.args.rewardFeeFraction)
			assert.Equal(t, res, tt.result)
		})
	}
}

func TestCalculateDelegatorReward(t *testing.T) {
	type args struct {
		delegation       near.Delegation
		validator        *model.Validator
		remainingRewards types.Amount
	}
	tests := []struct {
		name    string
		args    args
		result  types.Amount
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				delegation: near.Delegation{
					StakedBalance: "2000",
				},
				validator: &model.Validator{
					Stake: types.NewInt64Amount(10000),
				},
				remainingRewards: types.NewInt64Amount(9000),
			},
			result: types.NewInt64Amount(1800),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := CalculateDelegatorReward(tt.args.delegation, tt.args.validator, tt.args.remainingRewards)
			assert.Equal(t, res, tt.result)
		})
	}
}
