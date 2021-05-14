package types

import "time"

type QueryParams struct {
	From time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To   time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
}
