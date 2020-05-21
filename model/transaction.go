package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type Transaction struct {
	Model

	Height    types.Height `json:"height"`
	Time      time.Time    `json:"time"`
	Hash      string       `json:"hash"`
	Type      string       `json:"type"`
	Signer    string       `json:"signer"`
	SignerKey string       `json:"signer_key"`
	Signature string       `json:"signature"`
	Receiver  string       `json:"receiver"`
	Amount    types.Amount `json:"amount"`
	Fee       types.Amount `json:"fee"`
}

// Validate returns an error if transaction is invalid
func (t Transaction) Validate() error {
	if !t.Height.Valid() {
		return errors.New("height is invalid")
	}
	if t.Time.Year() == 1 {
		return errors.New("year is invalid")
	}
	if t.Hash == "" {
		return errors.New("hash is required")
	}
	if t.Signer == "" {
		return errors.New("signer is required")
	}
	return nil
}
