package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
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

// Import creates new validators in batch
func (s ValidatorsStore) Import(records []model.Validator) error {
	t := time.Now()

	rr := "INSERT INTO validators (  height,  time,  account_id,  epoch,  expected_blocks,  produced_blocks,  slashed,  stake,  efficiency,  reward_fee,  created_at,  updated_at) VALUES @values"
	return s.bulkImport(rr, len(records), func(i int) bulk.Row {
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
	rr := "WITH recent_heights AS (  SELECT DISTINCT height FROM validators ORDER BY height DESC LIMIT ? ) DELETE FROM validators WHERE height < (SELECT COALESCE(MIN(height), 0) FROM recent_heights)"
	result := s.db.Exec(rr, keepHeights)
	return result.RowsAffected, result.Error
}
