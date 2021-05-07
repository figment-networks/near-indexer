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

type HeightRange struct {
	Start uint64
	End   uint64
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (t TimeRange) FullRange(bucket string) (time.Time, time.Time, error) {
	now := time.Now()

	fromStart, _, err := getTimeRange(bucket, t.Start)
	if err != nil {
		return now, now, err
	}

	_, toEnd, err := getTimeRange(bucket, t.End)
	if err != nil {
		return now, now, err
	}

	return fromStart, toEnd, nil
}

// CreateBlockStats populates block stats for given block height range
func (s StatsStore) CreateBlockStats(bucket string, timeRange TimeRange) error {
	timeFrom, timeTo, err := timeRange.FullRange(bucket)
	if err != nil {
		return err
	}
	rr:= "INSERT INTO block_stats (time,bucket, blocks_count,block_time_avg,  validators_count,  transactions_count ) SELECT DATE_TRUNC('@bucket', time) AS time, '@bucket' AS bucket, COUNT(1) AS blocks_count, ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2) AS block_time_avg,  COALESCE(MAX(approvals_count), 0) AS validators_count, COALESCE(SUM(transactions_count), 0) AS transactions_count FROM blocks WHERE  time >= $1::timestamp AND time <= $2::timestamp GROUP BY DATE_TRUNC('@bucket', time) ON CONFLICT (time, bucket) DO UPDATE SET blocks_count       = excluded.blocks_count, block_time_avg     = excluded.block_time_avg,  validators_count   = excluded.validators_count,  transactions_count = excluded.transactions_count"
	query := s.prepareBucket(rr, bucket)
	return s.db.Exec(query, timeFrom, timeTo).Error
}

// CreateValidatorsStats populates validators stats for a time bucket
func (s StatsStore) CreateValidatorsStats(bucket string, ts time.Time) error {
	start, end, err := getTimeRange(bucket, ts)
	if err != nil {
		return err
	}
	rr:= "INSERT INTO validator_stats (  time,  bucket,  total_min, total_max, total_avg,  active_min, active_max, active_avg,  slashed_min, slashed_max, slashed_avg) SELECT  DATE_TRUNC('@bucket', time) AS time,  '@bucket' AS bucket, MIN(total_count),   MAX(total_count),   ROUND(AVG(total_count), 2),\n  MIN(active_count),  MAX(active_count),  ROUND(AVG(active_count), 2),  MIN(slashed_count), MAX(slashed_count), ROUND(AVG(slashed_count), 2) FROM  validator_counts WHERE  time >= $1 AND time <= $2 GROUP BY  DATE_TRUNC('@bucket', time) ON CONFLICT (time, bucket) DO UPDATE SET  total_min   = excluded.total_min,  total_max   = excluded.total_max,\n  total_avg   = excluded.total_avg,  active_min  = excluded.active_min,  active_max  = excluded.active_max,  active_avg  = excluded.active_avg,  slashed_min = excluded.slashed_min,  slashed_max = excluded.slashed_max,  slashed_avg = excluded.slashed_avg"
	query := s.prepareBucket(rr, bucket)

	return s.db.Exec(query, start, end).Error
}

// getTimeRange returns the start/end time for a given time bucket
func getTimeRange(bucket string, ts time.Time) (start time.Time, end time.Time, err error) {
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
