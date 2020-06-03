package near

import (
	"encoding/json"
	"time"
)

type GenesisConfig struct {
	ConfigVersion         int       `json:"config_version"`
	ProtocolVersion       int       `json:"protocol_version"`
	ChainID               string    `json:"chain_id"`
	GenesisHeight         uint64    `json:"genesis_height"`
	GenesisTime           time.Time `json:"genesis_time"`
	NumBlockProducerSeats int       `json:"num_block_producer_seats"`
	EpochLength           int       `json:"epoch_length"`
	TotalSupply           string    `json:"total_supply"`
	Validators            []struct {
		AccountID string `json:"account_id"`
		PublicKey string `json:"public_key"`
		Amount    string `json:"amount"`
	} `json:"validators"`
}

type GenesisRecords struct {
	Records    []json.RawMessage `json:"records"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}
