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
}

func (p *timesIntervalParams) setDefaults() {
	if p.Interval == "" {
		p.Interval = "1h"
	}
	if p.Period == "" {
		p.Period = "48h"
	}
}