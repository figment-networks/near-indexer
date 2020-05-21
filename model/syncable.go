package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

const (
	SyncableTypeBlock      = "block"
	SyncableTypeValidators = "validators"
)

// Syncable contains raw blockchain data
type Syncable struct {
	Model

	RunID       int64           `json:"run_id"`
	Height      types.Height    `json:"height"`
	Time        time.Time       `json:"time"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	ProcessedAt *time.Time      `json:"processed_at"`
}

// String returns a text representation of syncable
func (s Syncable) String() string {
	return fmt.Sprintf("type=%v height=%v", s.Type, s.Height)
}

// Validate returns an error if syncable is invalid
func (s Syncable) Validate() error {
	if s.Height <= 0 {
		return errors.New("height is invalid")
	}
	if s.Time.Year() == 1 {
		return errors.New("year is invalid")
	}
	if s.Type == "" {
		return errors.New("type is required")
	}
	if s.Data == nil {
		return errors.New("data is required")
	}
	return nil
}

// Decode decodes the raw data into a destination interface
func (s Syncable) Decode(dst interface{}) error {
	return json.Unmarshal(s.Data, dst)
}
