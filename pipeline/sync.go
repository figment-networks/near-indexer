package pipeline

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline/sync"
	"github.com/figment-networks/near-indexer/store"
)

type Stage struct {
	Name    string
	Handler sync.HandlerFunc
}

func RunSync(cfg *config.Config, db *store.Store, client *near.Client) error {
	ctx := sync.NewContext(db, client)

	stages := []Stage{
		{"create_height", sync.CreateHeight},
		{"create_run", sync.CreateRun},
		{"create_syncables", sync.CreateSyncables},
		{"process_syncables", sync.ProcessSyncables},
		{"finish_run", sync.FinishRun},
	}

	for _, stage := range stages {
		stage.Handler(ctx)

		if ctx.IsAborted() {
			if ctx.LastError() != nil {
				log.Printf("aborted on stage %s with error: %s", stage.Name, ctx.LastError())
			}
			sync.FinishRun(ctx)
			break
		}
	}

	return ctx.LastError()
}
