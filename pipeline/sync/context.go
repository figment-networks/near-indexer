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

	Status             *near.NodeStatus
	DefaultStartHeight uint64
	BlockHeight        uint64
	Height             *model.Height
	Run                *model.Run
	Lag                uint64

	Block      *near.Block
	Validators []near.Validator

	shouldAbort bool
	lastErr     error
	errors      []error
}

func NewContext(db *store.Store, client *near.Client) *Context {
	return &Context{
		Client: client,
		DB:     db,

		errors: []error{},
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

// FirstError returns the first available error
func (c *Context) FirstError() error {
	if len(c.errors) > 0 {
		return c.errors[0]
	}
	return nil
}

// LastError returns the last available error
func (c *Context) LastError() error {
	return c.lastErr
}

// IsAborted returns true if context is aborted
func (c *Context) IsAborted() bool {
	return c.shouldAbort
}
