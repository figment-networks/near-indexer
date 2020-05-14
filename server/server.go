package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/near/client"
)

// Server handles all HTTP calls
type Server struct {
	rpc *client.Client
	*gin.Engine
}

// New returns a new server
func New(rpc *client.Client) Server {
	s := Server{
		rpc:    rpc,
		Engine: gin.Default(),
	}

	s.GET("/health", s.GetHealth)
	s.GET("/status", s.GetStatus)
	s.GET("/height", s.GetHeight)
	s.GET("/block", s.GetCurrentBlock)
	s.GET("/validators", s.GetValidators)

	return s
}

// GetHealth renders the server health status
func (s Server) GetHealth(c *gin.Context) {
	c.String(200, "OK")
}

// GetStatus renders the node status
func (s Server) GetStatus(c *gin.Context) {
	block, err := s.rpc.Status()
	response(c, block, err)
}

// GetHeight returns the current height
func (s Server) GetHeight(c *gin.Context) {
	block, err := s.rpc.CurrentBlock()
	result := gin.H{
		"height": block.Header.Height,
		"time":   block.Header.Timestamp,
	}
	response(c, result, err)
}

// GetCurrentBlock renders the latest available block
func (s Server) GetCurrentBlock(c *gin.Context) {
	block, err := s.rpc.CurrentBlock()
	response(c, block, err)
}

// GetValidators renders all validators
func (s Server) GetValidators(c *gin.Context) {
	validators, err := s.rpc.Validators()
	response(c, validators, err)
}

func response(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, data)
}
