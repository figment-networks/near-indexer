package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type Transaction struct {
	Model

	Time      time.Time    `json:"time"`
	Height    types.Height `json:"height"`
	Hash      string       `json:"hash"`
	BlockHash string       `json:"block_hash"`
	Signer    string       `json:"signer"`
	SignerKey string       `json:"signer_key"`
	Signature string       `json:"signature"`
	Receiver  string       `json:"receiver"`
	Amount    types.Amount `json:"amount"`
	GasBurnt  string       `json:"gas_burnt"`
	Actions   []byte       `json:"actions"`
}

type TransactionAction struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Validate returns an error if transaction is invalid
func (t Transaction) Validate() error {
	if !t.Height.Valid() {
		return errors.New("height is invalid")
	}
	if t.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if t.Hash == "" {
		return errors.New("hash is required")
	}
	if t.Signer == "" {
		return errors.New("signer is required")
	}
	return nil
}
