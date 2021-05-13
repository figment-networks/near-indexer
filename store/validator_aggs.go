package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
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
	rr:= "INSERT INTO validator_epochs (  account_id,  epoch,  last_height,  last_time,  expected_blocks,  produced_blocks,  efficiency,  staking_balance,  reward_fee, reward_fee_fraction) VALUES @values ON CONFLICT (account_id, epoch) DO UPDATE SET  last_height     = excluded.last_height,  last_time       = excluded.last_time,  expected_blocks = excluded.expected_blocks, produced_blocks = excluded.produced_blocks,  efficiency      = ROUND(excluded.efficiency, 4), staking_balance = excluded.staking_balance, reward_fee      = COALESCE(excluded.reward_fee, validator_epochs.reward_fee), reward_fee_fraction = excluded.reward_fee_fraction "
	return s.bulkImport(rr, len(records), func(i int) bulk.Row {
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
			r.RewardFeeFraction,g
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
	rr := "INSERT INTO validator_aggregates( start_height, start_time,  last_height,  last_time,  account_id,  expected_blocks,  produced_blocks,  slashed,  stake,  efficiency,  active,  reward_fee,  created_at,  updated_at) VALUES @values ON CONFLICT(account_id) DO UPDATE SET last_height     = excluded.last_height, last_time       = excluded.last_time,  expected_blocks = COALESCE((SELECT SUM(expected_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),  produced_blocks = COALESCE((SELECT SUM(produced_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),  efficiency      = COALESCE((SELECT AVG(efficiency) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),  stake           = excluded.stake,  slashed         = excluded.slashed,  active          = excluded.active,  reward_fee      = COALESCE(excluded.reward_fee, validator_aggregates.reward_fee), updated_at      = excluded.updated_at"
	return s.bulkImport(rr, len(records), func(i int) bulk.Row {
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
