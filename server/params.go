package server

import (
	"errors"
	"time"

	"github.com/figment-networks/near-indexer/model"
)

type statsParams struct {
	Bucket string `form:"bucket"`
	Limit  uint   `form:"limit"`
}

type blockTimesParams struct {
	Limit int64 `form:"limit"`
}

type accountsIndexParams struct {
	Height int64 `form:"height"`
}

func (p *blockTimesParams) setDefaults() {
	if p.Limit <= 0 {
		p.Limit = 100
	}
	if p.Limit > 1000 {
		p.Limit = 1000
	}
}

func (p *statsParams) Validate() error {
	if p.Bucket == "" {
		p.Bucket = "h"
	}

	switch p.Bucket {
	case "d":
		if p.Limit == 0 {
			p.Limit = 30
		}
		if p.Limit > 90 {
			return errors.New("maximum daily limit is 90")
		}
	case "h":
		if p.Limit == 0 {
			p.Limit = 24
		}
		if p.Limit > 48 {
			return errors.New("max hourly limit is 48")
		}
	default:
		return errors.New("invalid time bucket: " + p.Bucket)
	}

	return nil
}

type rewardsParams struct {
	From     time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To       time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
	Interval string    `form:"interval" binding:"required" `
}

type delegatorRewardsParams struct {
	rewardsParams
	ValidatorId string `form:"validator_id" binding:"-" `
}

func (p *rewardsParams) Validate() error {
	if p.From.IsZero() && p.To.IsZero() {
		return errors.New("invalid time range: " + "")
	}

	var ok bool
	if _, ok = model.GetTypeForTimeInterval(p.Interval); !ok {
			return errors.New("time interval type is wrong")
	}

	return nil
}