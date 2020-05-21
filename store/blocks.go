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
