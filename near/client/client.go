package client

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/rpc/v2/json2"

	"github.com/figment-networks/near-indexer/near"
)

const (
	methodStatus     = "status"
	methodBlock      = "block"
	methodValidators = "validators"
)

// Client interacts with the node RPC API
type Client struct {
	endpoint string
	client   *http.Client
}

// New returns a new node client
func New(endpoint string) Client {
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
func (c Client) Status() (status near.Status, err error) {
	err = c.Call(methodStatus, nil, &status)
	return
}

// CurrentBlock returns the latest available block
func (c Client) CurrentBlock() (block near.Block, err error) {
	params := map[string]interface{}{"finality": "final"}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHeight returns a block for a given height
func (c Client) BlockByHeight(id uint64) (block near.Block, err error) {
	params := map[string]interface{}{"block_id": id}
	err = c.Call(methodBlock, params, &block)
	return
}

// BlockByHash returns a block for a given hash
func (c Client) BlockByHash(hash string) (block near.Block, err error) {
	params := map[string]interface{}{"block_id": hash}
	err = c.Call(methodBlock, params, &block)
	return
}

// Validators returns a list of available validators
func (c Client) Validators() ([]near.Validator, error) {
	result := struct {
		Validators []near.Validator `json:"current_validators"`
	}{}
	if err := c.Call(methodValidators, []interface{}{nil}, &result); err != nil {
		return nil, err
	}
	return result.Validators, nil
}
