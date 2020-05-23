package store

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
)

type ValidatorAggsStore struct {
	baseStore
}

// All returns all validator records
func (s ValidatorAggsStore) All() ([]model.ValidatorAgg, error) {
	result := []model.ValidatorAgg{}

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}

// Top returns top validators
func (s ValidatorAggsStore) Top() ([]model.ValidatorAgg, error) {
	result := []model.ValidatorAgg{}

	err := s.db.
		Order("efficiency DESC").
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindBy returns an validator agg record for a key and value
func (s ValidatorAggsStore) FindBy(key string, value interface{}) (*model.ValidatorAgg, error) {
	result := &model.ValidatorAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// Upsert creates or updates and existing agg record
func (s ValidatorAggsStore) Upsert(record *model.ValidatorAgg) error {
	agg, err := s.FindBy("account_id", record.AccountID)
	if err != nil {
		if isNotFound(err) {
			return s.Create(record)
		}
		return err
	}

	// Got an older record, should skip any updates
	if agg.LastHeight > record.LastHeight {
		return nil
	}

	agg.LastHeight = record.LastHeight
	agg.LastTime = record.LastTime
	agg.ExpectedBlocks = record.ExpectedBlocks
	agg.ProducedBlocks = record.ProducedBlocks
	agg.Slashed = record.Slashed
	agg.Stake = record.Stake
	agg.Efficiency = record.Efficiency

	return s.Update(agg)
}

func (s ValidatorAggsStore) BulkUpsert(records []model.ValidatorAgg) error {
	t := time.Now()

	return s.Import(sqlValidatorAggsBulkUpsert, len(records), func(i int) bulkRow {
		r := records[i]
		return bulkRow{
			r.StartHeight,
			r.StartTime,
			r.LastHeight,
			r.LastTime,
			r.AccountID,
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
	sqlValidatorAggsBulkUpsert = `
		INSERT INTO validator_aggregates(
			start_height,
			start_time,
			last_height,
			last_time,
			account_id,
			expected_blocks,
			produced_blocks,
			slashed,
			stake,
			efficiency,
			created_at,
			updated_at
		) 
		VALUES @values
		ON CONFLICT(account_id) DO UPDATE
		SET
			last_height     = excluded.last_height,
			last_time       = excluded.last_time,
			expected_blocks = excluded.expected_blocks,
			produced_blocks = excluded.produced_blocks,
			slashed         = excluded.slashed,
			efficiency      = excluded.efficiency,
			updated_at      = excluded.updated_at`
)
