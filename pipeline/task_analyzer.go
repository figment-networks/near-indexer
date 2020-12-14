package pipeline

import (
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

// AnalyzerTask performs data analysis on parsed heights data
type AnalyzerTask struct {
	db     *store.Store
	logger *logrus.Logger
}

// NewAnalyzerTask returns a new analyzer task
func NewAnalyzerTask(db *store.Store, logger *logrus.Logger) AnalyzerTask {
	return AnalyzerTask{
		db:     db,
		logger: logger,
	}
}

// Name returns the task name
func (t AnalyzerTask) Name() string {
	return analyzerTaskName
}

// ShouldRun returns true if there any heights to process
func (t AnalyzerTask) ShouldRun(payload *Payload) bool {
	return len(payload.Heights) > 0
}

// Run executes the analyzer task
func (t AnalyzerTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(t, time.Now())

	// Skip processing if no height data is available
	if len(payload.Heights) == 0 {
		return nil
	}

	events := []model.Event{}

	for _, h := range payload.Heights {
		numPrev := len(h.PreviousValidators)
		numCurrent := len(h.Validators)

		// Do not process any of events until epoch is complete
		if h.CurrentEpoch {
			continue
		}

		// No previous or current validators
		if numPrev == 0 && numCurrent == 0 {
			continue
		}

		t.logger.WithFields(logrus.Fields{
			"count_before": numPrev,
			"count_after":  numCurrent,
			"diff":         numCurrent - numPrev,
		}).Info("validator set changed")

		previousIds := map[string]*near.Validator{}
		currentIds := map[string]*near.Validator{}

		for _, v := range h.PreviousValidators {
			previousIds[v.AccountID] = &v
		}

		for _, v := range h.Validators {
			currentIds[v.AccountID] = &v

			if previousIds[v.AccountID] == nil {
				t.logger.
					WithField("account", v.AccountID).
					WithField("height", h.Height).
					Info("validator added to active set")

				event, err := mapper.ValidatorAddEvent(h.Block, &v)
				if err != nil {
					return err
				}
				events = append(events, *event)
			}

			// if previousIds[v.AccountID] != nil {
			// 	prevAmount := types.NewAmount(previousIds[v.AccountID].Stake)
			// 	curAmount := types.NewAmount(v.Stake)

			// 	t.logger.
			// 		WithField("account", v.AccountID).
			// 		WithField("before", prevAmount).
			// 		WithField("after", curAmount).
			// 		Info("validator staking balance change")

			// 	if curAmount.Compare(prevAmount) != 0 {
			// 		event, err := mapper.ValidatorStakingBalanceChangeEvent(h.Block, &v, &prevAmount, &curAmount)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		log.Println(event)
			// 		//events = append(events, *event)
			// 	}
			// }
		}

		for _, v := range h.PreviousValidators {
			if currentIds[v.AccountID] != nil {
				continue
			}
			t.logger.WithField("account", v.AccountID).Info("validator removed from active set")

			event, err := mapper.ValidatorRemoveEvent(h.Block, &v)
			if err != nil {
				return err
			}
			events = append(events, *event)
		}
	}

	return t.db.Events.Import(events)
}
