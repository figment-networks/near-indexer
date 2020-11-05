package model

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

type Block struct {
	ID                int64        `json:"-"`
	Height            types.Height `json:"height"`
	Time              time.Time    `json:"time"`
	Producer          string       `json:"producer"`
	Hash              string       `json:"hash"`
	PrevHash          string       `json:"prev_hash"`
	Epoch             string       `json:"epoch"`
	GasPrice          types.Amount `json:"gas_price"`
	GasLimit          int          `json:"gas_allowed"`
	GasUsed           int          `json:"gas_used"`
	RentPaid          types.Amount `json:"rent_paid"`
	ValidatorReward   types.Amount `json:"validator_reward"`
	TotalSupply       types.Amount `json:"total_supply"`
	Signature         string       `json:"signature"`
	ChunksCount       int          `json:"chunks_count"`
	TransactionsCount int          `json:"transactions_count"`
	ApprovalsCount    int          `json:"approvals_count"`
	CreatedAt         time.Time    `json:"created_at"`
}

// Validate returns an error if block data is invalid
func (b Block) Validate() error {
	if b.Hash == "" {
		return errors.New("hash is required")
	}
	if b.Producer == "" {
		return errors.New("producer is required")
	}
	if !b.Height.Valid() {
		return errors.New("height is invalid")
	}
	if b.Time.IsZero() {
		return errors.New("block time is required")
	}
	return nil
}
