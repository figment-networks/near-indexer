package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

// Server handles all HTTP calls
type Server struct {
	db  *store.Store
	rpc *near.Client
	*gin.Engine
}

// New returns a new server
func New(db *store.Store, rpc *near.Client) Server {
	s := Server{
		db:     db,
		rpc:    rpc,
		Engine: gin.Default(),
	}

	s.GET("/health", s.GetHealth)
	s.GET("/height", s.GetHeight)
	s.GET("/block", s.GetBlock)
	s.GET("/validators", s.GetValidators)

	return s
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
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{
		"height": block.Height,
		"time":   block.Time,
	})
}

// GetBlock renders the last indexed block
func (s Server) GetBlock(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, block)
}

// GetValidators renders the validators list for a height
func (s Server) GetValidators(c *gin.Context) {
	var height uint64

	fmt.Sscanf(c.Query("height"), "%d", &height)
	if height == 0 {
		b, err := s.db.Blocks.Recent()
		if err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		height = b.Height
	}

	validators, err := s.db.Validators.ByHeight(height)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, validators)
}
