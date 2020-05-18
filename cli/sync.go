package cli

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
)

func runSync(cfg *config.Config) error {
	client := near.NewClient(cfg.RPCEndpoint)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	runner := pipeline.NewSync(cfg, db, &client)
	return runner.Execute()
}
