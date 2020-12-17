package pipeline

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/store"
)

// RunCleanup performs the data cleanup
func RunCleanup(cfg *config.Config, db *store.Store, logger *logrus.Logger) error {
	lastBlock, err := db.Blocks.Last()
	if err != nil {
		return err
	}

	logger.WithField("height", lastBlock.ID).Info("starting cleanup")

	keepHeights := uint64(1000)
	if numRows, err := db.Validators.Cleanup(keepHeights); err == nil {
		logrus.WithField("count", numRows).Info("validators removed")
	} else {
		logrus.WithError(err).Error("validators cleanup failed")
	}

	return nil
}
