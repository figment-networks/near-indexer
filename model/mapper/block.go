package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// Block constructs the block record from the block data
func Block(block *near.Block) (*model.Block, error) {
	h := block.Header

	record := &model.Block{
		Time:            util.ParseTime(h.Timestamp),
		Height:          types.Height(h.Height),
		Producer:        block.Author,
		Hash:            h.Hash,
		PrevHash:        h.PrevHash,
		RentPaid:        types.NewAmount(h.RentPaid),
		ValidatorReward: types.NewAmount(h.ValidatorReward),
		TotalSupply:     types.NewAmount(h.TotalSupply),
		Signature:       h.Signature,
		ChunksCount:     h.ChunksIncluded,
		GasPrice:        types.NewAmount(h.GasPrice),
		GasLimit:        0, // TODO: calculate this
		GasUsed:         0, // TODO: calculate this
	}

	return record, record.Validate()
}
