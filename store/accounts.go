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

// Import imports new records and updates existing ones
func (s AccountsStore) Import(records []model.Account) error {
	t := time.Now()
	rr := "INSERT INTO accounts (  name,  start_height,  start_time,  last_height,  last_time, balance,  staking_balance,  created_at,  updated_at) VALUES @values ON CONFLICT (name) DO UPDATE SET  last_height     = excluded.last_height,  last_time       = excluded.last_time,  balance         = excluded.balance,  staking_balance = excluded.staking_balance,  updated_at      = excluded.updated_at"
	return s.bulkImport(rr, len(records), func(i int) bulk.Row {
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
