package store

import (
	"strconv"
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
func (s *DelegatorsStore) FetchRewardsByInterval(account string, validatorId string, from time.Time, to time.Time, timeInterval model.TimeInterval) ([]model.RewardsSummary, error) {
	q := "SELECT to_char(distributed_at_time, $INTERVAL) AS interval, validator_id as validator, SUM(reward) AS amount FROM delegator_epochs "

	var (
		args   []interface{}
		wherec []string
		i      = 1
	)

	if account != "" {
		wherec = append(wherec, ` account_id =  $`+strconv.Itoa(i))
		args = append(args, account)
		i++
	}
	if validatorId != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, validatorId)
		i++
	}
	if !from.IsZero() {
		wherec = append(wherec, ` distributed_at_time > $`+strconv.Itoa(i))
		args = append(args, from)
		i++
	}
	if !to.IsZero() {
		wherec = append(wherec, ` distributed_at_time < $`+strconv.Itoa(i))
		args = append(args, to)
		i++
	}

	wherec = append(wherec, ` reward >  $`+strconv.Itoa(i))
	args = append(args, 0)
	i++

	q += ` WHERE `
	q += strings.Join(wherec, " AND ")
	q += " GROUP BY  to_char(distributed_at_time, $INTERVAL), validator_id ORDER BY to_char(distributed_at_time, $INTERVAL)"
	q = strings.Replace(q, "$INTERVAL", "'"+timeInterval.String()+"'", -1)

	var res []model.RewardsSummary
	var err error
	rows, err := s.db.Raw(q, args).Rows()

	defer rows.Close()

	for rows.Next() {
		r := model.RewardsSummary{}
		if err = rows.Scan(&r.Interval, &r.Validator, &r.Amount); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

// FindDelegatorEpochBy returns delegator epoch by epoch and account id
func (s DelegatorsStore) FindDelegatorEpochBy(epoch string, accountId string, validatorId string) (*model.DelegatorEpoch, error) {
	res := &model.DelegatorEpoch{}
	err := s.db.Where("distributed_at_epoch = ? AND account_id = ? AND validator_id = ?", epoch, accountId, validatorId).Limit(1).Take(res).Error
	return res, checkErr(err)
}

// SearchDelegatorEpochs performs a delegator epoch search and returns matching records
func (s DelegatorsStore) SearchDelegatorEpochs(search DelegatorEpochsSearch) ([]model.DelegatorEpoch, error) {
	if err := search.Validate(); err != nil {
		return nil, err
	}

	scope := s.db.Model(&model.DelegatorEpoch{})

	if search.Epoch != "" {
		scope = scope.Where("epoch = ?", search.Epoch)
	}
	if search.ValidatorID != "" {
		scope = scope.Where("validator_id = ?", search.ValidatorID)
	}
	if search.AccountID != "" {
		scope = scope.Where("account_id = ?", search.AccountID)
	}

	delegatorEpochs := []model.DelegatorEpoch{}
	err := scope.Find(&delegatorEpochs).Error
	if err != nil {
		return nil, err
	}

	return delegatorEpochs, nil
}

// ImportDelegatorEpochs creates new validators in batch
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
				r.DistributedAtEpoch,
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
