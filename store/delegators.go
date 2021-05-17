package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
	"time"
)

// DelegatorsStore handles operations on blocks
type DelegatorsStore struct {
	baseStore
}

// FetchRewardsByInterval fetches reward by interval
func (s *DelegatorsStore) FetchRewardsByInterval(account string, validatorId string, from time.Time, to time.Time, timeInterval model.TimeInterval) (model.RewardsSummary, error) {
	var res model.RewardsSummary
	err := s.db.Raw(queries.DelegatorsRewards, timeInterval.String(), account, validatorId, from, to, timeInterval.String()).Scan(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil
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
