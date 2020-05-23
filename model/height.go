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

	Height     types.Height `json:"height"`
	Status     string       `json:"status"`
	RetryCount int          `json:"retry_count"`
	Error      *string      `json:"error"`
}

type HeightStatusCount struct {
	Status string
	Num    int
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

// ShouldRetry returns true if height is retriable
func (h Height) ShouldRetry() bool {
	return h.Status == HeightStatusError && h.RetryCount < 3
}

// ResetForRetry clears the errors
func (h *Height) ResetForRetry() {
	h.Error = nil
	h.RetryCount++
}
