package server

import "errors"

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
