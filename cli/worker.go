package cli

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/pipeline"
	"github.com/figment-networks/near-indexer/store"
)

func startSyncWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(cfg.SyncDuration())
	busy := false

	rpcEndpoints := strings.Split(cfg.RPCEndpoints, ",")
	clients := []near.Client{}
	for _, address := range rpcEndpoints {
		client := near.DefaultClient(address)
		client.SetDebug(cfg.Debug)
		client.SetTimeout(cfg.RPCClientTimeout())
		clients = append(clients, client)
	}

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
					lag, _ := pipeline.RunSync(cfg, db, clients)
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

func startCleanupWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store, logger *logrus.Logger) context.CancelFunc {
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
				pipeline.RunCleanup(cfg, db, logger)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startWorker(cfg *config.Config, logger *logrus.Logger) error {
	logger.Info("log level: ", cfg.LogLevel)
	logger.Info("using rpc endpoints: ", cfg.RPCEndpoints)
	logger.Info("sync will run every: ", cfg.SyncInterval)
	logger.Info("cleanup will run every: ", cfg.CleanupInterval)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	cancelSync := startSyncWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db, logger)

	s := <-initSignals()
	logger.Info("received signal: ", s)

	cancelSync()
	cancelCleanup()

	wg.Wait()
	return nil
}
