package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

type FetcherTask struct {
	rpc    *near.Client
	db     *store.Store
	logger *logrus.Logger

	batchSize   int
	startHeight uint64
}

func NewFetcherTask(
	db *store.Store,
	rpc *near.Client,
	config *config.Config,
	logger *logrus.Logger,
) FetcherTask {
	return FetcherTask{
		rpc:         rpc,
		db:          db,
		logger:      logger,
		batchSize:   config.SyncBatchSize,
		startHeight: config.StartHeight,
	}
}

func (t FetcherTask) Run(ctx context.Context, payload *Payload) error {
	defer logTaskDuration(FetcherTaskName, time.Now())

	var (
		lastHeight  uint64
		lastEpoch   string
		startHeight uint64
		endHeight   uint64
	)

	// Fetch the current block from the chain
	currentBlock, err := t.rpc.CurrentBlock()
	if err != nil {
		t.logger.WithError(err).Error("cant fetch current block")
		return err
	}
	payload.CurrentBlock = &currentBlock

	// Fetch last indexed block
	lastBlock, err := t.db.Blocks.Last()
	if err != nil && err != store.ErrNotFound {
		return nil
	}
	if lastBlock != nil {
		lastHeight = uint64(lastBlock.Height)
		lastEpoch = lastBlock.Epoch
	}

	// Start sync from the next block
	if lastHeight > 0 {
		startHeight = lastHeight + 1
	} else {
		startHeight = t.startHeight
	}

	// Get the genesis block height if no start height was provided
	if startHeight == 0 {
		genesis, err := t.rpc.GenesisConfig()
		if err != nil {
			return err
		}
		startHeight = genesis.GenesisHeight
	}

	// Determine the actual sync end height
	endHeight = startHeight + uint64(t.batchSize)
	if endHeight >= currentBlock.Header.Height {
		endHeight = currentBlock.Header.Height
	}

	// Fetch all heights data concurrenly
	t.fetchHeights(startHeight, endHeight, payload)

	for dataIdx, data := range payload.Heights {
		// Skip any heights with missing or non-existent blocks
		if data.Skip {
			continue
		}

		// Terminate early when fetch failed
		if data.Error != nil {
			return err
		}

		if data.Block.Header.EpochID != lastEpoch {
			t.logger.
				WithField("from", lastEpoch).
				WithField("to", data.Block.Header.EpochID).
				WithField("height", data.Block.Header.Height).
				Info("epoch changed")

			var lastBlockOfEpoch *near.Block
			for _, innerData := range payload.Heights {
				if innerData.Block == nil {
					continue
				}
				if innerData.Block.Header.EpochID == lastEpoch {
					lastBlockOfEpoch = innerData.Block
					break
				}
			}

			if lastBlockOfEpoch == nil {
				t.logger.WithField("epoch", lastEpoch).Info("no block with last epoch in scope, fetching from db")
				epochBlock, err := t.db.Blocks.LastInEpoch(lastEpoch)
				if err == nil {
					logrus.Info("no last block in epoch found!!!")
					lastBlockOfEpoch = &near.Block{
						Header: near.BlockHeader{
							Height: uint64(epochBlock.Height),
						},
					}
				} else {
					if err != store.ErrNotFound {
						return err
					}
				}
			}

			if lastBlockOfEpoch != nil {
				logrus.
					WithField("hash", lastBlockOfEpoch.Header.Hash).
					WithField("height", lastBlockOfEpoch.Header.Height).
					Info("last block of epoch")
			}

			var previousValidators *near.ValidatorsResponse
			if lastBlockOfEpoch != nil {
				logrus.WithField("height", lastBlockOfEpoch.Header.Height).Info("fetching previous validators")
				previousValidators, err = t.rpc.ValidatorsByHeight(lastBlockOfEpoch.Header.Height)
				if err != nil {
					return err
				}
			}

			logrus.WithField("height", data.Height).Info("fetching validators")
			validators, err := t.rpc.ValidatorsByHeight(data.Height)
			if err != nil {
				return err
			}

			data.Validators = validators.CurrentValidators
			data.PreviousEpochKickOut = validators.PreviousEpochKickout
			if previousValidators != nil {
				data.PreviousValidators = previousValidators.CurrentValidators
				data.PreviousBlock = lastBlockOfEpoch
			}

			lastEpoch = data.Block.Header.EpochID
		} else {
			isLastInBatch := dataIdx == len(payload.Heights)-1

			// Fetch validators in the current epoch in the last height of the batch
			if currentBlock.Header.EpochID == data.Block.Header.EpochID && isLastInBatch {
				logrus.WithField("height", data.Height).Debug("fetching validators")
				validators, err := t.rpc.ValidatorsByHeight(data.Height)
				if err != nil {
					return err
				}
				data.Validators = validators.CurrentValidators
				data.CurrentEpoch = true
			}
		}
	}

	return nil
}

// fetchHeight retrieves heights data in parallel
func (t FetcherTask) fetchHeights(startHeight, endHeight uint64, payload *Payload) {
	count := int(endHeight - startHeight)
	payloads := make([]*HeightPayload, count)

	wg := sync.WaitGroup{}
	wg.Add(count)

	idx := 0
	for height := startHeight; height < endHeight; height++ {
		go func(idx int, height uint64) {
			defer wg.Done()
			payloads[idx] = t.fetchHeightData(height)
		}(idx, height)
		idx++
	}

	wg.Wait()

	payload.Heights = []*HeightPayload{}
	for _, p := range payloads {
		if p.Skip {
			continue
		}
		if p.Error != nil {
			continue
		}

		payload.Heights = append(payload.Heights, p)
	}

	if len(payload.Heights) > 0 {
		payload.StartHeight = payload.Heights[0].Height
		payload.StartTime = util.ParseTime(payload.Heights[0].Block.Header.Timestamp)
		payload.EndHeight = payload.Heights[len(payload.Heights)-1].Height
		payload.EndTime = util.ParseTime(payload.Heights[len(payload.Heights)-1].Block.Header.Timestamp)
	} else {
		payload.StartHeight = startHeight
		payload.EndHeight = endHeight
	}
}

// fetchHeightData retrieves all the data for a given height
func (t FetcherTask) fetchHeightData(height uint64) (payload *HeightPayload) {
	logrus.WithField("height", height).Debug("fetching height data")
	defer func() {
		if payload.Error != nil && !payload.Skip {
			logrus.
				WithField("height", height).
				WithError(payload.Error).
				Error("fetch failed")
		}
	}()

	payload = &HeightPayload{
		Height: height,
	}

	block, err := t.rpc.BlockByHeight(payload.Height)
	if err != nil {
		payload.Skip = err == near.ErrBlockMissing || err == near.ErrBlockNotFound
		payload.Error = err
		return
	}
	payload.Block = &block

	if block.Header.ChunksIncluded > 0 {
		for _, blockChunk := range block.Chunks {
			// Skip chunks without any transactions
			if blockChunk.TxRoot == near.EmptyTxRoot {
				continue
			}

			// Fetch chunk details
			chunk, err := t.rpc.Chunk(blockChunk.ChunkHash)
			if err != nil {
				payload.Error = err
				return
			}
			payload.Chunks = append(payload.Chunks, chunk)

			// Build a list of transactions to fetch concurrently
			txHashes := []string{}
			for _, chunkTx := range chunk.Transactions {
				txHashes = append(txHashes, chunkTx.Hash)
			}

			// Fetch all transactions concurrently
			transactions, err := t.fetchBlockTransactions(&block, txHashes)
			if err != nil {
				payload.Error = err
				return
			}
			payload.Transactions = transactions
		}
	}

	return payload
}

type txFetchResult struct {
	transaction near.TransactionDetails
	err         error
}

// fetchBlockTransactions retrieves all transactions in parallel
func (t FetcherTask) fetchBlockTransactions(block *near.Block, hashes []string) ([]near.TransactionDetails, error) {
	results := []txFetchResult{}
	resultsChan := make(chan txFetchResult)

	wg := sync.WaitGroup{}
	wg.Add(len(hashes))

	for _, h := range hashes {
		go func(hash string) {
			defer wg.Done()

			tx, err := t.rpc.Transaction(hash)
			resultsChan <- txFetchResult{
				transaction: tx,
				err:         err,
			}
		}(h)
	}

	go func() {
		for {
			select {
			case res, ok := <-resultsChan:
				if !ok {
					return
				}
				results = append(results, res)
			}
		}
	}()

	wg.Wait()
	close(resultsChan)

	txlist := []near.TransactionDetails{}

	for _, res := range results {
		if res.err != nil {
			return nil, res.err
		}
		txlist = append(txlist, res.transaction)
	}

	return txlist, nil
}
