package cli

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/pipeline"
)

func startCleanup(cfg *config.Config, logger *logrus.Logger) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return pipeline.RunCleanup(cfg, db, logger)
}
