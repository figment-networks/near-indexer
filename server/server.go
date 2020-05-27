package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/config"
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
	router.GET("/status", s.GetStatus)
	router.GET("/leaderboard", s.GetTopValidators)
	router.GET("/height", s.GetHeight)
	router.GET("/block", s.GetRecentBlock)
	router.GET("/blocks", s.GetBlocks)
	router.GET("/blocks/:id", s.GetBlock)
	router.GET("/block_times", s.GetBlockTimes)
	router.GET("/block_times_interval", s.GetBlockTimesInterval)
	router.GET("/validators", s.GetValidators)
	router.GET("/validator_times_interval", s.GetValidatorTimesInterval)
	router.GET("/validators/:id", s.GetValidators)
	router.GET("/transactions/:id", s.GetTransaction)
	router.GET("/accounts/:id", s.GetAccount)

	return s
}

// Run runs the server
func (s Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetHealth renders the server health status
func (s Server) GetHealth(c *gin.Context) {
	if err := s.db.Test(); err != nil {
		jsonError(c, 500, "unhealthy")
		return
	}
	jsonOk(c, gin.H{"healthy": true})
}

// GetStatus returns the status of the service
func (s Server) GetStatus(c *gin.Context) {
	data := gin.H{
		"app_name":    config.AppName,
		"app_version": config.AppVersion,
		"git_commit":  config.GitCommit,
		"go_version":  config.GoVersion,
	}

	if block, err := s.db.Blocks.Recent(); err == nil {
		data["last_block_time"] = block.Time
		data["last_block_height"] = block.Height
	}

	jsonOk(c, data)
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

func (s Server) GetBlockTimes(c *gin.Context) {
	params := blockTimesParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.AvgRecentTimes(params.Limit)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

func (s Server) GetBlockTimesInterval(c *gin.Context) {
	params := timesIntervalParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.AvgTimesForInterval(params.Interval, params.Period)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// GetValidators renders the validators list for a height
func (s Server) GetValidators(c *gin.Context) {
	height := types.HeightFromString(c.Query("height"))
	if height == 0 {
		h, err := s.db.Heights.LastSuccessful()
		if shouldReturn(c, err) {
			return
		}
		height = h.Height
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

// GetValidatorTimesInterval returns active validators count over period of time
func (s Server) GetValidatorTimesInterval(c *gin.Context) {
	params := timesIntervalParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Validators.CountsForInterval(params.Interval, params.Period)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// GetTransaction returns a transaction details
func (s Server) GetTransaction(c *gin.Context) {
}

// GetAccount returns an account by name
func (s Server) GetAccount(c *gin.Context) {
	acc, err := s.db.Accounts.FindByName(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, acc)
}
