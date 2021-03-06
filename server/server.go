package server

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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
	rpc    near.Client
	log    *logrus.Logger
}

// New returns a new server
func New(cfg *config.Config, db *store.Store, logger *logrus.Logger, rpc near.Client) Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	if cfg.RollbarToken != "" {
		router.Use(RollbarMiddleware())
	}

	s := Server{
		router: router,
		db:     db,
		rpc:    rpc,
		log:    logger,
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
	router.GET("/block_stats", s.GetBlockStats)
	router.GET("/validators", s.GetValidators)
	router.GET("/validators/:id", s.GetValidator)
	router.GET("/validators/:id/epochs", s.GetValidatorEpochs)
	router.GET("/validators/:id/events", s.GetValidatorEvents)
	router.GET("/delegators/:id/rewards", s.GetDelegatorRewards)
	router.GET("/transactions", s.GetTransactions)
	router.GET("/transactions/:id", s.GetTransaction)
	router.GET("/accounts/:id", s.GetAccount)
	router.GET("/delegations/:id", s.GetDelegations)
	router.GET("/delegators", s.GetDelegators)
	router.GET("/events", s.GetEvents)
	router.GET("/events/:id", s.GetEvent)

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
			"/health":                 "Get service health",
			"/status":                 "Get service and network status",
			"/height":                 "Get current block height",
			"/block":                  "Get current block details",
			"/blocks":                 "Get latest blocks",
			"/blocks/:id":             "Get block details by height or hash",
			"/block_times":            "Get average block times",
			"/block_stats":            "Get block stats for a time bucket",
			"/epochs":                 "Get list of epochs",
			"/epochs/:id":             "Get epoch details",
			"/validators":             "List all validators",
			"/validators/:id":         "Get validator details",
			"/validators/:id/epochs":  "Get validator epochs performance",
			"/validators/:id/events":  "Get validator events",
			"/delegators/:id/rewards": "Get delegator rewards",
			"/transactions":           "List all recent transactions",
			"/transactions/:id":       "Get transaction details",
			"/accounts/:id":           "Get account details",
			"/delegations/:id":        "Get account delegations",
			"/events":                 "Get list of events",
			"/events/:id":             "Get event details",
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
		data["last_block_height"] = block.ID

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
		"height": block.ID,
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

// GetBlockStats returns block stats for a given time bucket
func (s Server) GetBlockStats(c *gin.Context) {
	params := statsParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	if err := params.Validate(); err != nil {
		badRequest(c, err)
		return
	}

	result, err := s.db.Blocks.BlockStats(params.Bucket, params.Limit)
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

// GetDelegatorRewards returns delegator rewards
func (s Server) GetDelegatorRewards(c *gin.Context) {
	var params delegatorRewardsParams
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}

	if err := params.Validate(); err != nil {
		badRequest(c, err)
		return
	}

	delegatorEpochs, err := s.db.Delegators.SearchDelegatorEpochs(store.DelegatorEpochsSearch{
		AccountID: c.Param("id"),
	})
	if shouldReturn(c, err) {
		return
	}
	if len(delegatorEpochs) == 0 {
		badRequest(c, errors.New("account is not found"))
		return
	}

	interval, _ := model.GetTypeForTimeInterval(params.Interval)

	resp, err := s.db.Delegators.FetchRewardsByInterval(c.Param("id"), params.ValidatorId, params.From, params.To, interval)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, resp)
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
		height = block.ID
	}

	validators, err := s.db.Validators.ByHeight(height)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, validators)
}

// GetTransactions returns a list of transactions that match query
func (s Server) GetTransactions(c *gin.Context) {
	search := store.TransactionsSearch{}
	if err := c.Bind(&search); err != nil {
		badRequest(c, err)
		return
	}

	transactions, err := s.db.Transactions.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, transactions)
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
	account, err := s.rpc.Account(c.Param("id"))
	if shouldReturn(c, err) {
		s.log.WithError(err).Error("unable to fetch account details")
		return
	}
	jsonOk(c, account)
}

// GetDelegations returns list of delegations for a given account
func (s Server) GetDelegations(c *gin.Context) {
	var (
		blockID uint64
		err     error
	)

	if val := c.Query("block_id"); val != "" {
		blockID, err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			badRequest(c, "Invalid block_id value")
			return
		}
	}

	rawDelegations, err := s.rpc.Delegations(c.Param("id"), blockID)
	if shouldReturn(c, err) {
		return
	}

	delegations, err := mapper.Delegations(rawDelegations)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, delegations)
}

// GetDelegators returns list of delegators
func (s Server) GetDelegators(c *gin.Context) {
	params := store.DelegatorEpochsSearch{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}

	delegatorEpochs, err := s.db.Delegators.SearchDelegatorEpochs(params)
	if shouldReturn(c, err) {
		return
	}

	delegators, err := mapper.Delegators(delegatorEpochs)
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, delegators)
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

// GetEvent returns a single event
func (s Server) GetEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	event, err := s.db.Events.FindByID(eventID)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, event)
}
