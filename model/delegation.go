package model

import (
	"errors"

	"github.com/figment-networks/near-indexer/model/types"
)

type Delegation struct {
	ID              int          `json:"-"`
	Account         string       `json:"account"`
	UnstakedBalance types.Amount `json:"unstaked_balance"`
	StakedBalance   types.Amount `json:"staked_balance"`
	CanWithdraw     bool         `json:"can_withdraw"`
}

func (Delegation) TableName() string {
	return "delegations"
}

func (d Delegation) Validate() error {
	if d.Account == "" {
		return errors.New("account is required")
	}
	if !d.StakedBalance.Valid() {
		return errors.New("staked balance is invalid")
	}
	if !d.UnstakedBalance.Valid() {
		return errors.New("unstaked balance is invalid")
	}
	return nil
}
