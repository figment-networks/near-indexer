package pipeline

import (
	"time"

	"github.com/figment-networks/near-indexer/near"
)

type Payload struct {
	Lag          int
	StartHeight  uint64
	StartTime    time.Time
	EndHeight    uint64
	EndTime      time.Time
	CurrentBlock *near.Block
	Heights      []*HeightPayload
}
