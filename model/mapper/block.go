package mapper

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

// Block constructs the block record from the chain inpu
func Block(block *near.Block) (*model.Block, error) {
	h := block.Header

	record := &model.Block{
		Time:            time.Unix(0, h.Timestamp),
		Height:          h.Height,
		Producer:        block.Author,
		Hash:            h.Hash,
		PrevHash:        h.PrevHash,
		GasPrice:        h.GasPrice,
		RentPaid:        h.RentPaid,
		ValidatorReward: h.ValidatorReward,
		TotalSupply:     h.TotalSupply,
		Signature:       h.Signature,
		ChunksCount:     h.ChunksIncluded,
	}

	return record, record.Validate()
}
