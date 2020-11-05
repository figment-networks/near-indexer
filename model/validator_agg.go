package model

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

// ValidatorAgg represents validator data at latest height
type ValidatorAgg struct {
	Model

	AccountID      string       `json:"account_id"`
	StartHeight    types.Height `json:"start_height"`
	StartTime      time.Time    `json:"start_time"`
	LastHeight     types.Height `json:"last_height"`
	LastTime       time.Time    `json:"last_time"`
	ExpectedBlocks int          `json:"expected_blocks"`
	ProducedBlocks int          `json:"produced_blocks"`
	Active         bool         `json:"active"`
	Slashed        bool         `json:"slashed"`
	Stake          types.Amount `json:"stake"`
	Efficiency     float64      `json:"efficiency"`
	RewardFee      *int         `json:"reward_fee"`
}

func (ValidatorAgg) TableName() string {
	return "validator_aggregates"
}
