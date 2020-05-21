package sync

import (
	"log"
	"strings"

	"github.com/figment-networks/near-indexer/model"
)

func FinishRun(c *Context) {
	finishRun(c)
	finishHeight(c)
}

func finishHeight(c *Context) {
	if c.Height == nil {
		return
	}

	if len(c.errors) == 0 {
		c.Height.Status = model.HeightStatusOK
	} else {
		c.Height.Status = model.HeightStatusError
		*c.Height.Error = c.LastError().Error()
	}

	if err := c.DB.Heights.Update(c.Height); err != nil {
		log.Println("cant update height:", err)
	}
}

func finishRun(c *Context) {
	if c.Run == nil {
		return
	}

	c.Run.Success = true
	c.Run.Error = nil

	if len(c.errors) > 0 {
		c.Run.Success = false
		*c.Run.Error = errorsString(c.errors)
	}

	if err := c.DB.Runs.Update(c.Run); err != nil {
		log.Println("cant update run:", err)
	}
}

func errorsString(errs []error) string {
	lines := make([]string, len(errs))
	for idx, e := range errs {
		lines[idx] = e.Error()
	}
	return strings.Join(lines, "\n")
}
