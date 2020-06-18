package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
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

	return s.Import(sqlAccountsBulkUpsert, len(records), func(i int) bulk.Row {
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

var (
	sqlAccountsBulkUpsert = `
		INSERT INTO accounts (
			name,
			start_height,
			start_time,
			last_height,
			last_time,
			balance,
			staking_balance,
			created_at,
			updated_at
		)
		VALUES @values
		ON CONFLICT (name) DO UPDATE
		SET
			last_height     = excluded.last_height,
			last_time       = excluded.last_time,
			balance         = excluded.balance,
			staking_balance = excluded.staking_balance,
			updated_at      = excluded.updated_at`
)
