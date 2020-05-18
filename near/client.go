package near

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/rpc/v2/json2"
)

const (
	methodStatus      = "status"
	methodBlock       = "block"
	methodChunk       = "chunk"
	methodValidators  = "validators"
	methodQuery       = "query"
	methodTransaction = "tx"
	methodGasPrice    = "gas_price"
)

// Client interacts with the node RPC API
type Client struct {
	endpoint string
	client   *http.Client
}

// NewClient returns a new node client
func NewClient(endpoint string) Client {
	return Client{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
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

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json2.DecodeClientResponse(resp.Body, out)
}

// Status returns current status of the node
func (c Client) Status() (status NodeStatus, err error) {
	err = c.Call(methodStatus, nil, &status)
	return
}

// CurrentBlock returns the latest available block
func (c Client) CurrentBlock() (block Block, err error) {
	params := map[string]interface{}{"finality": "final"}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHeight returns a block for a given height
func (c Client) BlockByHeight(id uint64) (block Block, err error) {
	params := map[string]interface{}{"block_id": id}
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
	// NOTE: There's a bug in docs/rpc where it says the second params
	// is optional, however it really requires it and will return an
	// error in case if its not provided.
	args := []interface{}{id, "reference"}
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

// Validators returns current validators
func (c Client) Validators() ([]Validator, error) {
	result := struct {
		CurrentValidators []Validator `json:"current_validators"`
		NextValidators    []Validator `json:"next_validators"`
	}{}
	params := []interface{}{nil}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, err
	}
	return result.CurrentValidators, nil
}

// ValidatorsByHeight returns validators for a given height
func (c Client) ValidatorsByHeight(height uint64) ([]Validator, error) {
	result := struct {
		CurrentValidators []Validator `json:"current_validators"`
		NextValidators    []Validator `json:"next_validators"`
	}{}
	params := []interface{}{height}

	if err := c.Call(methodValidators, params, &result); err != nil {
		return nil, err
	}
	return result.CurrentValidators, nil
}
