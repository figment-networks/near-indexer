package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
)

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
		UnstakedBalance: types.NewAmount(input.StakedBalance),
		StakedBalance:   types.NewAmount(input.StakedBalance),
		CanWithdraw:     input.CanWithdraw,
	}

	return delegation, delegation.Validate()
}
