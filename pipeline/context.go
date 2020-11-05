package pipeline

import (
	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

// Context contains sync data and clients
type Context struct {
	rpc *near.Client
	db  *store.Store

	batchSize          int
	defaultStartHeight uint64

	heightPayloads []*HeightPayload
}

// NewContext returns a new sync context
func NewContext(db *store.Store, rpc *near.Client, config *config.Config) *Context {
	return &Context{
		db:                 db,
		rpc:                rpc,
		batchSize:          config.SyncBatchSize,
		defaultStartHeight: config.StartHeight,
	}
}
