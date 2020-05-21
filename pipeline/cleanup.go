package pipeline

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/store"
)

func RunCleanup(cfg *config.Config, db *store.Store) error {
	run, _ := db.Runs.Last()
	if run == nil {
		return nil
	}

	maxHeight := uint64(run.Height) - uint64(cfg.CleanupThreshold)
	if maxHeight <= 0 {
		return nil
	}

	return db.Runs.Cleanup(maxHeight)
}
