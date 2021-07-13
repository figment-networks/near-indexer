package store

import (
	"errors"
	"regexp"
	"time"
)

var (
	reDate = regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}$`)
)

type TransactionsSearch struct {
	Pagination

	BlockHash   string `form:"block_hash"`
	BlockHeight uint64 `form:"block_height"`
	Sender      string `form:"sender"`
	Receiver    string `form:"receiver"`
	Account     string `form:"account"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`

	startTime *time.Time
	endTime   *time.Time
}

func (s *TransactionsSearch) Validate() error {
	if s.Page == 0 {
		s.Page = 1
	}
	if s.Limit == 0 {
		s.Limit = paginationLimit
	}
	if s.Limit >= paginationLimit {
		s.Limit = paginationLimit
	}

	if t, err := parseTimeFilter(s.StartDate); err == nil {
		s.startTime = t
	} else {
		return errors.New("start time is invalid")
	}
	if t, err := parseTimeFilter(s.EndDate); err == nil {
		s.endTime = t
	} else {
		return errors.New("end time is invalid")
	}

	return nil
}

func parseTimeFilter(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}

	var t time.Time
	var err error

	if reDate.MatchString(input) {
		t, err = time.Parse("2006-01-02", input)
	} else {
		t, err = time.Parse(time.RFC3339, input)
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}
