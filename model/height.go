package model

import (
	"errors"

	"github.com/figment-networks/near-indexer/model/types"
)

var (
	HeightStatusOK      = "success"  // Successful result
	HeightStatusError   = "error"    // Something went wrong
	HeightStatusSkip    = "skip"     // Marked as skip
	HeightStatusNoBlock = "no_block" // No block at height
)

type Height struct {
	Model

	Height types.Height `json:"height"`
	Status string       `json:"status"`
	Error  *string      `json:"error"`
}

func (h Height) Validate() error {
	if !h.Height.Valid() {
		return errors.New("height is invalid")
	}
	if h.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
