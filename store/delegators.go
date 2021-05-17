package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

// DelegatorsStore handles operations on blocks
type DelegatorsStore struct {
	baseStore
}

// Import creates new validators in batch
func (s DelegatorsStore) ImportDelegatorEpochs(records []model.DelegatorEpoch) error {
	return s.bulkImport(queries.DelegatorEpochsImport, len(records), func(i int) bulk.Row {
		r := records[i]
		return bulk.Row{
			r.AccountID,
			r.ValidatorID,
			r.Epoch,
			r.LastHeight,
			r.LastTime,
			r.StakedBalance,
			r.UnstakedBalance,
			r.Reward,
		}
	})
}

