package near

import (
	"bytes"
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
)

var (
	// ErrBlockMissing is returned when block has been GC'ed by the node
	ErrBlockMissing = errors.New("block is missing")

	// ErrBlockNotFound is returned when block does not exist at given height
	ErrBlockNotFound = errors.New("block not found")
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
type Client struct {
	endpoint string
	client   *http.Client
}

// DefaultClient returns a new default RPC client
func DefaultClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client:   defaultClient,
	}
}

// NewClient returns a new RPC client with overrides
func NewClient(endpoint string, httpClient *http.Client) *Client {
	return &Client{
		endpoint: endpoint,
		client:   httpClient,
	}
}

func (c *Client) SetTimeout(dur time.Duration) {
	c.client.Timeout = dur
}

func (c *Client) handleServerError(err *json2.Error) error {
	if err.Code == json2.E_SERVER {
		if msg, ok := err.Data.(string); ok {
			msg = strings.ToLower(msg)

			if strings.Contains(msg, "db not found error") {
				return ErrBlockNotFound
			}
			if strings.Contains(msg, "block missing") {
				return ErrBlockMissing
			}
		}
	}
	return errors.New(err.Message)
}

func (c *Client) handleRPCError(err error, method string, params interface{}) error {
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

// SetDebug changes the debug mode
func (c *Client) SetDebug(val bool) {
	logrus.SetLevel(logrus.DebugLevel)
}

// Call executes a RPC transaction
func (c Client) Call(method string, args interface{}, out interface{}) error {
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
func (c Client) GenesisConfig() (result GenesisConfig, err error) {
	err = c.Call(methodGenesisConfig, nil, &result)
	return
}

// GenesisRecords returns the chain genesis records
func (c Client) GenesisRecords(limit, offset int) (result GenesisRecords, err error) {
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	err = c.Call(methodGenesisRecords, []interface{}{args}, &result)
	return
}

// Status returns current status of the node
func (c Client) Status() (status NodeStatus, err error) {
	err = c.Call(methodStatus, nil, &status)
	return
}

// NetworkInfo returns current status of the network
func (c Client) NetworkInfo() (info NetworkInfo, err error) {
	err = c.Call(methodNetworkInfo, nil, &info)
	return
}

// CurrentBlock returns the latest available block
func (c Client) CurrentBlock() (block Block, err error) {
	params := map[string]interface{}{"finality": "final"}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHeight returns a block for a given height
func (c Client) BlockByHeight(height uint64) (block Block, err error) {
	params := map[string]interface{}{"block_id": height}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHash returns a block for a given hash
func (c Client) BlockByHash(hash string) (block Block, err error) {
	params := map[string]interface{}{"block_id": hash}
	err = c.Call(methodBlock, params, &block)

	return
}

// Chunk returns block chunk details by hash
func (c Client) Chunk(hash string) (chunk ChunkDetails, err error) {
	params := []string{hash}
	err = c.Call(methodChunk, params, &chunk)
	return
}

// Account returns an account by id
func (c Client) Account(id string) (acc Account, err error) {
	params := map[string]string{
		"request_type": "view_account",
		"finality":     "final",
		"account_id":   id,
	}
	err = c.Call(methodQuery, params, &acc)
	return
}

// Transaction returns a transaction by hash
func (c Client) Transaction(id string) (tran TransactionDetails, err error) {
	// NOTE: There's a bug in docs/rpc where it says the second param is optional,
	// however it really requires it and will return an error when it's missing.
	args := []interface{}{id, "near"}
	err = c.Call(methodTransaction, args, &tran)
	return
}

// GasPrice returns the current gas price
func (c Client) GasPrice(block string) (string, error) {
	result := GasPriceDetails{}
	args := []interface{}{nil}

	err := c.Call(methodGasPrice, args, &result)
	return result.GasPrice, err
}

// CurrentValidators returns the current validators
func (c Client) CurrentValidators() (*ValidatorsResponse, error) {
	result := &ValidatorsResponse{}
	params := []interface{}{nil}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ValidatorsByHeight returns validators for a given height
func (c Client) ValidatorsByHeight(height uint64) (*ValidatorsResponse, error) {
	result := &ValidatorsResponse{}
	params := []interface{}{height}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// BlockChanges returns a collection of change events in the block
func (c Client) BlockChanges(block interface{}) (result BlockChangesResponse, err error) {
	err = c.Call(methodChangesInBlock, map[string]interface{}{"block_id": block}, &result)
	return result, err
}

// RewardFee returns a reward fee for an account
func (c Client) RewardFee(account string) (*RewardFee, error) {
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
func (c Client) Delegations(account string, fromIndex uint64, limit uint64) ([]Delegation, error) {
	callArgs, err := argsToBase64(map[string]interface{}{
		"from_index": fromIndex,
		"limit":      limit,
	})
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"finality":     "final",
		"request_type": "call_function",
		"method_name":  "get_accounts",
		"account_id":   account,
		"args_base64":  callArgs,
	}

	resp := DelegationsResponse{}
	if err := c.Call(methodQuery, args, &resp); err != nil {
		return nil, err
	}

	delegations := []Delegation{}
	if len(resp.Result) == 0 {
		return delegations, nil
	}

	if err := json.Unmarshal(resp.Result, &delegations); err != nil {
		return nil, err
	}
	return delegations, nil
}

func reqWithTiming(c *http.Client, req *http.Request) (*http.Response, time.Duration, error) {
	ts := time.Now()
	resp, err := c.Do(req)
	te := time.Since(ts)

	return resp, te, err
}
