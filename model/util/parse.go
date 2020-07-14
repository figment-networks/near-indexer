package util

import (
	"strconv"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
)

// ParseTime returns a timestamp from a unix ns epoch
func ParseTime(src int64) time.Time {
	return time.Unix(0, src)
}

// ParseTimeFromString returns a timestamp from a unix ns epoch in string format
func ParseTimeFromString(src string) time.Time {
	val, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return ParseTime(val)
}

// ParseAmount returns amount value
func ParseAmount(src string) types.Amount {
	return types.NewAmount(src)
}

// ParseHeight returns a height
func ParseHeight(src uint64) types.Height {
	return types.Height(src)
}
