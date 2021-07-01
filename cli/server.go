package cli

import (
	"strings"

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

	rpcEndpoints := strings.Split(cfg.RPCEndpoints, ",")
	rpc := near.DefaultClient(rpcEndpoints[0])

	srv := server.New(cfg, db, logger, rpc, logger)

	logger.Info("Starting server on ", cfg.ListenAddr())
	return srv.Run(cfg.ListenAddr())
}
