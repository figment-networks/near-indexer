package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/store/queries"
)

// ValidatorsStore handles operations on blocks
type ValidatorsStore struct {
	baseStore
}

// ByHeight returns all transactions for a given height
func (s ValidatorsStore) ByHeight(height types.Height) ([]model.Validator, error) {
	result := []model.Validator{}

	err := s.db.
		Where("height = ?", height).
		Order("id DESC").
		Find(&result).
		Error

	return result, err
}

// Import creates new validators in batch
func (s ValidatorsStore) Import(records []model.Validator) error {
	t := time.Now()

	return s.bulkImport(queries.ValidatorsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.Height,
			r.Time,
			r.AccountID,
			r.Epoch,
			r.ExpectedBlocks,
			r.ProducedBlocks,
			r.Slashed,
			r.Stake,
			r.Efficiency,
			r.RewardFee,
			t,
			t,
		}
	})
}

// Cleanup removes old validator records and keeps the N most recent ones
func (s ValidatorsStore) Cleanup(keepHeights uint64) (int64, error) {
	result := s.db.Exec(queries.ValidatorsCleanup, keepHeights)
	return result.RowsAffected, result.Error
}
