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

// Import imports new records and updates existing ones
func (s AccountsStore) Import(records []model.Account) error {
	t := time.Now()

	return s.bulkImport(queries.AccountsImport, len(records), func(i int) bulk.Row {
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
