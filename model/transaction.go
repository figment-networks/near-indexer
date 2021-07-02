package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

var (
	errTxHeightInvalid   = errors.New("height is invalid")
	errTxTimeInvalid     = errors.New("time is invalid")
	errTxHashInvalid     = errors.New("hash is invalid")
	errTxSenderInvalid   = errors.New("sender is invalid")
	errTxReceiverInvalid = errors.New("receiver is invalid")
)

type Transaction struct {
	Model

	Time         time.Time       `json:"time"`
	Height       types.Height    `json:"height"`
	Hash         string          `json:"hash"`
	BlockHash    string          `json:"block_hash"`
	Sender       string          `json:"sender"`
	Receiver     string          `json:"receiver"`
	GasBurnt     string          `json:"gas_burnt"`
	Fee          string          `json:"fee"`
	PublicKey    string          `json:"public_key"`
	Signature    string          `json:"signature"`
	Actions      json.RawMessage `json:"actions"`
	ActionsCount int             `json:"actions_count"`
	Outcome      json.RawMessage `json:"outcome"`
	Receipt      json.RawMessage `json:"receipt"`
	Success      bool            `json:"success"`
}

type TransactionAction struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Validate returns an error if transaction is invalid
func (t Transaction) Validate() error {
	if !t.Height.Valid() {
		return errTxHeightInvalid
	}
	if t.Time.IsZero() {
		return errTxTimeInvalid
	}
	if t.Hash == "" {
		return errTxHashInvalid
	}
	if t.Sender == "" {
		return errTxSenderInvalid
	}
	if t.Receiver == "" {
		return errTxReceiverInvalid
	}
	return nil
}
