package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
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

// FindByHash returns a block with the matching hash
func (s BlocksStore) FindByHash(hash string) (*model.Block, error) {
	return s.FindBy("hash", hash)
}

// FindByHeight returns a block with the matching height
func (s BlocksStore) FindByHeight(height uint64) (*model.Block, error) {
	return s.FindBy("id", height)
}

// FindPrevious returns a block prior to the given height
func (s BlocksStore) FindPrevious(height uint64) (*model.Block, error) {
	block := &model.Block{}

	err := s.db.
		Order("id DESC").
		Limit(1).
		Find(block, "id < ?", height).
		Error

	return block, checkErr(err)
}

// Last returns the last indexed block
func (s BlocksStore) Last() (*model.Block, error) {
	block := &model.Block{}

	err := s.db.
		Order("id DESC").
		Limit(1).
		Find(&block).
		Error

	return block, checkErr(err)
}

// LastInEpoch returns the last block in the given epoch
func (s BlocksStore) LastInEpoch(epoch string) (*model.Block, error) {
	block := &model.Block{}

	err := s.db.
		Order("id DESC").
		Limit(1).
		Take(&block, "epoch = ?", epoch).
		Error

	return block, checkErr(err)
}

// Search returns matching blocks
func (s BlocksStore) Search() ([]model.Block, error) {
	result := []model.Block{}

	err := s.db.
		Order("id DESC").
		Limit(50).
		Find(&result).
		Error

	return result, err
}

// BlockTimes returns recent blocks averages
func (s BlocksStore) BlockTimes(limit int64) ([]byte, error) {
	return jsonquery.MustObject(s.db, queries.BlockTimes, limit)
}

// BlockStats returns block stats for a given interval
func (s BlocksStore) BlockStats(bucket string, limit uint) ([]byte, error) {
	return jsonquery.MustArray(s.db, queries.BlocksTimeStats, bucket, limit)
}

// Import creates block records in batch
func (s BlocksStore) Import(records []model.Block) error {
	now := time.Now()

	return s.bulkImport(queries.BlocksImport, len(records), func(i int) bulk.Row {
		r := records[i]

		return bulk.Row{
			r.ID,
			r.Time,
			r.Hash,
			r.Producer,
			r.Epoch,
			r.GasPrice,
			r.GasLimit,
			r.GasUsed,
			r.TotalSupply,
			r.ChunksCount,
			r.TransactionsCount,
			r.ApprovalsCount,
			now,
		}
	})
}
