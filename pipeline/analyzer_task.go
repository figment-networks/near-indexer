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
		if h.CurrentEpoch {
			continue
		}
		if len(h.PreviousValidators) == 0 && len(h.Validators) == 0 {
			continue
		}

		t.logger.WithFields(logrus.Fields{
			"previous": len(h.PreviousValidators),
			"current":  len(h.Validators),
		}).Info("validator set changed")

		previousIds := map[string]bool{}
		currentIds := map[string]bool{}

		for _, v := range h.PreviousValidators {
			previousIds[v.AccountID] = true
		}

		// Find new validators in the set
		for _, v := range h.Validators {
			currentIds[v.AccountID] = true
			if !previousIds[v.AccountID] {
				t.logger.WithField("account", v.AccountID).Info("validator added to active set")
				if err := t.createValidatorAddEvent(h.Block, &v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t AnalyzerTask) createValidatorAddEvent(block *near.Block, validator *near.Validator) error {
	event, err := mapper.ValidatorAddEvent(block, validator)
	if err != nil {
		return err
	}
	return t.db.Events.Create(event)
}
