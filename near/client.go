package near

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/sirupsen/logrus"
)

const (
	methodStatus         = "status"
	methodNetworkInfo    = "network_info"
	methodBlock          = "block"
	methodChunk          = "chunk"
	methodValidators     = "validators"
	methodQuery          = "query"
	methodTransaction    = "tx"
	methodGasPrice       = "gas_price"
	methodGenesisConfig  = "EXPERIMENTAL_genesis_config"
	methodGenesisRecords = "EXPERIMENTAL_genesis_records"
	methodChangesInBlock = "EXPERIMENTAL_changes_in_block"

	delegationsLimit = 100
)

var (
	// ErrBlockMissing is returned when block has been GC'ed by the node
	ErrBlockMissing = errors.New("block is missing")

	// ErrBlockNotFound is returned when block does not exist at given height
	ErrBlockNotFound = errors.New("block not found")

	// ErrEpochUnknown is returned when epoch can't be obtained from the node
	ErrEpochUnknown = errors.New("unknown epoch")

	// ErrValidatorsUnavailable is returned when requesting validator information using invalid epoch
	ErrValidatorsUnavailable = errors.New("validator info unavailable")
)

var (
	defaultClient = &http.Client{
		Timeout: time.Second * 15,
		Transport: &http.Transport{
			MaxConnsPerHost:     250,
			MaxIdleConnsPerHost: 250,
		},
	}
)

// Client interacts with the node RPC API
type Client interface {
	SetTimeout(time.Duration)
	SetDebug(bool)

	Call(string, interface{}, interface{}) error
	GenesisConfig() (GenesisConfig, error)
	GenesisRecords(int, int) (GenesisRecords, error)
	Status() (NodeStatus, error)
	NetworkInfo() (NetworkInfo, error)
	CurrentBlock() (Block, error)
	BlockByHeight(uint64) (Block, error)
	BlockByHash(string) (Block, error)
	Chunk(string) (ChunkDetails, error)
	Account(id string) (Account, error)
	AccountInfo(string, string, uint64) (*AccountInfo, error)
	Transaction(string) (TransactionDetails, error)
	GasPrice(string) (string, error)
	CurrentValidators() (*ValidatorsResponse, error)
	ValidatorsByEpoch(string) (*ValidatorsResponse, error)
	BlockChanges(interface{}) (BlockChangesResponse, error)
	RewardFee(string) (*RewardFee, error)
	Delegations(string, uint64) ([]AccountInfo, error)
}

// DefaultClient returns a new default RPc client
func DefaultClient(endpoint string) Client {
	return &client{
		endpoint: endpoint,
		client:   defaultClient,
	}
}

// NewClient returns a new RPc client with overrides
func NewClient(endpoint string, httpClient *http.Client) Client {
	return &client{
		endpoint: endpoint,
		client:   httpClient,
	}
}

type client struct {
	endpoint string
	client   *http.Client
}

// SetTimeout changes the client timeout
func (c *client) SetTimeout(dur time.Duration) {
	c.client.Timeout = dur
}

// SetDebug changes the debug mode
func (c *client) SetDebug(val bool) {
	logrus.SetLevel(logrus.DebugLevel)
}

// Call executes a RPC transaction
func (c client) Call(method string, args interface{}, out interface{}) error {
	data, err := json2.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}
	reqBody := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, c.endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, duration, err := reqWithTiming(c.client, req)

	logrus.WithFields(logrus.Fields{
		"method":   method,
		"args":     args,
		"duration": duration.String(),
	}).Debug("rpc call")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json2.DecodeClientResponse(resp.Body, out)
	return c.handleRPCError(err, method, args)
}

// GenesisConfig returns the chain genesis configuration
func (c client) GenesisConfig() (result GenesisConfig, err error) {
	err = c.Call(methodGenesisConfig, nil, &result)
	return
}

// GenesisRecords returns the chain genesis records
func (c client) GenesisRecords(limit, offset int) (result GenesisRecords, err error) {
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	err = c.Call(methodGenesisRecords, []interface{}{args}, &result)
	return
}

// Status returns current status of the node
func (c client) Status() (status NodeStatus, err error) {
	err = c.Call(methodStatus, nil, &status)
	return
}

// NetworkInfo returns current status of the network
func (c client) NetworkInfo() (info NetworkInfo, err error) {
	err = c.Call(methodNetworkInfo, nil, &info)
	return
}

// CurrentBlock returns the latest available block
func (c client) CurrentBlock() (block Block, err error) {
	params := map[string]interface{}{"finality": "final"}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHeight returns a block for a given height
func (c client) BlockByHeight(height uint64) (block Block, err error) {
	params := map[string]interface{}{"block_id": height}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHash returns a block for a given hash
func (c client) BlockByHash(hash string) (block Block, err error) {
	params := map[string]interface{}{"block_id": hash}
	err = c.Call(methodBlock, params, &block)

	return
}

// Chunk returns block chunk details by hash
func (c client) Chunk(hash string) (chunk ChunkDetails, err error) {
	params := []string{hash}
	err = c.Call(methodChunk, params, &chunk)
	return
}

// Account returns an account by id
func (c client) Account(id string) (acc Account, err error) {
	params := map[string]string{
		"request_type": "view_account",
		"finality":     "final",
		"account_id":   id,
	}
	err = c.Call(methodQuery, params, &acc)
	return
}

// AccountInfo returns account delegation balance for a given pool address
func (c client) AccountInfo(poolID string, lookupID string, blockID uint64) (*AccountInfo, error) {
	callArgs, err := argsToBase64(map[string]interface{}{
		"account_id": lookupID,
	})
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"request_type": "call_function",
		"method_name":  "get_account",
		"account_id":   poolID,
		"args_base64":  callArgs,
	}
	if blockID == 0 {
		args["finality"] = "final"
	} else {
		args["block_id"] = blockID
	}

	resp := QueryResponse{}
	if err := c.Call(methodQuery, args, &resp); err != nil {
		return nil, err
	}

	acc := &AccountInfo{}
	return acc, json.Unmarshal(resp.Result, acc)
}

// Transaction returns a transaction by hash
func (c client) Transaction(id string) (tran TransactionDetails, err error) {
	// NOTE: There's a bug in docs/rpc where it says the second param is optional,
	// however it really requires it and will return an error when it's missing.
	args := []interface{}{id, "near"}
	err = c.Call(methodTransaction, args, &tran)
	return
}

// GasPrice returns the current gas price
func (c client) GasPrice(block string) (string, error) {
	result := GasPriceDetails{}
	args := []interface{}{nil}

	err := c.Call(methodGasPrice, args, &result)
	return result.GasPrice, err
}

// CurrentValidators returns the current validators
func (c client) CurrentValidators() (*ValidatorsResponse, error) {
	result := &ValidatorsResponse{}
	params := []interface{}{nil}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ValidatorsByEpoch returns validators for a given height
func (c client) ValidatorsByEpoch(epoch string) (*ValidatorsResponse, error) {
	result := &ValidatorsResponse{}
	params := []interface{}{epoch}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, c.handleRPCError(err, methodValidators, params)
	}
	return result, nil
}

// BlockChanges returns a collection of change events in the block
func (c client) BlockChanges(block interface{}) (result BlockChangesResponse, err error) {
	err = c.Call(methodChangesInBlock, map[string]interface{}{"block_id": block}, &result)
	return result, err
}

// RewardFee returns a reward fee for an account
func (c client) RewardFee(account string) (*RewardFee, error) {
	callArgs, err := argsToBase64(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"finality":     "final",
		"request_type": "call_function",
		"method_name":  "get_reward_fee_fraction",
		"account_id":   account,
		"args_base64":  callArgs,
	}

	result := QueryResponse{}
	err = c.Call(methodQuery, args, &result)
	if err != nil {
		return nil, err
	}

	fee := &RewardFee{}
	return fee, json.Unmarshal(result.Result, fee)
}

// Delegations returns a list of delegations for a given account
func (c client) Delegations(account string, blockID uint64) ([]AccountInfo, error) {
	var (
		result   []AccountInfo
		startIdx int
	)

	for {
		callArgs, err := argsToBase64(map[string]interface{}{
			"from_index": startIdx,
			"limit":      delegationsLimit,
		})
		if err != nil {
			return nil, err
		}

		args := map[string]interface{}{
			"request_type": "call_function",
			"method_name":  "get_accounts",
			"account_id":   account,
			"args_base64":  callArgs,
		}
		if blockID == 0 {
			args["finality"] = "final"
		} else {
			args["block_id"] = blockID
		}

		resp := QueryResponse{}
		if err := c.Call(methodQuery, args, &resp); err != nil {
			return nil, err
		}
		if len(resp.Result) == 0 {
			break
		}

		delegations := []AccountInfo{}
		if err := json.Unmarshal(resp.Result, &delegations); err != nil {
			return nil, err
		}
		if len(delegations) == 0 {
			break
		}

		result = append(result, delegations...)
		startIdx += int(delegationsLimit)
	}

	return result, nil
}

func (c client) handleServerError(err *json2.Error) error {
	if err.Code == json2.E_SERVER {
		if msg, ok := err.Data.(string); ok {
			msg = strings.ToLower(msg)

			if strings.Contains(msg, "db not found error") {
				return ErrBlockNotFound
			}
			if strings.Contains(msg, "block missing") {
				return ErrBlockMissing
			}
			if strings.Contains(msg, "unknown epoch") {
				return ErrEpochUnknown
			}
			if strings.Contains(msg, "validator info unavailable") {
				return ErrValidatorsUnavailable
			}
		}
	}
	return errors.New(err.Message)
}

func (c client) handleRPCError(err error, method string, params interface{}) error {
	if err != nil {
		switch err.(type) {
		case *json2.Error:
			e := err.(*json2.Error)

			logrus.WithFields(logrus.Fields{
				"code":    e.Code,
				"message": e.Message,
				"data":    e.Data,
				"method":  method,
				"params":  params,
			}).Debug("rpc service error")

			return c.handleServerError(e)

		default:
			logrus.WithFields(logrus.Fields{
				"method": method,
				"params": params,
			}).WithError(err).Debug("rpc service error")
		}
	}
	return err
}

func reqWithTiming(c *http.Client, req *http.Request) (*http.Response, time.Duration, error) {
	ts := time.Now()
	resp, err := c.Do(req)
	te := time.Since(ts)

	return resp, te, err
}

func argsToBase64(input interface{}) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
