package mapper

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

// Validator constructs a new validator record from chain input
func Validator(block *near.Block, v *near.Validator) (*model.Validator, error) {
	efficiency := float32(0)
	if v.NumExpectedBlocks > 0 && v.NumProducedBlocks > 0 {
		efficiency = (float32(v.NumProducedBlocks) * 100.0) / float32(v.NumExpectedBlocks)
	}

	result := &model.Validator{
		Height:         block.Header.Height,
		Time:           time.Unix(0, block.Header.Timestamp),
		PublicKey:      v.PublicKey,
		AccountID:      v.AccountID,
		ExpectedBlocks: v.NumExpectedBlocks,
		ProducedBlocks: v.NumProducedBlocks,
		Stake:          v.Stake,
		Efficiency:     efficiency,
	}

	return result, result.Validate()
}
