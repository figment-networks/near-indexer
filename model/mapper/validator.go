package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// Validator constructs a new validator record from chain input
func Validator(block *near.Block, v *near.Validator) (*model.Validator, error) {
	result := &model.Validator{
		Height:         types.Height(block.Header.Height),
		Time:           util.ParseTime(block.Header.Timestamp),
		AccountID:      v.AccountID,
		Epoch:          block.Header.EpochID,
		ExpectedBlocks: v.NumExpectedBlocks,
		ProducedBlocks: v.NumProducedBlocks,
		Stake:          types.NewAmount(v.Stake),
		Efficiency:     util.Percentage(v.NumExpectedBlocks, v.NumProducedBlocks),
	}

	return result, result.Validate()
}

// ValidatorAgg constructs a new validator record from chain input
func ValidatorAgg(block *near.Block, v *near.Validator) (*model.ValidatorAgg, error) {
	height := types.Height(block.Header.Height)
	time := util.ParseTime(block.Header.Timestamp)

	result := &model.ValidatorAgg{
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
		AccountID:      v.AccountID,
		ExpectedBlocks: v.NumExpectedBlocks,
		ProducedBlocks: v.NumProducedBlocks,
		Stake:          types.NewAmount(v.Stake),
		Efficiency:     util.Percentage(v.NumExpectedBlocks, v.NumProducedBlocks),
		Active:         true,
		Slashed:        v.IsSlashed,
	}

	return result, nil
}
