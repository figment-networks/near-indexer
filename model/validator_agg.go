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
	Slashed        bool         `json:"slashed"`
	Stake          types.Amount `json:"stake"`
	Efficiency     float64      `json:"efficiency"`
}

func (ValidatorAgg) TableName() string {
	return "validator_aggregates"
}

type ValidatorEpoch struct {
	ID             int64        `json:"-"`
	AccountID      string       `json:"-"`
	Epoch          string       `json:"epoch"`
	LastHeight     types.Height `json:"last_height"`
	LastTime       time.Time    `json:"last_time"`
	ExpectedBlocks int          `json:"expected_blocks"`
	ProducedBlocks int          `json:"produced_blocks"`
	Efficiency     float64      `json:"efficiency"`
}

func (ValidatorEpoch) TableName() string {
	return "validator_epochs"
}
