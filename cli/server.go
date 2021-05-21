package cli

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strings"

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
	if len(rpcEndpoints) != 1 {
		return errors.New("only one rpc should be set for status command")
	}
	rpc := near.DefaultClient(rpcEndpoints[0])

	srv := server.New(cfg, db, logger, rpc)

	logger.Info("Starting server on ", cfg.ListenAddr())
	return srv.Run(cfg.ListenAddr())
}
