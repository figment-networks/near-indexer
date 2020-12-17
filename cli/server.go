package cli

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/server"
)

func startServer(cfg *config.Config, logger *logrus.Logger) error {
	server.SetGinDefaults(cfg)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	rpc := near.DefaultClient(cfg.RPCEndpoint)

	srv := server.New(cfg, db, logger, rpc)

	logger.Info("Starting server on ", cfg.ListenAddr())
	return srv.Run(cfg.ListenAddr())
}
