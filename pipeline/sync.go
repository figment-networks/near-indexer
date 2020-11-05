package pipeline

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

func RunSync(cfg *config.Config, db *store.Store, client *near.Client) (int, error) {
	startTime := time.Now()
	payload := &Payload{}
	logger := logrus.StandardLogger()

	switch cfg.LogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	defer func() {
		logrus.
			WithField("duration", time.Since(startTime).Milliseconds()).
			WithField("from", payload.StartHeight).
			WithField("to", payload.EndHeight).
			WithField("lag", payload.Lag).
			WithField("heights", len(payload.Heights)).
			Info("sync finished")
	}()

	fetcherTask := NewFetcherTask(db, client, cfg, logger)
	parserTask := NewParserTask(db, logger)
	persistorTask := NewPersistorTask(db, logger)
	analyzerTask := NewAnalyzerTask(db, logger)

	tasks := []func(context.Context, *Payload) error{
		fetcherTask.Run,
		parserTask.Run,
		persistorTask.Run,
		analyzerTask.Run,
	}

	var err error

	for _, task := range tasks {
		if err = task(context.Background(), payload); err != nil {
			logger.WithError(err).Error("task failed with error")
			break
		}
	}

	return payload.Lag, err
}
