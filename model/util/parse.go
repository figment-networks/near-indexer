package util

import (
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

// ParseTime returns a timestamp from a unix ns epoch
func ParseTime(src int64) time.Time {
	return time.Unix(0, src)
}

// ParseAmount returns amount value
func ParseAmount(src string) types.Amount {
	return types.NewAmount(src)
}

// ParseHeight returns a height
func ParseHeight(src uint64) types.Height {
	return types.Height(src)
}
