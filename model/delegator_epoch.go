package model

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type DelegatorEpoch struct {
	ID              int64        `json:"-"`
	AccountID       string       `json:"account_id"`
	ValidatorID     string       `json:"validator_id"`
	Epoch           string       `json:"epoch"`
	LastHeight      types.Height `json:"last_height"`
	LastTime        time.Time    `json:"last_time"`
	StakedBalance   types.Amount `json:"staked_balance"`
	UnstakedBalance types.Amount `json:"unstaked_balance"`
	Reward          types.Amount `json:"reward"`
}

func (DelegatorEpoch) TableName() string {
	return "delegator_epochs"
}
