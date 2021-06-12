package util

import (
	"testing"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"

	"github.com/stretchr/testify/assert"
)

func TestCalculateValidatorReward(t *testing.T) {
	type args struct {
		rewards           types.Amount
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
				rewards: types.NewInt64Amount(10000),
				rewardFeeFraction: near.RewardFee{
					Numerator:   10,
					Denominator: 100,
				},
			},
			result: types.NewInt64Amount(1000),
		},
		{
			name: "error case denominator",
			args: args{
				rewards: types.NewInt64Amount(10000),
				rewardFeeFraction: near.RewardFee{
					Numerator:   10,
					Denominator: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateValidatorReward(tt.args.rewards, tt.args.rewardFeeFraction)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res, tt.result)
			}
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
		{
			name: "error case staked balance",
			args: args{
				delegation: near.Delegation{},
				validator: &model.Validator{
					Stake: types.NewInt64Amount(10000),
				},
				remainingRewards: types.NewInt64Amount(9000),
			},
			result:  types.NewInt64Amount(1800),
			wantErr: true,
		},
		{
			name: "error case validator stake",
			args: args{
				delegation: near.Delegation{
					StakedBalance: "2000",
				},
				validator:        &model.Validator{},
				remainingRewards: types.NewInt64Amount(9000),
			},
			result:  types.NewInt64Amount(1800),
			wantErr: true,
		},
		{
			name: "error case remaining rewards",
			args: args{
				delegation: near.Delegation{
					StakedBalance: "2000",
				},
				validator: &model.Validator{
					Stake: types.NewInt64Amount(10000),
				},
			},
			result:  types.NewInt64Amount(1800),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateDelegatorReward(tt.args.delegation, tt.args.validator, tt.args.remainingRewards)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res, tt.result)
			}
		})
	}
}
