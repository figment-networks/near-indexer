package model

import (
	"github.com/figment-networks/near-indexer/model/types"
)

type RewardsResponse struct {
	Month  string       `json:"month"`
	Amount types.Amount `json:"amount"`
}
