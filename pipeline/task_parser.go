package pipeline

import (
	"context"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/store"
	"github.com/sirupsen/logrus"
)

// ParserTask performs raw block data parsing
type ParserTask struct {
	db     *store.Store
	logger *logrus.Logger
}

// NewParserTask returns a new parser task
func NewParserTask(db *store.Store, logger *logrus.Logger) ParserTask {
	return ParserTask{
		db:     db,
		logger: logger,
	}
}

// ShouldRun returns true if there any heights to process
func (t ParserTask) ShouldRun(payload *Payload) bool {
	return len(payload.Heights) > 0
}

// Name returns the task name
func (t ParserTask) Name() string {
	return parserTaskName
}

// Run executes the parser task
func (t ParserTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(t, time.Now())

	for _, h := range payload.Heights {
		parsed := &ParsedPayload{}
		h.Parsed = parsed

		block, err := mapper.Block(h.Block)
		if err != nil {
			return err
		}
		parsed.Block = block

		epoch := &model.Epoch{
			ID:              h.Block.Header.EpochID,
			StartTime:       block.Time,
			StartHeight:     uint64(block.ID),
			EndTime:         block.Time,
			EndHeight:       uint64(block.ID),
			ValidatorsCount: 0,
		}
		parsed.Epoch = epoch

		for _, v := range h.Validators {
			validator, err := mapper.Validator(h.Block, &v)
			if err != nil {
				return err
			}
			if fee, ok := h.RewardFees[v.AccountID]; ok {
				validator.RewardFee = &fee.Numerator
				rff, _ := util.Divide(fee.Numerator, fee.Denominator)
				reward := types.NewInt64Amount(0)
				reward.Mul(validator.Stake.Int, rff)
				validator.Reward = reward
			}
			parsed.Validators = append(parsed.Validators, *validator)

			validatorAgg, err := mapper.ValidatorAgg(h.Block, &v)
			if err != nil {
				return err
			}
			if fee, ok := h.RewardFees[v.AccountID]; ok {
				validatorAgg.RewardFee = &fee.Numerator
			}
			parsed.ValidatorAggs = append(parsed.ValidatorAggs, *validatorAgg)

			account, err := mapper.AccountFromValidator(h.Block, &v)
			if err != nil {
				return err
			}
			parsed.Accounts = append(parsed.Accounts, *account)

			parsed.ValidatorEpochs = append(parsed.ValidatorEpochs, model.ValidatorEpoch{
				AccountID:         validator.AccountID,
				Epoch:             validator.Epoch,
				LastHeight:        validator.Height,
				LastTime:          validator.Time,
				ExpectedBlocks:    validator.ExpectedBlocks,
				ProducedBlocks:    validator.ProducedBlocks,
				Efficiency:        validator.Efficiency,
				StakingBalance:    validator.Stake,
				RewardFee:         validator.RewardFee,
				Reward: validator.Reward,
			})
		}

		for _, v := range h.PreviousValidators {
			validator, err := mapper.Validator(h.PreviousBlock, &v)
			if err != nil {
				return err
			}

			parsed.ValidatorEpochs = append(parsed.ValidatorEpochs, model.ValidatorEpoch{
				AccountID:         validator.AccountID,
				Epoch:             validator.Epoch,
				LastHeight:        validator.Height,
				LastTime:          validator.Time,
				ExpectedBlocks:    validator.ExpectedBlocks,
				ProducedBlocks:    validator.ProducedBlocks,
				Efficiency:        validator.Efficiency,
				StakingBalance:    validator.Stake,
				RewardFee:         validator.RewardFee,
				Reward: validator.Reward,
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
		block.TransactionsCount = len(transactions)
		parsed.Block.TransactionsCount = len(transactions)
	}

	return nil
}
