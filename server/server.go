package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/store"
)

// Server handles all HTTP calls
type Server struct {
	router *gin.Engine
	db     *store.Store
}

// New returns a new server
func New(db *store.Store) Server {
	router := gin.Default()

	s := Server{
		router: router,
		db:     db,
	}

	router.GET("/health", s.GetHealth)
	router.GET("/leaderboard", s.GetTopValidators)
	router.GET("/height", s.GetHeight)
	router.GET("/block", s.GetRecentBlock)
	router.GET("/blocks", s.GetBlocks)
	router.GET("/blocks/:id", s.GetBlock)
	router.GET("/validators", s.GetValidators)
	router.GET("/validators/:id", s.GetValidators)
	router.GET("/transactions/:id", s.GetTransaction)

	return s
}

// Run runs the server
func (s Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetHealth renders the server health status
func (s Server) GetHealth(c *gin.Context) {
	if err := s.db.Test(); err != nil {
		c.String(http.StatusInternalServerError, "ERROR")
		return
	}
	c.String(200, "OK")
}

// GetHeight renders the last indexed height
func (s Server) GetHeight(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, gin.H{
		"height": block.Height,
		"time":   block.Time,
	})
}

// GetRecentBlock renders the last indexed block
func (s Server) GetRecentBlock(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, block)
}

// GetBlocks renders blocks that match search params
func (s Server) GetBlocks(c *gin.Context) {
	blocks, err := s.db.Blocks.Search()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, blocks)
}

// GetBlock renders a block for a given height or hash
func (s Server) GetBlock(c *gin.Context) {
	var block *model.Block
	var err error

	rid := resourceID(c, "id")
	if rid.IsNumeric() {
		block, err = s.db.Blocks.FindByHeight(rid.UInt64())
	} else {
		block, err = s.db.Blocks.FindByHash(rid.String())
	}
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, block)
}

// GetValidators renders the validators list for a height
func (s Server) GetValidators(c *gin.Context) {
	height := types.HeightFromString(c.Query("height"))
	if height == 0 {
		b, err := s.db.Blocks.Recent()
		if shouldReturn(c, err) {
			return
		}
		height = b.Height
	}

	validators, err := s.db.Validators.ByHeight(height)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, validators)
}

// GetTopValidators returns top validators
func (s Server) GetTopValidators(c *gin.Context) {
	validators, err := s.db.ValidatorAggs.Top()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, validators)
}

// GetTransaction returns a transaction details
func (s Server) GetTransaction(c *gin.Context) {
}
