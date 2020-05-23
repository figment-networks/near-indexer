package sync

import (
	"log"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/store"
)

func CreateHeight(c *Context) {
	status, err := c.Client.Status()
	if err != nil {
		c.Abort(err)
		return
	}
	latest := status.SyncInfo.LatestBlockHeight
	c.Status = &status

	defer func() {
		if c.BlockHeight > 0 {
			log.Printf("started height=%d retries=%d lag=%d",
				c.BlockHeight,
				c.Height.RetryCount,
				latest-c.BlockHeight,
			)
		}
	}()

	h, err := c.DB.Heights.Last()
	if err != nil {
		if err == store.ErrNotFound {
			createNewHeight(latest, c)
			return
		}
		c.Abort(err)
		return
	}
	hval := uint64(h.Height)

	// Retry the last height if it was not successful
	if h.ShouldRetry() {
		retryLastHeight(h, c)
		return
	}

	// Node is behind for some reason
	if latest < hval {
		c.Abort("chain height is behind")
		return
	}

	// We're up-to-date, no need to process anything
	if latest == hval {
		c.Abort(nil)
		return
	}

	createNewHeight(uint64(hval)+1, c)
}

func createNewHeight(val uint64, c *Context) {
	c.BlockHeight = val
	c.Height = &model.Height{
		Height: types.Height(val),
	}
	if err := c.DB.Heights.Create(c.Height); err != nil {
		c.Abort(err)
	}
}

func retryLastHeight(h *model.Height, c *Context) {
	h.ResetForRetry()

	if err := c.DB.Heights.Update(h); err != nil {
		c.Abort(err)
		return
	}

	c.Height = h
	c.BlockHeight = uint64(h.Height)
}
