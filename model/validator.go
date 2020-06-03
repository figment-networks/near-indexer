package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type Validator struct {
	Model

	Height         types.Height `json:"height"`
	Time           time.Time    `json:"time"`
	AccountID      string       `json:"account_id"`
	Epoch          string       `json:"epoch"`
	ExpectedBlocks int          `json:"expected_blocks"`
	ProducedBlocks int          `json:"produced_blocks"`
	Slashed        bool         `json:"slashed"`
	Stake          types.Amount `json:"stake"`
	Efficiency     float64      `json:"efficiency"`
}

type ValidatorIntervalStat struct {
	TimeInterval string `json:"time_interval"`
	Count        int    `json:"count"`
}

func (Validator) TableName() string {
	return "validators"
}

func (v Validator) Validate() error {
	if !v.Height.Valid() {
		return errors.New("height is invalid")
	}
	if v.Time.Year() == 1 {
		return errors.New("time is invalid")
	}
	if v.AccountID == "" {
		return errors.New("account id is required")
	}
	return nil
}
