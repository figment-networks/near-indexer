package pipeline

import (
	"log"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline/sync"
	"github.com/figment-networks/near-indexer/store"
)

type Task struct {
	Name    string
	Handler sync.HandlerFunc
}

func RunSync(cfg *config.Config, db *store.Store, client *near.Client) (uint64, error) {
	ctx := sync.NewContext(db, client)
	ctx.DefaultStartHeight = cfg.StartHeight

	tasks := []Task{
		{"create_height", sync.CreateHeight},
		{"create_run", sync.CreateRun},
		{"fetch_data", sync.FetchChainData},
		{"process_data", sync.ProcessChainData},
		{"finish_height", sync.FinishHeight},
	}

	for _, task := range tasks {
		task.Handler(ctx)

		if ctx.IsAborted() {
			if ctx.LastError() != nil {
				log.Printf("aborted on %s with error: %s", task.Name, ctx.LastError())
			}
			sync.FinishHeight(ctx)
			break
		}
	}

	return ctx.Lag, ctx.LastError()
}
