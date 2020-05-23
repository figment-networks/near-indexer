package cli

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
	"github.com/figment-networks/near-indexer/store"
)

func startSyncWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	client := near.NewClient(cfg.RPCEndpoint)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-time.Tick(cfg.SyncDuration()):
				pipeline.RunSync(cfg, db, client)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startCleanupWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()

		for {
			select {
			case <-time.Tick(cfg.CleanupDuration()):
				pipeline.RunCleanup(cfg, db)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startWorker(cfg *config.Config) error {
	log.Println("using rpc endpoint", cfg.RPCEndpoint)
	log.Println("sync will run every", cfg.SyncInterval)
	log.Println("cleanup will run every", cfg.CleanupInterval)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}

	cancelSync := startSyncWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db)

	s := <-initSignals()

	log.Println("received signal", s)
	cancelSync()
	cancelCleanup()

	wg.Wait()
	return nil
}
