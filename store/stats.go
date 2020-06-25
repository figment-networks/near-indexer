package store

import (
	"errors"
	"strings"
	"time"

	"github.com/figment-networks/near-indexer/model/util"
)

const (
	BucketHour = "h"
	BucketDay  = "d"
)

type StatsStore struct {
	baseStore
}

// CreateBlockStats populates block stats for a time bucket
func (s StatsStore) CreateBlockStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}
	query := s.prepareBucket(sqlCreateBlockStats, bucket)

	return s.db.Exec(query, start, end).Error
}

// CreateValidatorsStats populates validators stats for a time bucket
func (s StatsStore) CreateValidatorsStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}
	query := s.prepareBucket(sqlCreateValidatorsStats, bucket)

	return s.db.Exec(query, start, end).Error
}

// getTimeRange returns the start/end time for a given time bucket
func (s StatsStore) getTimeRange(bucket string, ts time.Time) (start time.Time, end time.Time, err error) {
	switch bucket {
	case BucketHour:
		start, end = util.HourInterval(ts)
	case BucketDay:
		start, end = util.DayInterval(ts)
	default:
		err = errors.New("invalid time bucket")
	}
	return
}

// prepareBucket replaces references of time bucket in the query
func (s StatsStore) prepareBucket(q, bucket string) string {
	return strings.ReplaceAll(q, "@bucket", bucket)
}

var (
	sqlCreateBlockStats = `
		INSERT INTO block_stats (
			time,
			bucket,
			blocks_count,
			block_time_avg
		)
		SELECT
			DATE_TRUNC('@bucket', time) AS time,
			'@bucket' AS bucket,
			COUNT(1) AS blocks_count,
			ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2) AS block_time_avg
		FROM
			blocks
		WHERE
			time >= $1 AND time <= $2
		GROUP BY
			DATE_TRUNC('@bucket', time)
		ON CONFLICT (time, bucket) DO UPDATE
		SET
			blocks_count   = excluded.blocks_count,
			block_time_avg = excluded.block_time_avg`

	sqlCreateValidatorsStats = `
		INSERT INTO validator_stats (
			time,
			bucket,
			total_min, total_max, total_avg,
			active_min, active_max, active_avg,
			slashed_min, slashed_max, slashed_avg
		)
		SELECT
			DATE_TRUNC('@bucket', time) AS time,
			'@bucket' AS bucket,
			MIN(total_count),   MAX(total_count),   ROUND(AVG(total_count), 2),
			MIN(active_count),  MAX(active_count),  ROUND(AVG(active_count), 2),
			MIN(slashed_count), MAX(slashed_count), ROUND(AVG(slashed_count), 2)
		FROM
			validator_counts
		WHERE
			time >= $1 AND time <= $2
		GROUP BY
			DATE_TRUNC('@bucket', time)
		ON CONFLICT (time, bucket) DO UPDATE
		SET
			total_min   = excluded.total_min,
			total_max   = excluded.total_max,
			total_avg   = excluded.total_avg,
			active_min  = excluded.active_min,
			active_max  = excluded.active_max,
			active_avg  = excluded.active_avg,
			slashed_min = excluded.slashed_min,
			slashed_max = excluded.slashed_max,
			slashed_avg = excluded.slashed_avg`
)
