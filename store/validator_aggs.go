package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
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
		Order("efficiency DESC, produced_blocks DESC, last_height DESC").
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s ValidatorAggsStore) FindValidatorEpochs(account string, limit int) ([]model.ValidatorEpoch, error) {
	result := []model.ValidatorEpoch{}

	err := s.db.
		Model(&model.ValidatorEpoch{}).
		Where("account_id = ?", account).Order("last_height DESC").
		Limit(limit).
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

func (s ValidatorAggsStore) FindDetails(id string) ([]byte, error) {
	return jsonquery.MustObject(s.db, jsonquery.Prepare(queries.ValidatorAggDetails), id)
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

func (s ValidatorAggsStore) ImportValidatorEpochs(records []model.ValidatorEpoch) error {
	return s.Import(queries.ValidatorEpochsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.AccountID,
			r.Epoch,
			r.LastHeight,
			r.LastTime,
			r.ExpectedBlocks,
			r.ProducedBlocks,
			r.Efficiency,
		}
	})
}

// UpdateCountsForHeight creates a count tracking record for a given height
func (s ValidatorAggsStore) UpdateCountsForHeight(height uint64) error {
	return s.db.Exec(queries.ValidatorCountsImport, height).Error
}

func (s ValidatorAggsStore) BulkUpsert(records []model.ValidatorAgg) error {
	t := time.Now()

	return s.Import(queries.ValidatorAggImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
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
