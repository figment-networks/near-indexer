package store

import (
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

const batchSize = 100

// DelegatorsStore handles operations on blocks
type DelegatorsStore struct {
	baseStore
}

// FetchRewardsByInterval fetches reward by interval
func (s *DelegatorsStore) FetchRewardsByInterval(account string, validatorId string, from time.Time, to time.Time, timeInterval model.TimeInterval) (model.RewardsSummary, error) {
	var res model.RewardsSummary
	q := strings.Replace(queries.DelegatorsRewards, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	var err error
	if validatorId == "" {
		q = strings.Replace(q, "AND validator_id = ?", "", -1)
		err = s.db.Raw(q, account, from, to).Scan(&res).Error
	} else {
		err = s.db.Raw(q, account, validatorId, from, to).Scan(&res).Error
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// Import creates new validators in batch
func (s DelegatorsStore) ImportDelegatorEpochs(records []model.DelegatorEpoch) error {
	var err error
	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}
		err = s.bulkImport(queries.DelegatorEpochsImport, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.AccountID,
				r.ValidatorID,
				r.Epoch,
				r.DistributedAtHeight,
				r.DistributedAtTime,
				r.StakedBalance,
				r.UnstakedBalance,
				r.Reward,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}
