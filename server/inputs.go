package server

type blockTimesParams struct {
	Limit int64 `form:"limit"`
}

type timesIntervalParams struct {
	Interval string `form:"interval"`
	Period   string `form:"period"`
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

func (p *timesIntervalParams) setDefaults() {
	if p.Interval == "" {
		p.Interval = "h"
	}

	if p.Period == "" {
		if p.Interval == "h" {
			p.Period = "48"
		}
		if p.Interval == "d" {
			p.Period = "30"
		}
	}
}
