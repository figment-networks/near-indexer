package sync

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/figment-networks/near-indexer/model"
)

func CreateSyncables(c *Context) {
	log.Println("creating syncables for height", c.Height.Height)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go createBlockSyncable(c, wg)
	go createValidatorsSyncable(c, wg)

	wg.Wait()
}

func createBlockSyncable(c *Context, wg *sync.WaitGroup) {
	defer wg.Done()

	block, err := c.Client.BlockByHeight(c.BlockHeight)
	if err != nil {
		c.Abort(err)
		return
	}
	c.Block = &block

	if err := createSyncable(c, model.SyncableTypeBlock, block); err != nil {
		c.Abort(err)
	}
}

func createValidatorsSyncable(c *Context, wg *sync.WaitGroup) {
	defer wg.Done()

	validators, err := c.Client.ValidatorsByHeight(c.BlockHeight)
	if err != nil {
		c.Abort(err)
		return
	}

	if err := createSyncable(c, model.SyncableTypeValidators, validators); err != nil {
		c.Abort(err)
	}
}

func createSyncable(c *Context, kind string, data interface{}) error {
	jsondata, err := json.Marshal(data)
	if err != nil {
		return err
	}

	syncable := model.Syncable{
		RunID:  c.Run.ID,
		Height: c.Height.Height,
		Time:   time.Now(),
		Type:   kind,
		Data:   jsondata,
	}

	if err := c.DB.Syncables.Create(&syncable); err != nil {
		return err
	}
	c.Syncables = append(c.Syncables, syncable)

	return nil
}
