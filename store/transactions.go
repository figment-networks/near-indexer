package store

import "github.com/figment-networks/near-indexer/model"

type TransactionsStore struct {
	baseStore
}

// FindByBlock returns a collection of transactions for a block hash
func (s TransactionsStore) FindByBlock(hash string) ([]model.Transaction, error) {
	result := []model.Transaction{}

	err := s.db.
		Model(&model.Transaction{}).
		Order("id DESC").
		Limit(100).
		Find(&result, "block_hash = ?", hash).
		Error

	return result, checkErr(err)
}

// FindByHash returns a transaction record by hash
func (s TransactionsStore) FindByHash(hash string) (*model.Transaction, error) {
	result := &model.Transaction{}

	err := s.db.
		Model(result).
		Take(result, "hash = ?", hash).
		Limit(1).
		Error

	return result, checkErr(err)
}

// Recent returns the most recent N transactions
func (s TransactionsStore) Recent(n int) ([]model.Transaction, error) {
	result := []model.Transaction{}

	err := s.db.
		Model(&model.Transaction{}).
		Order("id DESC").
		Limit(n).
		Find(&result).
		Error

	return result, checkErr(err)
}
