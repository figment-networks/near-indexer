package sync

import (
	"log"
	"strings"
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

func FinishHeight(c *Context) {
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
		c.Height.Error = &msg

		switch err {
		case near.ErrBlockNotFound:
			c.Height.Status = model.HeightStatusNoBlock
		case near.ErrBlockMissing:
			c.Height.Status = model.HeightStatusMissing
		default:
			c.Height.Status = model.HeightStatusError
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

	c.Run.Duration = time.Since(c.Height.CreatedAt).Milliseconds()
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
