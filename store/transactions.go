package store

import (
	"strings"
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

// Search performs a transaction search and returns matching records
func (s TransactionsStore) Search(search TransactionsSearch) (*PaginatedResult, error) {
	if err := search.Validate(); err != nil {
		return nil, err
	}

	scope := s.db.
		Model(&model.Transaction{}).
		Order("id DESC").
		Limit(search.Limit)

	if search.BlockHash != "" {
		scope = scope.Where("block_hash = ?", search.BlockHash)
	}
	if search.BlockHeight > 0 {
		scope = scope.Where("height = ?", search.BlockHeight)
	}

	if search.Account != "" {
		accounts := strings.Split(search.Account, ",")
		scope = scope.Where("sender IN (?) OR receiver IN (?)", accounts, accounts)
	} else {
		if search.Sender != "" {
			scope = scope.Where("sender = ?", search.Sender)
		}
		if search.Receiver != "" {
			scope = scope.Where("receiver = ?", search.Receiver)
		}
	}

	if search.startTime != nil {
		scope = scope.Where("time >= ?", search.startTime)
	}
	if search.endTime != nil {
		scope = scope.Where("time <= ?", search.endTime)
	}

	var count uint
	if err := scope.Count(&count).Error; err != nil {
		return nil, err
	}

	transactions := []model.Transaction{}

	err := scope.
		Offset((search.Page - 1) * search.Limit).
		Limit(search.Limit).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	result := &PaginatedResult{
		Page:    search.Page,
		Limit:   search.Limit,
		Count:   count,
		Records: transactions,
	}

	return result.update(), nil
}

// Import imports transactions in bulk
func (s TransactionsStore) Import(records []model.Transaction) error {
	t := time.Now()

	return s.bulkImport(queries.TransactionsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.Hash,
			r.BlockHash,
			r.Height,
			r.Time,
			r.Sender,
			r.Receiver,
			r.Amount,
			r.GasBurnt,
			r.Success,
			string(r.Actions),
			r.ActionsCount,
			t,
			t,
		}
	})
}
