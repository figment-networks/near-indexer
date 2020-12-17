package pipeline

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

func RunSync(cfg *config.Config, db *store.Store, client near.Client) (int, error) {
	var err error
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
		fields := logrus.Fields{
			"duration": time.Since(startTime).Milliseconds(),
			"from":     payload.StartHeight,
			"to":       payload.EndHeight,
			"lag":      payload.Lag,
			"heights":  len(payload.Heights),
		}

		if err != nil {
			fields["error"] = err
		}

		logrus.WithFields(fields).Info("sync finished")
	}()

	fetcherTask := NewFetcherTask(db, client, cfg, logger)
	parserTask := NewParserTask(db, logger)
	persistorTask := NewPersistorTask(db, logger)
	analyzerTask := NewAnalyzerTask(db, logger)

	tasks := []Task{
		fetcherTask,
		parserTask,
		persistorTask,
		analyzerTask,
	}

	for _, t := range tasks {
		if !t.ShouldRun(payload) {
			logger.WithField("task", t.Name()).Info("task execution skipped")
			continue
		}

		if err = t.Run(context.Background(), payload); err != nil {
			logger.
				WithError(err).
				WithField("task", t.Name()).
				Error("task execution failed")

			break
		}
	}

	return payload.Lag, err
}
