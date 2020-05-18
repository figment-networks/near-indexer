package cli

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/server"
)

func startServer(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	client := near.NewClient(cfg.RPCEndpoint)
	srv := server.New(db, &client)

	log.Println("Starting server on", cfg.ListenAddr())
	return srv.Run(cfg.ListenAddr())
}
