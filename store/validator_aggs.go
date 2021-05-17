package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
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

// FindValidatorEpochs returns the last N validator epochs
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

// FetchRewardsByInterval fetches reward by interval
func (s *ValidatorAggsStore) FetchRewardsByInterval(account string, from time.Time, to time.Time, timeInterval model.TimeInterval) (model.RewardsSummary, error) {
	var res model.RewardsSummary
	err := s.db.Raw(queries.ValidatorsRewards, timeInterval.String(), account, from, to, timeInterval.String()).Scan(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

// PaginateValidatorEpochs returns a paginated search of validator epochs
func (s ValidatorAggsStore) PaginateValidatorEpochs(account string, pagination Pagination) (*PaginatedResult, error) {
	if err := pagination.Validate(); err != nil {
		return nil, err
	}

	scope := s.db.
		Model(&model.ValidatorEpoch{}).
		Where("account_id = ?", account).
		Order("last_height DESC")

	var count uint
	if err := scope.Count(&count).Error; err != nil {
		return nil, err
	}

	result := []model.ValidatorEpoch{}

	err := scope.
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	paginatedResult := &PaginatedResult{
		Page:    pagination.Page,
		Limit:   pagination.Limit,
		Count:   count,
		Records: result,
	}

	return paginatedResult.update(), nil
}

// FindBy returns an validator agg record for a key and value
func (s ValidatorAggsStore) FindBy(key string, value interface{}) (*model.ValidatorAgg, error) {
	result := &model.ValidatorAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// ImportValidatorEpochs imports validator epochs records in batch
func (s ValidatorAggsStore) ImportValidatorEpochs(records []model.ValidatorEpoch) error {
	return s.bulkImport(queries.ValidatorEpochsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.AccountID,
			r.Epoch,
			r.LastHeight,
			r.LastTime,
			r.ExpectedBlocks,
			r.ProducedBlocks,
			r.Efficiency,
			r.StakingBalance,
			r.RewardFee,
			r.Reward,
		}
	})
}

// Import create validator aggregates in batch
func (s ValidatorAggsStore) Import(records []model.ValidatorAgg) error {
	t := time.Now()

	// Mark all validators as inactive
	if err := s.db.Exec("UPDATE validator_aggregates SET active = false").Error; err != nil {
		return err
	}

	return s.bulkImport(queries.ValidatorAggImport, len(records), func(i int) bulk.Row {
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
			r.Active,
			r.RewardFee,
			t,
			t,
		}
	})
}
