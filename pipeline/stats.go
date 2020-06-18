package pipeline

import (
	"github.com/figment-networks/near-indexer/store"
)

func RunStats(db *store.Store) error {
	block, err := db.Blocks.Recent()
	if err != nil {
		return err
	}
	t := block.Time

	if err := db.Stats.CreateBlockStats(store.BucketHour, t); err != nil {
		return err
	}

	if err := db.Stats.CreateBlockStats(store.BucketDay, t); err != nil {
		return err
	}

	if err := db.Stats.CreateValidatorsStats(store.BucketHour, t); err != nil {
		return err
	}

	if err := db.Stats.CreateValidatorsStats(store.BucketDay, t); err != nil {
		return err
	}

	return nil
}
