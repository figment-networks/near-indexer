package store

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
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

	return s.Import(sqlValidatorsBulkInsert, len(records), func(i int) bulkRow {
		r := records[i]
		return bulkRow{
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

var (
	sqlValidatorsBulkInsert = `
		INSERT INTO validators (
			height,
			time,
			account_id,
			epoch,
			expected_blocks,
			produced_blocks,
			slashed,
			stake,
			efficiency,
			created_at,
			updated_at
		)
		VALUES @values`
)
