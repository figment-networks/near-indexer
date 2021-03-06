package near

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	EmptyTxRoot = "11111111111111111111111111111111"
)

type Version struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

func (v Version) String() string {
	return fmt.Sprintf("%s-%s", v.Version, v.Build)
}

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

type SyncInfo struct {
	LatestBlockHash   string    `json:"latest_block_hash"`
	LatestBlockHeight uint64    `json:"latest_block_height"`
	LatestStateRoot   string    `json:"latest_state_root"`
	LatestBlockTime   time.Time `json:"latest_block_time"`
	Syncing           bool      `json:"syncing"`
}

type NodeStatus struct {
	Version  Version  `json:"version"`
	ChainID  string   `json:"chain_id"`
	RPCAddr  string   `json:"rpc_addr"`
	SyncInfo SyncInfo `json:"sync_info"`
}

type Block struct {
	Author string       `json:"author"`
	Header BlockHeader  `json:"header"`
	Chunks []BlockChunk `json:"chunks"`
}

type BlockHeader struct {
	Height                uint64              `json:"height"`
	EpochID               string              `json:"epoch_id"`
	NextEpochID           string              `json:"next_epoch_id"`
	Hash                  string              `json:"hash"`
	PrevHash              string              `json:"prev_hash"`
	PrevStateRoot         string              `json:"prev_state_root"`
	ChunkReceiptsRoot     string              `json:"chunk_receipts_root"`
	ChunkHeadersRoot      string              `json:"chunk_headers_root"`
	ChunkTxRoot           string              `json:"chunk_tx_root"`
	OutcomeRoot           string              `json:"outcome_root"`
	ChunksIncluded        int                 `json:"chunks_included"`
	ChallengesRoot        string              `json:"challenges_root"`
	Timestamp             int64               `json:"timestamp"`
	TimestampNanosec      string              `json:"timestamp_nanosec"`
	RandomValue           string              `json:"random_value"`
	ValidatorProposals    []ValidatorProposal `json:"validator_proposals"`
	ChunkMask             []bool              `json:"chunk_mask"`
	GasPrice              string              `json:"gas_price"`
	RentPaid              string              `json:"rent_paid"`
	ValidatorReward       string              `json:"validator_reward"`
	TotalSupply           string              `json:"total_supply"`
	ChallengesResult      []interface{}       `json:"challenges_result"`
	LastFinalBlock        string              `json:"last_final_block"`
	LastDsFinalBlock      string              `json:"last_ds_final_block"`
	NextBpHash            string              `json:"next_bp_hash"`
	BlockMerkleRoot       string              `json:"block_merkle_root"`
	Approvals             []interface{}       `json:"approvals"`
	Signature             string              `json:"signature"`
	LatestProtocolVersion int                 `json:"latest_protocol_version"`
}

type BlockChunk struct {
	ChunkHash            string              `json:"chunk_hash"`
	PrevBlockHash        string              `json:"prev_block_hash"`
	OutcomeRoot          string              `json:"outcome_root"`
	PrevStateRoot        string              `json:"prev_state_root"`
	EncodedMerkleRoot    string              `json:"encoded_merkle_root"`
	EncodedLength        int                 `json:"encoded_length"`
	HeightCreated        uint64              `json:"height_created"`
	HeightIncluded       uint64              `json:"height_included"`
	ShardID              int                 `json:"shard_id"`
	GasUsed              int                 `json:"gas_used"`
	GasLimit             int64               `json:"gas_limit"`
	RentPaid             string              `json:"rent_paid"`
	ValidatorReward      string              `json:"validator_reward"`
	BalanceBurnt         string              `json:"balance_burnt"`
	OutgoingReceiptsRoot string              `json:"outgoing_receipts_root"`
	TxRoot               string              `json:"tx_root"`
	ValidatorProposals   []ValidatorProposal `json:"validator_proposals"`
	Signature            string              `json:"signature"`
}

type Account struct {
	Amount        string `json:"amount"`
	Locked        string `json:"locked"`
	CodeHash      string `json:"code_hash"`
	StorageUsage  int    `json:"storage_usage"`
	StoragePaidAt int    `json:"storage_paid_at"`
	BlockHeight   uint64 `json:"block_height"`
	BlockHash     string `json:"block_hash"`
}

type Transaction struct {
	Hash       string        `json:"hash"`
	Nonce      int           `json:"nonce"`
	PublicKey  string        `json:"public_key"`
	ReceiverID string        `json:"receiver_id"`
	Signature  string        `json:"signature"`
	SignerID   string        `json:"signer_id"`
	Actions    []interface{} `json:"actions"`
}

type Validator struct {
	AccountID         string `json:"account_id"`
	IsSlashed         bool   `json:"is_slashed"`
	NumExpectedBlocks int    `json:"num_expected_blocks"`
	NumProducedBlocks int    `json:"num_produced_blocks"`
	PublicKey         string `json:"public_key"`
	Shards            []int  `json:"shards"`
	Stake             string `json:"stake"`
}

type ReceiptsOutcome struct {
	BlockHash string  `json:"block_hash"`
	ID        string  `json:"id"`
	Outcome   Outcome `json:"outcome"`
}

type Status struct {
	SuccessValue     *string     `json:"SuccessValue"`
	SuccessReceiptID *string     `json:"SuccessReceiptId"`
	Failure          interface{} `json:"Failure"`
}

type ActionError struct {
	Index int         `json:"index"`
	Kind  interface{} `json:"kind"`
}

type Outcome struct {
	GasBurnt   int64         `json:"gas_burnt"`
	Logs       []interface{} `json:"logs"`
	ReceiptIds []string      `json:"receipt_ids"`
	Status     Status        `json:"status"`
}

type TransactionOutcome struct {
	BlockHash string  `json:"block_hash"`
	ID        string  `json:"id"`
	Outcome   Outcome `json:"outcome"`
}

type TransactionDetails struct {
	ReceiptsOutcome    []ReceiptsOutcome  `json:"receipts_outcome"`
	Status             interface{}        `json:"status"`
	Transaction        Transaction        `json:"transaction"`
	TransactionOutcome TransactionOutcome `json:"transaction_outcome"`
}

type GasPriceDetails struct {
	GasPrice string `json:"gas_price"`
}

type ChunkDetails struct {
	Header       BlockChunk    `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

type Fisher struct {
	AccountID string `json:"account_id"`
	PublicKey string `json:"public_key"`
	Stake     string `json:"stake"`
}

type AccessKey struct {
	Nonce      int         `json:"nonce"`
	Permission interface{} `json:"permission"`
}

type NetworkInfo struct {
	NumActivePeers int `json:"num_active_peers"`
	MaxPeersCount  int `json:"peer_max_count"`

	KnownProducers []struct {
		ID        string  `json:"id"`
		Address   string  `json:"addr"`
		AccountID *string `json:"account_id"`
	} `json:"known_producers"`

	ActivePeers []struct {
		ID        string  `json:"id"`
		Address   string  `json:"addr"`
		AccountID *string `json:"account_id"`
	} `json:"active_peers"`
}

type BlockChange struct {
	Type    string `json:"type"`
	Account string `json:"account_id"`
}

type BlockChangesResponse struct {
	BlockHash string        `json:"block_hash"`
	Changes   []BlockChange `json:"changes"`
}

type ValidatorProposal struct {
	AccountID string `json:"account_id"`
	PublicKey string `json:"public_key"`
	Stake     string `json:"stake"`
}

type ValidatorKickout struct {
	Account string        `json:"account_id"`
	Reason  KickoutReason `json:"reason"`
}

type ValidatorsResponse struct {
	EpochStartHeight     uint64              `json:"epoch_start_height"`
	CurrentValidators    []Validator         `json:"current_validators"`
	CurrentProposales    []ValidatorProposal `json:"current_proposals"`
	NextValidators       []Validator         `json:"next_validators"`
	PreviousEpochKickout []ValidatorKickout  `json:"prev_epoch_kickout"`
}

type AccountInfo struct {
	Account         string `json:"account_id"`
	UnstakedBalance string `json:"unstaked_balance"`
	StakedBalance   string `json:"staked_balance"`
	CanWithdraw     bool   `json:"can_withdraw"`
}

type RewardFee struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

type QueryResponse struct {
	BlockHash   string `json:"block_hash"`
	BlockHeight uint64 `json:"block_height"`
	Result      []byte `json:"result"`
}
