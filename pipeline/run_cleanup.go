package pipeline

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/store"
)

func RunCleanup(cfg *config.Config, db *store.Store) error {
	lastBlock, err := db.Blocks.Last()
	if err != nil {
		return err
	}

	maxHeight := uint64(lastBlock.ID) - uint64(cfg.CleanupThreshold)
	if maxHeight == 0 {
		log.Println("nothing to cleanup")
		return nil
	}
	log.Println("starting cleanup, max height:", maxHeight, "threshold:", cfg.CleanupThreshold)

	return nil

	if numRows, err := db.Validators.Cleanup(maxHeight); err == nil {
		log.Println("validators removed:", numRows)
	} else {
		log.Println("validators cleanup error:", err)
	}

	if err := db.Validators.CleanupCounts(); err != nil {
		log.Println("validator counts cleanup error:", err)
	}

	return nil
}
