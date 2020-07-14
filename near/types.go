package near

import (
	"fmt"
	"time"
)

const (
	// https://docs.near.org/docs/concepts/transaction
	ActionCreateAccount  = "CreateAccount"  // to make a new account (for a person, company, contract, car, refrigerator, etc)
	ActionDeployContract = "DeployContract" // to deploy a new contract (with its own account)
	ActionFunctionCall   = "FunctionCall"   // to invoke a method on a contract (with budget for compute and storage)
	ActionTransfer       = "Transfer"       // to transfer tokens from one account to another
	ActionStake          = "Stake"          // to express interest in becoming a proof-of-stake validator at the next available opportunity
	ActionAddKey         = "AddKey"         // to add a key to an existing account (either FullAccess or FunctionCall access)
	ActionDeleteKey      = "DeleteKey"      // to delete an existing key from an account
	ActionDeleteAccount  = "DeleteAccount"  // to delete an account (and transfer balance to a beneficiary account)
)

type Version struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

func (v Version) String() string {
	return fmt.Sprintf("%s-%s", v.Version, v.Build)
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
	Height             uint64        `json:"height"`
	EpochID            string        `json:"epoch_id"`
	NextEpochID        string        `json:"next_epoch_id"`
	Hash               string        `json:"hash"`
	PrevHash           string        `json:"prev_hash"`
	PrevStateRoot      string        `json:"prev_state_root"`
	ChunkReceiptsRoot  string        `json:"chunk_receipts_root"`
	ChunkHeadersRoot   string        `json:"chunk_headers_root"`
	ChunkTxRoot        string        `json:"chunk_tx_root"`
	OutcomeRoot        string        `json:"outcome_root"`
	ChunksIncluded     int           `json:"chunks_included"`
	ChallengesRoot     string        `json:"challenges_root"`
	Timestamp          string        `json:"timestamp"`
	RandomValue        string        `json:"random_value"`
	ValidatorProposals []interface{} `json:"validator_proposals"`
	ChunkMask          []bool        `json:"chunk_mask"`
	GasPrice           string        `json:"gas_price"`
	RentPaid           string        `json:"rent_paid"`
	ValidatorReward    string        `json:"validator_reward"`
	TotalSupply        string        `json:"total_supply"`
	ChallengesResult   []interface{} `json:"challenges_result"`
	LastFinalBlock     string        `json:"last_final_block"`
	LastDsFinalBlock   string        `json:"last_ds_final_block"`
	NextBpHash         string        `json:"next_bp_hash"`
	Approvals          []interface{} `json:"approvals"`
	Signature          string        `json:"signature"`
}

type BlockChunk struct {
	ChunkHash            string        `json:"chunk_hash"`
	PrevBlockHash        string        `json:"prev_block_hash"`
	OutcomeRoot          string        `json:"outcome_root"`
	PrevStateRoot        string        `json:"prev_state_root"`
	EncodedMerkleRoot    string        `json:"encoded_merkle_root"`
	EncodedLength        int           `json:"encoded_length"`
	HeightCreated        uint64        `json:"height_created"`
	HeightIncluded       uint64        `json:"height_included"`
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
	BlockHeight   uint64 `json:"block_height"`
	BlockHash     string `json:"block_hash"`
}

// Transaction is a collection of Actions augmented with critical information
type Transaction struct {
	Hash       string        `json:"hash"`
	Nonce      int           `json:"nonce"`
	PublicKey  string        `json:"public_key"`
	ReceiverID string        `json:"receiver_id"`
	Signature  string        `json:"signature"`
	SignerID   string        `json:"signer_id"`
	Actions    []interface{} `json:"actions"`
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

type ReceiptsOutcome struct {
	BlockHash string `json:"block_hash"`
	ID        string `json:"id"`
	Outcome   struct {
		GasBurnt   int64         `json:"gas_burnt"`
		Logs       []interface{} `json:"logs"`
		ReceiptIds []interface{} `json:"receipt_ids"`
		Status     struct {
			SuccessValue string `json:"SuccessValue"`
		} `json:"status"`
	} `json:"outcome"`
}

type Status struct {
	SuccessValue     string `json:"SuccessValue"`
	SuccessReceiptID string `json:"SuccessReceiptId"`
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
	Status             Status             `json:"status"`
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

type Transfer struct {
	Deposit string `json:"deposit"`
}

type Stake struct {
	PublicKey string `json:"public_key"`
	Amount    string `json:"stake"`
}

type Fisher struct {
	AccountID string `json:"account_id"`
	PublicKey string `json:"public_key"`
	Stake     string `json:"stake"`
}
