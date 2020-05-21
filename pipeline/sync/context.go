package sync

import (
	"errors"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

type HandlerFunc func(*Context)

type Context struct {
	Client *near.Client
	DB     *store.Store

	BlockHeight uint64
	Block       *near.Block

	Height    *model.Height
	Run       *model.Run
	Syncables []model.Syncable

	shouldAbort bool
	lastErr     error
	errors      []error
}

func NewContext(db *store.Store, client *near.Client) *Context {
	return &Context{
		Client: client,
		DB:     db,
	}
}

// Abort aborts the chain
func (c *Context) Abort(val interface{}) {
	var err error

	switch val.(type) {
	case error:
		err = val.(error)
	case string:
		err = errors.New(val.(string))
	}

	c.lastErr = err
	c.errors = append(c.errors, err)
	c.shouldAbort = true
}

func (c *Context) LastError() error {
	return c.lastErr
}

func (c *Context) IsAborted() bool {
	return c.shouldAbort
}
