package model

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type ValidatorEpoch struct {
	ID                int64        `json:"-"`
	AccountID         string       `json:"-"`
	Epoch             string       `json:"epoch"`
	LastHeight        types.Height `json:"last_height"`
	LastTime          time.Time    `json:"last_time"`
	ExpectedBlocks    int          `json:"expected_blocks"`
	ProducedBlocks    int          `json:"produced_blocks"`
	Efficiency        float64      `json:"efficiency"`
	StakingBalance    types.Amount `json:"staking_balance"`
	RewardFee         *int         `json:"reward_fee"`
	Reward types.Amount    `json:"reward"`
}

func (ValidatorEpoch) TableName() string {
	return "validator_epochs"
}
