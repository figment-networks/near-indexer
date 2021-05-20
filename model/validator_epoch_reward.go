package model

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type ValidatorEpochReward struct {
	ID                int64        `json:"-"`
	AccountID         string       `json:"-"`
	Epoch             string       `json:"epoch"`
	DistributedHeight types.Height `json:"distributed_height"`
	DistributedTime   time.Time    `json:"distributed_time"`
	RewardFee         *int         `json:"reward_fee"`
	Reward            types.Amount `json:"reward"`
}

func (ValidatorEpochReward) TableName() string {
	return "validator_epochs_rewards"
}
