package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
)

type DelegatorInfo struct {
	Epoch           string       `json:"epoch"`
	AccountId       string       `json:"account_id"`
	ValidatorId     string       `json:"validator_id"`
	UnstakedBalance types.Amount `json:"unstaked_balance"`
	StakedBalance   types.Amount `json:"staked_balance"`
}

// Delegations constructs a set of delegation records
func Delegations(input []near.AccountInfo) ([]model.Delegation, error) {
	result := make([]model.Delegation, len(input))

	for i, d := range input {
		delegation, err := Delegation(&d)
		if err != nil {
			return nil, err
		}
		result[i] = *delegation
	}

	return result, nil
}

// Delegation constructs a new delegation record
func Delegation(input *near.AccountInfo) (*model.Delegation, error) {
	delegation := &model.Delegation{
		Account:         input.Account,
		UnstakedBalance: types.NewAmount(input.UnstakedBalance),
		StakedBalance:   types.NewAmount(input.StakedBalance),
		CanWithdraw:     input.CanWithdraw,
	}

	return delegation, delegation.Validate()
}

// Delegators constructs a set of delegator records from delegator epochs info
func Delegators(input []model.DelegatorEpoch) ([]DelegatorInfo, error) {
	result := make([]DelegatorInfo, len(input))

	for i, d := range input {
		delegation := &DelegatorInfo{
			Epoch:           d.Epoch,
			AccountId:       d.AccountID,
			ValidatorId:     d.ValidatorID,
			UnstakedBalance: d.UnstakedBalance,
			StakedBalance:   d.StakedBalance,
		}
		result[i] = *delegation
	}

	return result, nil
}
