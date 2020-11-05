package pipeline

import (
	"context"
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store"
	"github.com/sirupsen/logrus"
)

type PersistorTask struct {
	db     *store.Store
	logger *logrus.Logger
}

func NewPersistorTask(db *store.Store, logger *logrus.Logger) PersistorTask {
	return PersistorTask{
		db:     db,
		logger: logger,
	}
}

func (t PersistorTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(PersistorTaskName, time.Now())

	if len(payload.Heights) == 0 {
		return nil
	}

	blocks := []model.Block{}
	transactions := []model.Transaction{}
	epochs := []model.Epoch{}
	epochIds := map[string]bool{}

	for _, h := range payload.Heights {
		if h.Skip || h.Block == nil {
			continue
		}

		blocks = append(blocks, *h.Parsed.Block)
		transactions = append(transactions, h.Parsed.Transactions...)

		if !epochIds[h.Parsed.Epoch.UUID] {
			epochIds[h.Parsed.Epoch.UUID] = true
			epochs = append(epochs, *h.Parsed.Epoch)
		}
	}

	if err := t.db.Blocks.Import(blocks); err != nil {
		return err
	}

	if err := t.db.Epochs.Import(epochs); err != nil {
		return err
	}

	if err := t.db.Transactions.Import(transactions); err != nil {
		return err
	}

	for _, h := range payload.Heights {
		if h.Parsed == nil {
			continue
		}

		if err := t.processHeight(h, h.Parsed); err != nil {
			return err
		}
	}

	if err := t.createStats(payload); err != nil {
		return err
	}

	epochsToUpdate := []string{}
	for _, h := range payload.Heights {
		if h.PreviousBlock != nil {
			epochsToUpdate = append(epochsToUpdate, h.Block.Header.EpochID, h.Block.Header.EpochID)
		}
	}
	if len(epochsToUpdate) > 0 {
		if err := t.db.Epochs.UpdateCounters(epochsToUpdate); err != nil {
			return err
		}
	}

	lastHeight := payload.Heights[len(payload.Heights)-1]
	payload.Lag = int(payload.CurrentBlock.Header.Height - uint64(lastHeight.Height))

	return nil
}

func (t PersistorTask) processHeight(h *HeightPayload, parsed *ParsedPayload) error {
	if len(parsed.ValidatorAggs) > 0 {
		t.logger.WithField("count", len(parsed.ValidatorAggs)).Debug("saving validator aggs")
		if err := t.db.ValidatorAggs.Import(parsed.ValidatorAggs); err != nil {
			return err
		}
	}

	if len(parsed.Validators) > 0 {
		t.logger.WithField("count", len(parsed.Validators)).Debug("saving validators")
		if err := t.db.Validators.Import(parsed.Validators); err != nil {
			return err
		}
	}

	if len(parsed.ValidatorEpochs) > 0 {
		t.logger.WithField("count", len(parsed.ValidatorEpochs)).Debug("saving validator epochs")
		if err := t.db.ValidatorAggs.ImportValidatorEpochs(parsed.ValidatorEpochs); err != nil {
			return err
		}
	}

	if len(parsed.Accounts) > 0 {
		t.logger.WithField("count", len(parsed.Accounts)).Debug("saving accounts")
		if err := t.db.Accounts.Import(parsed.Accounts); err != nil {
			return err
		}
	}

	if len(parsed.Events) > 0 {
		t.logger.WithField("count", len(parsed.Events)).Debug("saving events")
		for _, event := range parsed.Events {
			if err := t.db.Events.Create(&event); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t PersistorTask) createStats(payload *Payload) error {
	timeRange := store.TimeRange{
		Start: payload.StartTime,
		End:   payload.EndTime,
	}

	for _, bucket := range []string{store.BucketHour, store.BucketDay} {
		t.logger.WithField("bucket", bucket).Debug("creating block stats")
		if err := t.db.Stats.CreateBlockStats(bucket, timeRange); err != nil {
			return err
		}
	}

	return nil
}
