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

// CountsForInterval returns validator counts for a period of time
func (s ValidatorsStore) CountsForInterval(interval, period string) ([]model.ValidatorIntervalStat, error) {
	rows, err := s.db.Raw(sqlValidatorCountsForInterval, interval, period).Rows()
	if err != nil {
		return nil, checkErr(err)
	}
	defer rows.Close()

	result := []model.ValidatorIntervalStat{}

	for rows.Next() {
		row := model.ValidatorIntervalStat{}
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, err
}

var (
	sqlValidatorCountsForInterval = `
		SELECT
  		time_bucket($1, time) AS time_interval,
  		ROUND(AVG(active_count)) AS count
		FROM
  		validator_counts
		WHERE
  		(
    		SELECT time
    		FROM validator_counts
    		ORDER BY time DESC
    		LIMIT 1
  		) - $2::INTERVAL < time
		GROUP BY time_interval
		ORDER BY time_interval ASC;`

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
