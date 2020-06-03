package sync

import "sync"

func FetchChainData(c *Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go fetchBlockData(wg, c)
	go fetchValidatorsData(wg, c)

	wg.Wait()
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
