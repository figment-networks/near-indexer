package pipeline

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

type Sync struct {
	cfg    *config.Config
	db     *store.Store
	client *near.Client

	latestHash   string
	latestHeight uint64
	height       uint64
	block        *near.Block
	blockTime    time.Time
	blockHash    string

	errors    []error
	syncables []model.Syncable

	run *model.Run
}

func NewSync(cfg *config.Config, db *store.Store, client *near.Client) Sync {
	return Sync{
		db:     db,
		client: client,

		errors:    []error{},
		syncables: []model.Syncable{},
	}
}

func (s *Sync) Execute() error {
	return runChain(
		func() error { return s.getCurrentHeight() },
		func() error { return s.createRun() },
		func() error { return s.createSyncables() },
		func() error { return s.processSyncables() },
		func() error { return s.finishRun() },
	)
}

func (s *Sync) getCurrentHeight() error {
	log.Println("fetching current height")
	status, err := s.client.Status()
	if err != nil {
		return err
	}
	log.Println("latest height is:", status.SyncInfo.LatestBlockHeight)

	s.latestHash = status.SyncInfo.LatestBlockHash
	s.latestHeight = status.SyncInfo.LatestBlockHeight

	return nil
}

func (s *Sync) createRun() error {
	log.Println("fetching most recent height")
	height, err := s.db.Syncables.Height()
	if err != nil {
		return err
	}
	if height > 0 {
		log.Println("recent indexed height is:", height)
	}

	run := &model.Run{}

	if height == 0 {
		run.Height = s.latestHeight
	} else {
		run.Height = height + 1
	}
	s.run = run
	s.height = run.Height

	log.Printf("creating run for height=%v lag=%v\n", s.height, s.latestHeight-s.height)

	return s.db.Runs.Create(s.run)
}

func (s *Sync) createSyncables() error {
	log.Println("creating syncables for height", s.height)

	err := runChain(
		func() error { return s.createBlockSyncable() },
		func() error { return s.createValidatorsSyncable() },
	)
	if err != nil {
		return err
	}

	if len(s.syncables) == 0 {
		s.createSyncable("empty", true)
	}

	return nil
}

func (s *Sync) createBlockSyncable() error {
	block, err := s.client.BlockByHeight(s.height)
	if err != nil {
		// Workaround for server errors
		log.Println("got error:", err)
		if shouldIgnoreError(err) {
			return nil
		}
		return err
	}
	s.block = &block
	s.blockTime = time.Unix(0, block.Header.Timestamp)
	s.blockHash = block.Header.Hash

	return s.createSyncable(model.SyncableTypeBlock, block)
}

func (s *Sync) createValidatorsSyncable() error {
	validators, err := s.client.ValidatorsByHeight(s.height)
	if err != nil {
		log.Println("got error:", err)
		if shouldIgnoreError(err) {
			return nil
		}
		return err
	}
	return s.createSyncable(model.SyncableTypeValidators, validators)
}

func (s *Sync) createSyncable(kind string, data interface{}) error {
	jsondata, err := json.Marshal(data)
	if err != nil {
		return err
	}

	syncable := model.Syncable{
		RunID:  s.run.ID,
		Height: s.height,
		Time:   time.Now(),
		Type:   kind,
		Data:   jsondata,
	}

	if err := s.db.Syncables.Create(&syncable); err != nil {
		return err
	}

	s.syncables = append(s.syncables, syncable)
	return nil
}

func (s *Sync) finishRun() error {
	s.run.Success = true
	s.run.Duration = time.Since(s.run.CreatedAt).Milliseconds()
	return s.db.Runs.Update(s.run)
}

func (s *Sync) processSyncables() (err error) {
	for _, syncable := range s.syncables {
		log.Println("processing syncable", syncable)

		switch syncable.Type {
		case model.SyncableTypeBlock:
			err = s.processBlockSyncable(&syncable)
		case model.SyncableTypeValidators:
			err = s.processValidatorsSyncable(&syncable)
		}

		if err != nil {
			return
		}

		err = s.db.Syncables.MarkProcessed(&syncable)
	}

	return
}

func (s *Sync) processBlockSyncable(syncable *model.Syncable) error {
	block := near.Block{}
	if err := syncable.Decode(&block); err != nil {
		return err
	}

	if err := s.db.Blocks.DeleteByHeight(s.height); err != nil {
		return err
	}

	record, err := mapper.Block(&block)
	if err != nil {
		return err
	}

	return s.db.Blocks.Create(record)
}

func (s *Sync) processValidatorsSyncable(syncable *model.Syncable) error {
	validators := []near.Validator{}
	if err := syncable.Decode(&validators); err != nil {
		return err
	}

	if err := s.db.Validators.DeleteByHeight(s.height); err != nil {
		return err
	}

	for _, v := range validators {
		record, err := mapper.Validator(s.block, &v)
		if err != nil {
			return err
		}
		if err := s.db.Validators.Create(&record); err != nil {
			return err
		}
	}

	return nil
}

func shouldIgnoreError(err error) bool {
	return strings.Contains(err.Error(), "Server error")
}
