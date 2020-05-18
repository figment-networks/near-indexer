package model

import (
	"errors"
	"time"
)

type Validator struct {
	Model

	Height         uint64    `json:"height"`
	Time           time.Time `json:"time"`
	PublicKey      string    `json:"public_key"`
	AccountID      string    `json:"account_id"`
	ExpectedBlocks int       `json:"expected_blocks"`
	ProducedBlocks int       `json:"produced_blocks"`
	Stake          string    `json:"stake"`
	Efficiency     float32   `json:"efficiency"`
}

func (v Validator) Validate() error {
	if v.Height == 0 {
		return errors.New("height is invalid")
	}
	if v.Time.Year() == 1 {
		return errors.New("time is invalid")
	}
	if v.PublicKey == "" {
		return errors.New("public key is required")
	}
	if v.AccountID == "" {
		return errors.New("account id is required")
	}
	return nil
}
