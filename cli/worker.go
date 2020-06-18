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
	timer := time.NewTimer(cfg.SyncDuration())

	client := near.NewClient(cfg.RPCEndpoint)
	client.SetDebug(cfg.Debug)

	go func() {
		defer func() {
			timer.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-timer.C:
				lag, _ := pipeline.RunSync(cfg, db, client)
				if lag > 60 {
					timer = time.NewTimer(time.Millisecond * 50)
				} else {
					timer = time.NewTimer(cfg.SyncDuration())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startStatsWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second)

	go func() {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case t := <-ticker.C:
				if t.Second() == 0 {
					if err := pipeline.RunStats(db); err != nil {
						log.Println("stats error:", err)
					}
				}
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
	ticker := time.NewTicker(cfg.CleanupDuration())

	go func() {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ticker.C:
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
	cancelStats := startStatsWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db)

	s := <-initSignals()

	log.Println("received signal", s)
	cancelSync()
	cancelStats()
	cancelCleanup()

	wg.Wait()
	return nil
}
