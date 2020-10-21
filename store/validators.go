package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
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

// BulkInsert creates new validators or updates existing ones
func (s ValidatorsStore) BulkInsert(records []model.Validator) error {
	t := time.Now()

	return s.Import(queries.ValidatorsImport, len(records), func(i int) bulk.Row {
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
			t,
			t,
		}
	})
}

// CountsForInterval returns validator counts for a period of time
func (s ValidatorsStore) CountsForInterval(interval, period string) ([]byte, error) {
	return jsonquery.MustArray(s.db, queries.ValidatorsCountsForInterval, interval, period)
}

// Cleanup removes any records before a certain height
func (s ValidatorsStore) Cleanup(maxHeight uint64) error {
	return s.db.Delete(s.model, "height < ?", maxHeight).Error
}

func (s ValidatorsStore) CleanupCounts() error {
	return s.db.Exec(queries.ValidatorsPurgeCounts).Error
}
