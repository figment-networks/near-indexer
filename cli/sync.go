package cli

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/pipeline"
)

func runSync(cfg *config.Config, logger *logrus.Logger) error {
	clients := initClients(cfg)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = pipeline.RunSync(cfg, db, clients)
	return err
}

func runSyncHistoricalDelegators(cfg *config.Config, logger *logrus.Logger) error {
	clients := initClients(cfg)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	err = pipeline.RunSyncHistoricalDelegators(cfg, db, clients)
	return err
}
