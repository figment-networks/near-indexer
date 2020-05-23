package sync

import (
	"log"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/near"
)

func ProcessSyncables(c *Context) {
	var err error

	for _, syncable := range c.Syncables {
		switch syncable.Type {
		case model.SyncableTypeBlock:
			err = processBlockSyncable(c, &syncable)
		case model.SyncableTypeValidators:
			err = processValidatorsSyncable(c, &syncable)
		}

		if err == nil {
			err = c.DB.Syncables.MarkProcessed(&syncable)
		}

		if err != nil {
			log.Println("syncable error:", err)
			c.Abort(err)
			return
		}
	}
}

func processBlockSyncable(c *Context, syncable *model.Syncable) error {
	block := near.Block{}
	if err := syncable.Decode(&block); err != nil {
		return err
	}

	record, err := mapper.Block(&block)
	if err != nil {
		return err
	}
	record.AppVersion = c.Status.Version.String()

	if err := c.DB.Blocks.DeleteByHeight(c.BlockHeight); err != nil {
		return err
	}

	return c.DB.Blocks.Create(record)
}

func processValidatorsSyncable(c *Context, syncable *model.Syncable) error {
	validators := []near.Validator{}
	if err := syncable.Decode(&validators); err != nil {
		return err
	}
	if len(validators) == 0 {
		return nil
	}

	validatorRecords := make([]model.Validator, len(validators))
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
	if err := c.DB.ValidatorAggs.BulkUpsert(validatorAggRecords); err != nil {
		return err
	}
	if err := c.DB.Accounts.BulkUpsert(accountRecords); err != nil {
		return err
	}

	return nil
}
