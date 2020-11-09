package store

const (
	paginationLimit = 50
)

type Pagination struct {
	Page  uint `form:"page"`
	Limit uint `form:"limit"`
}

type PaginatedResult struct {
	Page    uint        `json:"page"`
	Pages   uint        `json:"pages"`
	Limit   uint        `json:"limit"`
	Count   uint        `json:"count"`
	Records interface{} `json:"records"`
}

func (p *Pagination) Validate() error {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = paginationLimit
	}
	if p.Limit >= paginationLimit {
		p.Limit = paginationLimit
	}
	return nil
}

func (p *PaginatedResult) Update() *PaginatedResult {
	if p.Pages == 0 && p.Count > 0 {
		pages := p.Count / p.Limit
		if pages*p.Limit < p.Count {
			pages++
		}
		p.Pages = pages
	}
	return p
}
