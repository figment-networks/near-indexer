package cli

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/server"
)

func startServer(cfg *config.Config) error {
	server.SetGinDefaults(cfg)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := server.New(db)

	log.Println("Starting server on", cfg.ListenAddr())
	return srv.Run(cfg.ListenAddr())
}
