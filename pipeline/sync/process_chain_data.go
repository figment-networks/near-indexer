package sync

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
)

func ProcessChainData(c *Context) {
	if err := processBlockData(c); err != nil {
		c.Abort(err)
		return
	}

	if err := processValidatorsData(c); err != nil {
		c.Abort(err)
		return
	}
}

func processBlockData(c *Context) error {
	record, err := mapper.Block(c.Block)
	if err != nil {
		return err
	}
	record.AppVersion = c.Status.Version.String()

	// TODO: Do not remove the block
	if err := c.DB.Blocks.DeleteByHeight(c.BlockHeight); err != nil {
		return err
	}
	c.BlockTime = record.Time

	return c.DB.Blocks.Create(record)
}

func processValidatorsData(c *Context) error {
	validators := c.Validators
	if len(validators) == 0 {
		return nil
	}

	validatorRecords := make([]model.Validator, len(validators))
	validatorEpochRecords := make([]model.ValidatorEpoch, len(validators))
	validatorAggRecords := make([]model.ValidatorAgg, len(validators))
	accountRecords := make([]model.Account, len(validators))

	for idx, v := range validators {
		// Prepare a new validator record
		validator, err := mapper.Validator(c.Block, &v)
		if err != nil {
			return err
		}
		validatorRecords[idx] = *validator

		// Prepare a new validator aggregate record
		validatorAgg, err := mapper.ValidatorAgg(c.Block, &v)
		if err != nil {
			return err
		}
		validatorAggRecords[idx] = *validatorAgg

		// Prepare a new validator epoch record
		validatorEpochRecords[idx] = model.ValidatorEpoch{
			AccountID:      validator.AccountID,
			Epoch:          validator.Epoch,
			LastHeight:     validator.Height,
			LastTime:       validator.Time,
			ExpectedBlocks: validator.ExpectedBlocks,
			ProducedBlocks: validator.ProducedBlocks,
			Efficiency:     validator.Efficiency,
		}

		// Prepare a new account from the validator details
		account, err := mapper.AccountFromValidator(c.Block, &v)
		if err != nil {
			return err
		}
		accountRecords[idx] = *account
	}

	if err := c.DB.Validators.BulkInsert(validatorRecords); err != nil {
		return err
	}
	if err := c.DB.ValidatorAggs.ImportValidatorEpochs(validatorEpochRecords); err != nil {
		return err
	}
	if err := c.DB.ValidatorAggs.BulkUpsert(validatorAggRecords); err != nil {
		return err
	}
	if err := c.DB.ValidatorAggs.UpdateCountsForHeight(c.BlockHeight); err != nil {
		return err
	}
	if err := c.DB.Accounts.BulkUpsert(accountRecords); err != nil {
		return err
	}

	return nil
}
