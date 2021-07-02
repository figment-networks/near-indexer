package pipeline

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

func RunSync(cfg *config.Config, db *store.Store, clients []near.Client) (int, error) {
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

	fetcherTask := NewFetcherTask(db, clients, cfg, logger)
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

func RunSyncHistoricalDelegators(cfg *config.Config, db *store.Store, clients []near.Client) error {
	logger := logrus.StandardLogger()
	fetcherTask := NewFetcherTask(db, clients, cfg, logger)

	for {
		transactions, err := db.Transactions.FindUnIndexedTransactionFees()
		if err != nil {
			return err
		}
		if len(transactions) == 0 {
			break
		}
		var hashes []string
		for _, t := range transactions {
			hashes = append(hashes, t.Hash)
		}

		trxs, err := fetcherTask.fetchBlockTransactions(hashes)
		if err != nil {
			return err
		}
		for _, trx := range trxs {
			t := &model.Transaction{
				Hash:      trx.Transaction.Hash,
				Signature: trx.Transaction.Signature,
				PublicKey: trx.Transaction.PublicKey,
			}

			fee, err := util.CalculateTransactionFee(trx)
			if err != nil {
				return err
			}
			t.Fee = fee

			outcome, err := json.Marshal(trx.TransactionOutcome)
			if err != nil {
				return err
			}
			t.Outcome = outcome

			receipt, err := json.Marshal(trx.ReceiptsOutcome)
			if err != nil {
				return err
			}
			t.Receipt = receipt

			err = db.Transactions.UpdateTransactionsHistoricalInfo(*t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
