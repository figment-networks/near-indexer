package server

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/mapper"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

// Server handles all HTTP calls
type Server struct {
	router *gin.Engine
	db     *store.Store
	rpc    *near.Client
}

// New returns a new server
func New(cfg *config.Config, db *store.Store, rpc *near.Client) Server {
	router := gin.Default()

	s := Server{
		router: router,
		db:     db,
		rpc:    rpc,
	}

	if cfg.RollbarToken != "" {
		router.Use(RollbarMiddleware())
	}

	router.GET("/", s.GetEndpoints)
	router.GET("/health", s.GetHealth)
	router.GET("/status", s.GetStatus)
	router.GET("/height", s.GetHeight)
	router.GET("/epochs", s.GetEpochs)
	router.GET("/epochs/:id", s.GetEpoch)
	router.GET("/block", s.GetRecentBlock)
	router.GET("/blocks", s.GetBlocks)
	router.GET("/blocks/:id", s.GetBlock)
	router.GET("/block_times", s.GetBlockTimes)
	router.GET("/block_times_interval", s.GetBlockTimesInterval)
	router.GET("/validators", s.GetValidators)
	router.GET("/validators/:id", s.GetValidator)
	router.GET("/validators/:id/epochs", s.GetValidatorEpochs)
	router.GET("/validators/:id/events", s.GetValidatorEvents)
	router.GET("/validator_times_interval", s.GetValidatorTimesInterval)
	router.GET("/transactions", s.GetTransactions)
	router.GET("/transactions/:id", s.GetTransaction)
	router.GET("/accounts/:id", s.GetAccount)
	router.GET("/delegations/:id", s.GetDelegations)
	router.GET("/events", s.GetEvents)

	return s
}

// Run runs the server
func (s Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetEndpoints returns a list of all available endpoints
func (s Server) GetEndpoints(c *gin.Context) {
	jsonOk(c, gin.H{
		"endpoints": gin.H{
			"/health":                "Get service health",
			"/status":                "Get service and network status",
			"/height":                "Get current block height",
			"/block":                 "Get current block details",
			"/blocks":                "Get latest blocks",
			"/blocks/:id":            "Get block details by height or hash",
			"/block_times":           "Get average block times",
			"/block_times_interval":  "Get average block times for a given interval",
			"/epochs":                "Get list of epochs",
			"/epochs/:id":            "Get epoch details",
			"/validators":            "List all validators",
			"/validators/:id":        "Get validator details",
			"/validators/:id/epochs": "Get validator epochs performance",
			"/validators/:id/events": "Get validator events",
			"/transactions":          "List all recent transactions",
			"/transactions/:id":      "Get transaction details",
			"/accounts/:id":          "Get accoun details",
			"/delegations/:id":       "Get account delegations",
			"/events":                "Get list of events",
		},
	})
}

// GetHealth renders the server health status
func (s Server) GetHealth(c *gin.Context) {
	dbErr := s.db.Test()
	_, nodeErr := s.rpc.Status()

	if dbErr != nil || nodeErr != nil {
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
		"sync_status": "stale",
	}

	if block, err := s.db.Blocks.Last(); err == nil {
		data["last_block_time"] = block.Time
		data["last_block_height"] = block.Height

		if time.Since(block.Time).Seconds() <= 300 {
			data["sync_status"] = "current"
		}
	}

	if status, err := s.rpc.Status(); err == nil {
		data["network_name"] = status.ChainID
		data["network_version"] = status.Version.Version
		data["node_block_time"] = status.SyncInfo.LatestBlockTime.Format(time.RFC3339)
		data["node_block_height"] = status.SyncInfo.LatestBlockHeight
	}

	jsonOk(c, data)
}

// GetHeight renders the last indexed height
func (s Server) GetHeight(c *gin.Context) {
	block, err := s.db.Blocks.Last()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, gin.H{
		"height": block.Height,
		"time":   block.Time,
	})
}

// GetEpochs returns a list of recent epochs
func (s Server) GetEpochs(c *gin.Context) {
	epochs, err := s.db.Epochs.Recent(100)
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, epochs)
}

// GetEpoch returns a single epoch details
func (s Server) GetEpoch(c *gin.Context) {
	epoch, err := s.db.Epochs.FindByID(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, epoch)
}

// GetRecentBlock renders the last indexed block
func (s Server) GetRecentBlock(c *gin.Context) {
	block, err := s.db.Blocks.Last()
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

// GetBlockTimes returns an average block time for the last N blocks
func (s Server) GetBlockTimes(c *gin.Context) {
	params := blockTimesParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.BlockTimes(params.Limit)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

// GetBlockTimesInterval returns average block times for a given time period
func (s Server) GetBlockTimesInterval(c *gin.Context) {
	params := timesIntervalParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.BlockStats(params.Interval, params.Period)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// GetValidators returns recent validators
func (s Server) GetValidators(c *gin.Context) {
	validators, err := s.db.ValidatorAggs.Top()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, validators)
}

// GetValidatorEpochs returns validator epoch participation
func (s Server) GetValidatorEpochs(c *gin.Context) {
	validator, err := s.db.ValidatorAggs.FindBy("account_id", c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	pagination := store.Pagination{}
	if err := c.Bind(&pagination); err != nil {
		badRequest(c, err)
		return
	}

	result, err := s.db.ValidatorAggs.PaginateValidatorEpochs(validator.AccountID, pagination)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// GetValidatorEvents returns validator events
func (s Server) GetValidatorEvents(c *gin.Context) {
	validator, err := s.db.ValidatorAggs.FindBy("account_id", c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	search := store.EventsSearch{
		ItemID:   validator.AccountID,
		ItemType: "validator",
	}

	if err := c.Bind(&search); err != nil {
		badRequest(c, err)
		return
	}

	events, err := s.db.Events.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, events)
}

// GetValidator returns validator details
func (s Server) GetValidator(c *gin.Context) {
	info, err := s.db.ValidatorAggs.FindBy("account_id", c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	account, err := s.db.Accounts.FindByName(info.AccountID)
	if err != nil {
		if err != store.ErrNotFound {
			serverError(c, err)
			return
		}
		account = nil
	}

	epochs, err := s.db.ValidatorAggs.FindValidatorEpochs(info.AccountID, 30)
	if shouldReturn(c, err) {
		return
	}

	blocks, err := s.db.Blocks.Search()
	if shouldReturn(c, err) {
		return
	}

	events, err := s.db.Events.Search(store.EventsSearch{
		ItemID:   info.AccountID,
		ItemType: "validator",
	})
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, gin.H{
		"validator": info,
		"account":   account,
		"blocks":    blocks,
		"epochs":    epochs,
		"events":    events.Records,
	})
}

// GetValidatorsByHeight renders the validators list for a height
func (s Server) GetValidatorsByHeight(c *gin.Context) {
	height := types.HeightFromString(c.Query("height"))
	if height == 0 {
		block, err := s.db.Blocks.Last()
		if shouldReturn(c, err) {
			return
		}
		height = block.Height
	}

	validators, err := s.db.Validators.ByHeight(height)
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

// GetTransactions returns a list of transactions that match query
func (s Server) GetTransactions(c *gin.Context) {
	var (
		txs []model.Transaction
		err error
	)

	if blockHash := c.Query("block_hash"); blockHash != "" {
		txs, err = s.db.Transactions.FindByBlock(blockHash)
	} else {
		txs, err = s.db.Transactions.Recent(100)
	}

	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, txs)
}

// GetTransaction returns a transaction details
func (s Server) GetTransaction(c *gin.Context) {
	tx, err := s.db.Transactions.FindByHash(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, tx)
}

// GetAccount returns an account by name
func (s Server) GetAccount(c *gin.Context) {
	acc, err := s.db.Accounts.FindByName(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, acc)
}

// GetDelegations returns list of delegations for a given account
func (s Server) GetDelegations(c *gin.Context) {
	rawDelegations, err := s.rpc.Delegations(c.Param("id"), 0, 10000)
	if shouldReturn(c, err) {
		return
	}

	delegations, err := mapper.Delegations(rawDelegations)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, delegations)
}

// GetEvents returns a list of events
func (s Server) GetEvents(c *gin.Context) {
	search := store.EventsSearch{}
	if err := c.Bind(&search); err != nil {
		badRequest(c, err)
		return
	}

	events, err := s.db.Events.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, events)
}
