package sync

import (
	"github.com/figment-networks/near-indexer/model"
)

func CreateRun(c *Context) {
	c.Run = &model.Run{
		Height: c.Height.Height,
	}
	if err := c.DB.Runs.Create(c.Run); err != nil {
		c.Abort(err)
	}
}
