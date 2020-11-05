package cli

import (
	"context"
	"sync"
	"time"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
	"github.com/figment-networks/near-indexer/store"
	"github.com/sirupsen/logrus"
)

func startSyncWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(cfg.SyncDuration())
	busy := false

	client := near.DefaultClient(cfg.RPCEndpoint)
	client.SetDebug(cfg.Debug)

	go func() {
		defer func() {
			timer.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-timer.C:
				if !busy {
					busy = true
					lag, _ := pipeline.RunSync(cfg, db, client)
					busy = false
					if lag > 60 {
						timer = time.NewTimer(time.Millisecond * 10)
					} else {
						timer = time.NewTimer(cfg.SyncDuration())
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
	logrus.Info("log level:", cfg.LogLevel)
	logrus.Info("using rpc endpoint: ", cfg.RPCEndpoint)
	logrus.Info("sync will run every: ", cfg.SyncInterval)
	logrus.Info("cleanup will run every: ", cfg.CleanupInterval)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}

	cancelSync := startSyncWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db)

	s := <-initSignals()

	logrus.Info("received signal: ", s)
	cancelSync()
	cancelCleanup()

	wg.Wait()
	return nil
}
