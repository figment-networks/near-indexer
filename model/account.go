package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type Account struct {
	ID             int64        `json:"-"`
	Name           string       `json:"name"`
	StartHeight    types.Height `json:"start_height"`
	StartTime      time.Time    `json:"start_time"`
	LastHeight     types.Height `json:"last_height"`
	LastTime       time.Time    `json:"last_time"`
	Balance        types.Amount `json:"balance"`
	StakingBalance types.Amount `json:"staking_balance"`
	CreatedAt      time.Time    `json:"-"`
	UpdatedAt      time.Time    `json:"-"`
}

// Validate returns an error if account is invalid
func (acc Account) Validate() error {
	if acc.Name == "" {
		return errors.New("name is invalid")
	}
	if !acc.StartHeight.Valid() {
		return errors.New("start height is invalid")
	}
	if acc.StartTime.Year() == 1 {
		return errors.New("start time is invalid")
	}
	return nil
}
