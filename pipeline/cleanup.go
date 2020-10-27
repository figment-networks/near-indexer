package pipeline

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/store"
)

func RunCleanup(cfg *config.Config, db *store.Store) error {
	h, err := db.Heights.LastSuccessful()
	if err != nil {
		return err
	}

	maxHeight := uint64(h.Height) - uint64(cfg.CleanupThreshold)
	if maxHeight == 0 {
		log.Println("nothing to cleanup")
		return nil
	}

	log.Println("starting cleanup, max height:", maxHeight, "threshold:", cfg.CleanupThreshold)

	if numRows, err := db.Runs.Cleanup(maxHeight); err == nil {
		log.Println("runs removed:", numRows)
	} else {
		log.Println("runs cleanup error:", err)
	}

	if numRows, err := db.Heights.Cleanup(maxHeight); err == nil {
		log.Println("heights removed:", numRows)
	} else {
		log.Println("heights cleanup error:", err)
	}

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
