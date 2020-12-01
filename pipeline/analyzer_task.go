package pipeline

import (
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

type AnalyzerTask struct {
	db     *store.Store
	logger *logrus.Logger
}

func NewAnalyzerTask(db *store.Store, logger *logrus.Logger) AnalyzerTask {
	return AnalyzerTask{
		db:     db,
		logger: logger,
	}
}

func (t AnalyzerTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(AnalyzerTaskName, time.Now())

	// Skip processing if no height data is available
	if len(payload.Heights) == 0 {
		return nil
	}

	for _, h := range payload.Heights {
		// Do not process any of events until epoch is complete
		if h.CurrentEpoch {
			continue
		}

		numPrev := len(h.PreviousValidators)
		numCurrent := len(h.Validators)

		if numPrev == 0 && numCurrent == 0 {
			continue
		}

		t.logger.WithFields(logrus.Fields{
			"previous": numPrev,
			"current":  numCurrent,
			"diff":     numCurrent - numPrev,
		}).Info("validator set has changed")

		previousIds := map[string]bool{}
		currentIds := map[string]bool{}

		for _, v := range h.PreviousValidators {
			previousIds[v.AccountID] = true
		}

		for _, v := range h.Validators {
			currentIds[v.AccountID] = true
			if !previousIds[v.AccountID] {
				t.logger.WithField("account", v.AccountID).Info("validator added to active set")
				if err := t.createValidatorAddEvent(h.Block, &v); err != nil {
					return err
				}
			}
		}

		for _, v := range h.PreviousValidators {
			if !currentIds[v.AccountID] {
				t.logger.WithField("account", v.AccountID).Info("validator removed from active set")
				if err := t.createValidatorRemoveEvent(h.Block, &v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t AnalyzerTask) createValidatorAddEvent(block *near.Block, validator *near.Validator) error {
	event, err := mapper.ValidatorAddEvent(block, validator)
	if err == nil {
		err = t.db.Events.Create(event)
	}
	return err
}

func (t AnalyzerTask) createValidatorRemoveEvent(block *near.Block, validator *near.Validator) error {
	event, err := mapper.ValidatorRemoveEvent(block, validator)
	if err == nil {
		err = t.db.Events.Create(event)
	}
	return err
}
