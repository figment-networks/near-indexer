package pipeline

import (
	"context"
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/store"
	"github.com/sirupsen/logrus"
)

type ParserTask struct {
	db     *store.Store
	logger *logrus.Logger
}

func NewParserTask(db *store.Store, logger *logrus.Logger) ParserTask {
	return ParserTask{
		db:     db,
		logger: logger,
	}
}

func (t ParserTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(ParserTaskName, time.Now())

	if len(payload.Heights) == 0 {
		return nil
	}

	for _, h := range payload.Heights {
		if h.Block == nil {
			continue
		}

		parsed := &ParsedPayload{}
		h.Parsed = parsed

		block, err := mapper.Block(h.Block)
		if err != nil {
			return err
		}
		parsed.Block = block

		epoch := &model.Epoch{
			UUID:            h.Block.Header.EpochID,
			StartTime:       block.Time,
			StartHeight:     uint64(block.Height),
			EndTime:         block.Time,
			EndHeight:       uint64(block.Height),
			ValidatorsCount: 0,
		}
		parsed.Epoch = epoch

		for _, v := range h.Validators {
			validator, err := mapper.Validator(h.Block, &v)
			if err != nil {
				return err
			}
			parsed.Validators = append(parsed.Validators, *validator)

			validatorAgg, err := mapper.ValidatorAgg(h.Block, &v)
			if err != nil {
				return err
			}
			parsed.ValidatorAggs = append(parsed.ValidatorAggs, *validatorAgg)

			account, err := mapper.AccountFromValidator(h.Block, &v)
			if err != nil {
				return err
			}
			parsed.Accounts = append(parsed.Accounts, *account)

			parsed.ValidatorEpochs = append(parsed.ValidatorEpochs, model.ValidatorEpoch{
				AccountID:      validator.AccountID,
				Epoch:          validator.Epoch,
				LastHeight:     validator.Height,
				LastTime:       validator.Time,
				ExpectedBlocks: validator.ExpectedBlocks,
				ProducedBlocks: validator.ProducedBlocks,
				Efficiency:     validator.Efficiency,
			})
		}

		for _, v := range h.PreviousValidators {
			validator, err := mapper.Validator(h.PreviousBlock, &v)
			if err != nil {
				return err
			}

			parsed.ValidatorEpochs = append(parsed.ValidatorEpochs, model.ValidatorEpoch{
				AccountID:      validator.AccountID,
				Epoch:          validator.Epoch,
				LastHeight:     validator.Height,
				LastTime:       validator.Time,
				ExpectedBlocks: validator.ExpectedBlocks,
				ProducedBlocks: validator.ProducedBlocks,
				Efficiency:     validator.Efficiency,
			})
		}

		if len(h.PreviousEpochKickOut) > 0 {
			for _, kick := range h.PreviousEpochKickOut {
				event, err := mapper.ValidatorKickoutEvent(h.Block, &kick)
				if err != nil {
					return err
				}
				parsed.Events = append(parsed.Events, *event)
			}
		}

		transactions, err := mapper.Transactions(h.Block, h.Transactions)
		if err != nil {
			t.logger.
				WithError(err).
				WithField("block", h.Block.Header.Height).
				Error(err)

			return err
		}
		parsed.Transactions = append(parsed.Transactions, transactions...)
	}

	return nil
}
