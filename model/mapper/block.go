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
		ID:             types.Height(h.Height),
		Hash:           h.Hash,
		Time:           util.ParseTime(h.Timestamp),
		Producer:       block.Author,
		TotalSupply:    types.NewAmount(h.TotalSupply),
		Epoch:          h.EpochID,
		ChunksCount:    h.ChunksIncluded,
		ApprovalsCount: len(block.Header.Approvals),
		GasPrice:       types.NewAmount(h.GasPrice),
	}

	return record, record.Validate()
}
