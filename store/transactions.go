package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

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

func (s TransactionsStore) Import(records []model.Transaction) error {
	t := time.Now()

	return s.bulkImport(queries.TransactionsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.Hash,
			r.BlockHash,
			r.Height,
			r.Time,
			r.Signer,
			r.SignerKey,
			r.Receiver,
			r.Signature,
			r.Amount,
			r.GasBurnt,
			r.Success,
			r.Actions,
			t,
			t,
		}
	})
}
