package model

import (
	"errors"
	"time"
)

type Block struct {
	Model

	Height          uint64    `json:"height"`
	Time            time.Time `json:"time"`
	Producer        string    `json:"producer"`
	Hash            string    `json:"hash"`
	PrevHash        string    `json:"prev_hash"`
	GasPrice        string    `json:"gas_price"`
	RentPaid        string    `json:"rent_paid"`
	ValidatorReward string    `json:"validator_reward"`
	TotalSupply     string    `json:"total_supply"`
	Signature       string    `json:"signature"`
	ChunksCount     int       `json:"chunks_count"`
}

// Validate returns an error if block data is invalid
func (b Block) Validate() error {
	if b.Height == 0 {
		return errors.New("height is invalid")
	}
	if b.Time.Year() == 1 {
		return errors.New("time is invalid")
	}
	return nil
}
