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
		return nil
	}

	log.Println("removing run records")
	if err := db.Runs.Cleanup(maxHeight); err != nil {
		log.Println("runs cleanup error:", err)
	}

	log.Println("removing validator records")
	if err := db.Validators.Cleanup(maxHeight); err != nil {
		log.Println("validators cleanup error:", err)
	}

	if err := db.Validators.CleanupCounts(); err != nil {
		log.Println("validator counts cleanup error:", err)
	}

	return nil
}
