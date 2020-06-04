package sync

import (
	"log"
	"time"

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
	c.Lag = latest - c.BlockHeight

	defer func() {
		if c.BlockHeight > 0 {
			log.Printf("started height=%d retries=%d lag=%d",
				c.BlockHeight,
				c.Height.RetryCount,
				c.Lag,
			)
		}
	}()

	// Fetch the latest successfully processed height
	h, err := c.DB.Heights.Last()
	if err != nil {
		if err == store.ErrNotFound {
			if c.DefaultStartHeight > 0 {
				// Start with configured height
				log.Println("creating height from initial height config value")
				createNewHeight(c.DefaultStartHeight, c)
			} else {
				// Fetch start height from the genesis config
				log.Println("creating height from genesis config")
				createHeightFromGenesis(c)
			}
			return
		}
		c.Abort(err)
		return
	}
	if c.IsAborted() {
		return
	}

	hval := uint64(h.Height)

	// Retry the last height if it was not successful
	if h.ShouldRetry() {
		retryLastHeight(h, c)
		return
	}

	// Chain node is running behind the latest synced block.
	// This is an indication of a testnet reset.
	if latest < hval {
		c.Abort("chain height is behind")
		return
	}

	// We're up-to-date, no need to process anything
	if latest == hval {
		log.Println("already at latest height", latest)
		time.Sleep(time.Second)
		c.Abort(nil)
		return
	}

	createNewHeight(uint64(hval)+1, c)
}

func createHeightFromGenesis(c *Context) {
	config, err := c.Client.GenesisConfig()
	if err != nil {
		c.Abort(err)
		return
	}
	createNewHeight(config.GenesisHeight, c)
}

func createNewHeight(val uint64, c *Context) {
	c.BlockHeight = val
	c.Height = &model.Height{
		Status: model.HeightStatusPending,
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
