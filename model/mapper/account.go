package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// AccountFromValidator constructs an account from the validator data
func AccountFromValidator(block *near.Block, input *near.Validator) (*model.Account, error) {
	height := types.Height(block.Header.Height)
	time := util.ParseTime(block.Header.Timestamp)

	acc := &model.Account{
		Name:           input.AccountID,
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
		StakingBalance: types.NewAmount(input.Stake),
	}

	return acc, acc.Validate()
}
