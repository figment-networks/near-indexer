package store

import (
	"github.com/figment-networks/near-indexer/model"
)

// BlocksStore handles operations on blocks
type BlocksStore struct {
	baseStore
}

// CreateIfNotExists creates the block if it does not exist
func (s BlocksStore) CreateIfNotExists(block *model.Block) error {
	_, err := s.FindByHash(block.Hash)
	if isNotFound(err) {
		return s.Create(block)
	}
	return nil
}

// FindBy returns a block for a matching attribute
func (s BlocksStore) FindBy(key string, value interface{}) (*model.Block, error) {
	result := &model.Block{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a block with matching ID
func (s BlocksStore) FindByID(id int64) (*model.Block, error) {
	return s.FindBy("id", id)
}

// FindByHash returns a block with the matching hash
func (s BlocksStore) FindByHash(hash string) (*model.Block, error) {
	return s.FindBy("hash", hash)
}

// FindByHeight returns a block with the matching height
func (s BlocksStore) FindByHeight(height uint64) (*model.Block, error) {
	return s.FindBy("height", height)
}

// Recent returns the most recent block
func (s BlocksStore) Recent() (*model.Block, error) {
	block := &model.Block{}

	err := s.db.
		Order("height DESC").
		First(block).
		Error

	return block, checkErr(err)
}

// Search returns matching blocks
func (s BlocksStore) Search() ([]model.Block, error) {
	result := []model.Block{}

	err := s.db.
		Order("height DESC").
		Limit(25).
		Find(&result).
		Error

	return result, err
}

// AvgRecentTimes returns recent blocks averages
func (s BlocksStore) AvgRecentTimes(limit int64) (*model.BlockAvgStat, error) {
	res := &model.BlockAvgStat{}

	err := s.db.
		Raw(blockTimesForRecentBlocksQuery, limit).
		Scan(res).
		Error

	return res, checkErr(err)
}

// AvgTimesForInterval returns block stats for a given interval
func (s BlocksStore) AvgTimesForInterval(interval, period string) ([]model.BlockIntervalStat, error) {
	rows, err := s.db.Raw(blockTimesForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, checkErr(err)
	}
	defer rows.Close()

	result := []model.BlockIntervalStat{}

	for rows.Next() {
		row := model.BlockIntervalStat{}
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, err
}

const (
	blockTimesForRecentBlocksQuery = `
		SELECT 
			MIN(height) start_height, 
			MAX(height) end_height, 
			MIN(time) start_time,
			MAX(time) end_time,
			COUNT(*) count, 
			EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff, 
			EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
		FROM
			( 
				SELECT * FROM blocks
				ORDER BY height DESC
				LIMIT ?
			) t;`

	blockTimesForIntervalQuery = `
		SELECT
			time_bucket($1, time) AS time_interval,
			COUNT(*) AS count,
			EXTRACT(EPOCH FROM (last(time, time) - first(time, time)) / COUNT(*)) AS avg
		FROM
			blocks
		WHERE
			(
				SELECT time
				FROM blocks
				ORDER BY time DESC
				LIMIT 1
			) - $2::INTERVAL < time
		GROUP BY time_interval
		ORDER BY time_interval ASC;`
)
