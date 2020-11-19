package cli

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
)

func runSync(cfg *config.Config) error {
	client := near.DefaultClient(cfg.RPCEndpoint)
	client.SetDebug(true)
	client.SetTimeout(cfg.RPCClientTimeout())

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = pipeline.RunSync(cfg, db, client)
	return err
}
