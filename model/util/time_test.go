package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	assert.Equal(t, int64(0), ParseTime(0).Unix())
	assert.Equal(t, int64(1596166782), ParseTime(1596166782911378000).Unix())
	assert.Equal(t, int64(1596166782911378000), ParseTime(1596166782911378000).UnixNano())
	assert.Equal(t, "2020-07-30T22:39:42-05:00", ParseTime(1596166782911378000).Format(time.RFC3339))
	assert.Equal(t, "2020-07-30T22:39:42.911378-05:00", ParseTime(1596166782911378000).Format(time.RFC3339Nano))
}

func TestHourInterval(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2020-07-30T22:39:42-05:00")
	assert.NoError(t, err)

	start, end := HourInterval(now)
	assert.Equal(t, "2020-07-30T22:00:00-05:00", start.Format(time.RFC3339))
	assert.Equal(t, "2020-07-30T22:59:59-05:00", end.Format(time.RFC3339))
}

func TestDayInterval(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2020-07-30T22:39:42-05:00")
	assert.NoError(t, err)

	start, end := DayInterval(now)
	assert.Equal(t, "2020-07-30T00:00:00-05:00", start.Format(time.RFC3339))
	assert.Equal(t, "2020-07-30T23:59:59-05:00", end.Format(time.RFC3339))
}
