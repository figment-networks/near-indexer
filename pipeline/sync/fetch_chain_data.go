package sync

import (
	"log"
)

func FetchChainData(c *Context) {
	fetchBlockData(c)
	if c.IsAborted() {
		return
	}
	if c.Block == nil {
		return
	}

	fetchValidatorsData(c)
	if c.IsAborted() {
		return
	}

	// Chunks might be included into blocks multiple times but only applied once
	if c.Block.Header.ChunksIncluded > 0 {
		fetchChunksData(c)
		fetchTransactionsData(c)
	}
}

func fetchBlockData(c *Context) {
	block, err := c.Client.BlockByHeight(c.BlockHeight)
	if err != nil {
		c.Abort(err)
		return
	}
	c.Block = &block
}

func fetchValidatorsData(c *Context) {
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
