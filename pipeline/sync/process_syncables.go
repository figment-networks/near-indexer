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

	for _, v := range validators {
		validator, err := mapper.Validator(c.Block, &v)
		if err != nil {
			return err
		}

		validatorAgg, err := mapper.ValidatorAgg(c.Block, &v)
		if err != nil {
			return err
		}

		account, err := mapper.AccountFromValidator(c.Block, &v)
		if err != nil {
			return err
		}

		if err := c.DB.ValidatorAggs.Upsert(validatorAgg); err != nil {
			return err
		}

		if err := c.DB.Validators.Create(&validator); err != nil {
			return err
		}

		if err := c.DB.Accounts.Upsert(account); err != nil {
			return err
		}
	}

	return nil
}
