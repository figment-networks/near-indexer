package util

import (
	"time"
)

// ParseTime returns a timestamp from a unix ns epoch
func ParseTime(src int64) time.Time {
	return time.Unix(0, src)
}

// HourInterval returns a time interval for an hour
func HourInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()

	start := time.Date(year, month, day, t.Hour(), 0, 0, 0, t.Location())
	end := time.Date(year, month, day, t.Hour(), 59, 59, 0, t.Location())

	return start, end
}

// DayInterval returns a time interval for 24h
func DayInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	end := time.Date(year, month, day, 23, 59, 59, 0, t.Location())

	return start, end
}
