package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

const (
	ScopeStaking           = "staking"
	ActionValidatorAdded   = "joined_active_set"
	ActionValidatorRemoved = "left_active_set"
	ActionValidatorKicked  = "kicked"
	ItemTypeValidator      = "validator"
)

type Event struct {
	ID          int       `json:"id"`
	Scope       string    `json:"scope"`
	Action      string    `json:"action"`
	BlockHeight uint64    `json:"block_height"`
	BlockTime   time.Time `json:"block_time"`
	Epoch       string    `json:"epoch"`
	ItemID      string    `json:"item_id"`
	ItemType    string    `json:"item_type"`
	Metadata    types.Map `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Event) TableName() string {
	return "events"
}

func (e Event) Validate() error {
	if e.Scope == "" {
		return errors.New("scope is required")
	}
	if e.Action == "" {
		return errors.New("action is required")
	}
	return nil
}
