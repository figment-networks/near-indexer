package model

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type ValidatorEpochReward struct {
	ID                  int64        `json:"-"`
	AccountID           string       `json:"-"`
	Epoch               string       `json:"epoch"`
	DistributedAtHeight types.Height `json:"distributed_at_height"`
	DistributedAtTime   time.Time    `json:"distributed_at_time"`
	RewardFee           *int         `json:"reward_fee"`
	Reward              types.Amount `json:"reward"`
}

func (ValidatorEpochReward) TableName() string {
	return "validator_epochs_rewards"
}
