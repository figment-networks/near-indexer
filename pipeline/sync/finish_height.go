package sync

import (
	"log"
	"strings"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

func FinishRun(c *Context) {
	finishRun(c)
	finishHeight(c)
}

func finishHeight(c *Context) {
	// No height was set, the chain was aborted at the very beginning
	if c.Height == nil {
		return
	}

	if err := c.FirstError(); err == nil {
		c.Height.Status = model.HeightStatusOK
		c.Height.Error = nil
	} else {
		msg := err.Error()
		c.Height.Status = model.HeightStatusError
		c.Height.Error = &msg

		// Indicate that height does not have a block
		if err == near.ErrNotFound {
			c.Height.Status = model.HeightStatusNoBlock
		}
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
		err := errorsString(c.errors)

		c.Run.Success = false
		c.Run.Error = &err
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
