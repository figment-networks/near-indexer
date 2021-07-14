package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
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

func RunSyncHistoricalDelegators(cfg *config.Config, db *store.Store, clients []near.Client, logger *logrus.Logger) error {
	fetcherTask := NewFetcherTask(db, clients, cfg, logger)

	epochs, err := db.Epochs.FindUnIndexedDelegatorsEpochs()
	if err != nil {
		return err
	}
	logger.WithField("count", len(epochs)).Info("found historical epochs to sync")

	for _, epoch := range epochs {
		logger.WithField("epoch", epoch.ID).Info("fetching historical delegators")

		validatorEpochs, err := db.ValidatorAggs.FindValidatorEpochsByEpoch(epoch.ID)
		if err != nil {
			return err
		}
		accounts := make([]string, len(validatorEpochs))
		for idx, validator := range validatorEpochs {
			accounts[idx] = validator.AccountID
		}

		results := []delegationsFetchResult{}
		resultsLock := &sync.Mutex{}

		doConcurrently(accounts, fetcherTask.concurrencyLevel, func(account string) {
			var dlgs []near.AccountInfo
			for i := 1; i <= fetcherTask.maxRetryCount; i++ {
				dlgs, err = fetcherTask.RPC().Delegations(account, epoch.EndHeight)
				if err == nil {
					break
				}
				fetcherTask.logger.WithError(err).WithField("validator", account).Error("can not fetch delegations")
			}
			resultsLock.Lock()
			defer resultsLock.Unlock()

			results = append(results, delegationsFetchResult{
				account:     account,
				delegations: dlgs,
				err:         err,
			})
		})

		delegationsByValidator := map[string][]near.AccountInfo{}
		for _, res := range results {
			if res.err != nil {
				fetcherTask.logger.WithError(res.err).WithField("validator", res.account).Error("can not fetch delegations ")
				return res.err
			}
			delegationsByValidator[res.account] = res.delegations
		}
		var delegatorEpochs []model.DelegatorEpoch
		for _, ve := range validatorEpochs {
			if delegations, ok := delegationsByValidator[ve.AccountID]; ok {
				for _, d := range delegations {
					de := model.DelegatorEpoch{
						AccountID:       d.Account,
						ValidatorID:     ve.AccountID,
						Epoch:           epoch.ID,
						StakedBalance:   types.NewAmount(d.StakedBalance),
						UnstakedBalance: types.NewAmount(d.UnstakedBalance),
					}
					delegatorEpochs = append(delegatorEpochs, de)
				}
			}
		}
		err = db.Delegators.ImportDelegatorEpochs(delegatorEpochs)
		if err != nil {
			return err
		}
	}
	logger.Info(" finishing fetching historical delegators")

	return nil
}
