package sync

import (
	"log"
	"sync"
)

func FetchChainData(c *Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go fetchBlockData(wg, c)
	go fetchValidatorsData(wg, c)

	wg.Wait()

	if c.Block == nil {
		return
	}

	fetchChunksData(c)
	fetchTransactionsData(c)
}

func fetchBlockData(wg *sync.WaitGroup, c *Context) {
	defer wg.Done()

	block, err := c.Client.BlockByHeight(c.BlockHeight)
	if err != nil {
		c.Abort(err)
		return
	}
	c.Block = &block
}

func fetchValidatorsData(wg *sync.WaitGroup, c *Context) {
	defer wg.Done()

	validators, err := c.Client.ValidatorsByHeight(c.BlockHeight)
	if err != nil {
		c.Abort(err)
		return
	}
	c.Validators = validators
}

func fetchChunksData(c *Context) {
	for _, blockChunk := range c.Block.Chunks {
		chunk, err := c.Client.Chunk(blockChunk.ChunkHash)
		if err != nil {
			log.Println("cant fetch chunk:", err)
			continue
		}
		c.Chunks = append(c.Chunks, chunk)
	}
}

func fetchTransactionsData(c *Context) {
	for _, chunk := range c.Chunks {
		for _, tx := range chunk.Transactions {
			transaction, err := c.Client.Transaction(tx.Hash)
			if err != nil {
				log.Println("cant fetch transaction:", err)
				continue
			}

			c.Transactions = append(c.Transactions, transaction)
		}
	}
}
