package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

var (
	errBlockInvalidHeight   = errors.New("height is invalid")
	errBlockInvalidHash     = errors.New("hash is invalid")
	errBlockInvalidProducer = errors.New("producer is invalid")
	errBlockInvalidTime     = errors.New("block time is required")
	errBlockInvalidEpoch    = errors.New("epoch is invalid")
)

type Block struct {
	ID                types.Height `json:"id"`
	Time              time.Time    `json:"time"`
	Producer          string       `json:"producer"`
	Hash              string       `json:"hash"`
	Epoch             string       `json:"epoch"`
	GasPrice          types.Amount `json:"gas_price"`
	GasLimit          uint         `json:"gas_allowed"`
	GasUsed           uint         `json:"gas_used"`
	TotalSupply       types.Amount `json:"total_supply"`
	ChunksCount       int          `json:"chunks_count"`
	TransactionsCount int          `json:"transactions_count"`
	ApprovalsCount    int          `json:"approvals_count"`
	CreatedAt         time.Time    `json:"created_at"`
}

// Validate returns an error if block data is invalid
func (b Block) Validate() error {
	if !b.ID.Valid() {
		return errBlockInvalidHeight
	}
	if b.Hash == "" {
		return errBlockInvalidHash
	}
	if b.Producer == "" {
		return errBlockInvalidProducer
	}
	if b.Epoch == "" {
		return errBlockInvalidEpoch
	}
	if b.Time.IsZero() {
		return errBlockInvalidTime
	}
	return nil
}
