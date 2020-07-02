package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
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
	return jsonquery.MustObject(s.db, jsonquery.Prepare(sqlValidatorDetails), id)
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
	return s.Import(sqlValidatorEpochsUpsert, len(records), func(i int) bulk.Row {
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
	return s.db.Exec(sqlValidatorCountsUpsert, height).Error
}

func (s ValidatorAggsStore) BulkUpsert(records []model.ValidatorAgg) error {
	t := time.Now()

	return s.Import(sqlValidatorAggsBulkUpsert, len(records), func(i int) bulk.Row {
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

var (
	sqlValidatorDetails = `
		SELECT
			validator_aggregates.*,
			{{ array }}
				SELECT
					validator_epochs.epoch,
					validator_epochs.last_height,
					validator_epochs.last_time,
					validator_epochs.expected_blocks,
					validator_epochs.produced_blocks,
					validator_epochs.efficiency
			{{ end_array }} AS epochs
		FROM
			validator_aggregates
		LEFT JOIN validator_epochs
			ON validator_epochs.account_id = validator_aggregates.account_id
		WHERE
			validator_aggregates.account_id = ?
		LIMIT 1`

	sqlValidatorEpochsUpsert = `
		INSERT INTO validator_epochs (
			account_id,
			epoch,
			last_height,
			last_time,
			expected_blocks,
			produced_blocks,
			efficiency
		)
		VALUES @values
		ON CONFLICT (account_id, epoch) DO UPDATE
		SET
			last_height     = excluded.last_height,
			last_time       = excluded.last_time,
			expected_blocks = excluded.expected_blocks,
			produced_blocks = excluded.produced_blocks,
			efficiency      = excluded.efficiency;
	`

	sqlValidatorCountsUpsert = `
		INSERT INTO validator_counts (
  		height,
  		time,
  		total_count,
  		active_count,
  		slashed_count
		)
		SELECT
  		blocks.height,
  		blocks.time,
  		(SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height) total_validators,
  		(SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height AND slashed = false AND efficiency > 0) active_validators,
  		(SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height AND slashed = true) slashed_validators
		FROM
			blocks
		WHERE
			blocks.height = $1
		ON CONFLICT (height) DO UPDATE
		SET
		  time          = excluded.time,
		  total_count   = excluded.total_count,
		  active_count  = excluded.active_count,
		  slashed_count = excluded.slashed_count;`

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
			expected_blocks = COALESCE((SELECT SUM(expected_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
			produced_blocks = COALESCE((SELECT SUM(produced_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
			efficiency      = COALESCE((SELECT AVG(efficiency) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
			slashed         = excluded.slashed,
			updated_at      = excluded.updated_at`
)
