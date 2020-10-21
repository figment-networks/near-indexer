package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

// AccountsStore handles operations on account
type AccountsStore struct {
	baseStore
}

// FindByName returns an account for a name
func (s AccountsStore) FindByName(name string) (*model.Account, error) {
	result := &model.Account{}
	err := findBy(s.db, result, "name", name)
	return result, checkErr(err)
}

// Upsert creates a new account or updates the existing one
func (s AccountsStore) Upsert(acc *model.Account) error {
	existing, err := s.FindByName(acc.Name)
	if err != nil {
		if isNotFound(err) {
			return s.Create(acc)
		}
		return err
	}

	existing.LastHeight = acc.LastHeight
	existing.LastTime = acc.LastTime
	existing.Balance = acc.Balance
	existing.StakingBalance = acc.StakingBalance

	return s.Update(existing)
}

// BulkUpsert imports new records and updates existing ones
func (s AccountsStore) BulkUpsert(records []model.Account) error {
	t := time.Now()

	return s.Import(queries.AccountsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.Name,
			r.StartHeight,
			r.StartTime,
			r.LastHeight,
			r.LastTime,
			r.Balance,
			r.StakingBalance,
			t,
			t,
		}
	})
}
