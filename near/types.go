package near

import (
	"time"
)

type Version struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

type SyncInfo struct {
	LatestBlockHash   string    `json:"latest_block_hash"`
	LatestBlockHeight uint64    `json:"latest_block_height"`
	LatestStateRoot   string    `json:"latest_state_root"`
	LatestBlockTime   time.Time `json:"latest_block_time"`
	Syncing           bool      `json:"syncing"`
}

type Status struct {
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
	Height              int             `json:"height"`
	EpochID             string          `json:"epoch_id"`
	NextEpochID         string          `json:"next_epoch_id"`
	Hash                string          `json:"hash"`
	PrevHash            string          `json:"prev_hash"`
	PrevStateRoot       string          `json:"prev_state_root"`
	ChunkReceiptsRoot   string          `json:"chunk_receipts_root"`
	ChunkHeadersRoot    string          `json:"chunk_headers_root"`
	ChunkTxRoot         string          `json:"chunk_tx_root"`
	OutcomeRoot         string          `json:"outcome_root"`
	ChunksIncluded      int             `json:"chunks_included"`
	ChallengesRoot      string          `json:"challenges_root"`
	Timestamp           int64           `json:"timestamp"`
	RandomValue         string          `json:"random_value"`
	Score               int             `json:"score"`
	ValidatorProposals  []interface{}   `json:"validator_proposals"`
	ChunkMask           []bool          `json:"chunk_mask"`
	GasPrice            string          `json:"gas_price"`
	RentPaid            string          `json:"rent_paid"`
	ValidatorReward     string          `json:"validator_reward"`
	TotalSupply         string          `json:"total_supply"`
	ChallengesResult    []interface{}   `json:"challenges_result"`
	LastQuorumPreVote   string          `json:"last_quorum_pre_vote"`
	LastQuorumPreCommit string          `json:"last_quorum_pre_commit"`
	LastDsFinalBlock    string          `json:"last_ds_final_block"`
	NextBpHash          string          `json:"next_bp_hash"`
	Approvals           [][]interface{} `json:"approvals"`
	Signature           string          `json:"signature"`
}

type BlockChunk struct {
	ChunkHash            string        `json:"chunk_hash"`
	PrevBlockHash        string        `json:"prev_block_hash"`
	OutcomeRoot          string        `json:"outcome_root"`
	PrevStateRoot        string        `json:"prev_state_root"`
	EncodedMerkleRoot    string        `json:"encoded_merkle_root"`
	EncodedLength        int           `json:"encoded_length"`
	HeightCreated        int           `json:"height_created"`
	HeightIncluded       int           `json:"height_included"`
	ShardID              int           `json:"shard_id"`
	GasUsed              int           `json:"gas_used"`
	GasLimit             int64         `json:"gas_limit"`
	RentPaid             string        `json:"rent_paid"`
	ValidatorReward      string        `json:"validator_reward"`
	BalanceBurnt         string        `json:"balance_burnt"`
	OutgoingReceiptsRoot string        `json:"outgoing_receipts_root"`
	TxRoot               string        `json:"tx_root"`
	ValidatorProposals   []interface{} `json:"validator_proposals"`
	Signature            string        `json:"signature"`
}

type Account struct {
	Amount        string `json:"amount"`
	Locked        string `json:"locked"`
	CodeHash      string `json:"code_hash"`
	StorageUsage  int    `json:"storage_usage"`
	StoragePaidAt int    `json:"storage_paid_at"`
	BlockHeight   int    `json:"block_height"`
	BlockHash     string `json:"block_hash"`
}

type Transaction struct {
	Hash       string   `json:"hash"`
	Nonce      int      `json:"nonce"`
	PublicKey  string   `json:"public_key"`
	ReceiverID string   `json:"receiver_id"`
	Signature  string   `json:"signature"`
	SignerID   string   `json:"signer_id"`
	Actions    []Action `json:"actions"`
}

type Action struct {
	FunctionCall FunctionCall `json:"function_call"`
}

type FunctionCall struct {
	Args       string `json:"args"`
	Deposit    string `json:"deposit"`
	Gas        int64  `json:"gas"`
	MethodName string `json:"method_name"`
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
