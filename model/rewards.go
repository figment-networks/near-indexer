package model

import (
	"github.com/figment-networks/near-indexer/model/types"
)

type RewardsSummary struct {
	Interval string       `json:"interval"`
	Amount   types.Amount `json:"amount"`
}

type TimeInterval uint

const (
	TimeIntervalDaily TimeInterval = iota
	TimeIntervalMonthly
	TimeIntervalYearly
)

var (
	TimeIntervalTypes = map[string]TimeInterval{
		"daily":   TimeIntervalDaily,
		"monthly": TimeIntervalMonthly,
		"yearly":  TimeIntervalYearly,
	}
)

func GetTypeForTimeInterval(s string) (TimeInterval, bool) {
	t, ok := TimeIntervalTypes[s]
	return t, ok
}

func (k TimeInterval) String() string {
	switch k {
	case TimeIntervalDaily:
		return "DD-MM-YYYY"
	case TimeIntervalMonthly:
		return "MM-YYYY"
	case TimeIntervalYearly:
		return "YYYY"
	default:
		return "unknown"
	}
}
