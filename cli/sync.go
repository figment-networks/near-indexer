package cli

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
)

func runSync(cfg *config.Config) error {
	client := near.NewClient(cfg.RPCEndpoint)
	client.SetDebug(true)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return pipeline.RunSync(cfg, db, client)
}
