package pipeline

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/store"
)

type Cleanup struct {
	cfg *config.Config
	db  *store.Store
}

func NewCleanup(cfg *config.Config, db *store.Store) Cleanup {
	return Cleanup{
		cfg: cfg,
		db:  db,
	}
}

func (c Cleanup) Execute() error {
	return nil
}
