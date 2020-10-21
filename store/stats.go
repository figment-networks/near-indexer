package store

import (
	"errors"
	"strings"
	"time"

	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/store/queries"
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
	query := s.prepareBucket(queries.StatsCreateBlocks, bucket)

	return s.db.Exec(query, start, end).Error
}

// CreateValidatorsStats populates validators stats for a time bucket
func (s StatsStore) CreateValidatorsStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}
	query := s.prepareBucket(queries.StatsCreateValidators, bucket)

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
