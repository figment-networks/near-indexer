package cli

import (
	"log"
	"sync"
	"time"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
	"github.com/figment-networks/near-indexer/store"
)

func startSync(cfg *config.Config, db *store.Store) error {
	log.Println("sync will run every", cfg.SyncInterval)
	duration, err := time.ParseDuration(cfg.SyncInterval)
	if err != nil {
		return err
	}

	log.Println("using rpc endpoint", cfg.RPCEndpoint)
	client := near.NewClient(cfg.RPCEndpoint)

	for range time.Tick(duration) {
		if err := pipeline.RunSync(cfg, db, client); err != nil {
			log.Println("sync error:", err)
		}
	}

	return nil
}

func startCleanup(cfg *config.Config, db *store.Store) error {
	log.Println("cleanup will run every", cfg.CleanupInterval)
	duration, err := time.ParseDuration(cfg.CleanupInterval)
	if err != nil {
		return err
	}

	for range time.Tick(duration) {
		if err := pipeline.RunCleanup(cfg, db); err != nil {
			log.Println("cleanup error:", err)
		}
	}

	return nil
}

func startWorker(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := startSync(cfg, db); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	go func() {
		if err := startCleanup(cfg, db); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
