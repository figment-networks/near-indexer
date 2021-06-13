package pipeline

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

// Payload contains data for a single sync run
type Payload struct {
	Lag         int
	StartHeight uint64
	StartTime   time.Time
	EndHeight   uint64
	EndTime     time.Time
	Tip         *near.Block
	Heights     []*HeightPayload
}

// HeightPayload contains all raw data for a single height
type HeightPayload struct {
	Height uint64
	Error  error
	Skip   bool

	Block                  *near.Block
	Validators             []near.Validator
	Chunks                 []near.ChunkDetails
	Transactions           []near.TransactionDetails
	Delegations            []near.AccountInfo
	Accounts               []near.Account
	RewardFees             map[string]near.RewardFee
	DelegationsByValidator map[string][]near.AccountInfo
	CurrentEpoch           bool
	FirstBlockOfNewEpoch   bool
	PreviousValidators     []near.Validator
	PreviousEpochKickOut   []near.ValidatorKickout
	PreviousBlock          *near.Block

	Parsed *ParsedPayload
}

// SkipWithError marks the payload as skipped
func (p *HeightPayload) SkipWithError(err error) {
	p.Error = err
	p.Skip = true
}

// ParsedPayload contains parsed data for a single height
type ParsedPayload struct {
	Block                  *model.Block
	Epoch                  *model.Epoch
	Transactions           []model.Transaction
	Validators             []model.Validator
	ValidatorAggs          []model.ValidatorAgg
	ValidatorEpochs        []model.ValidatorEpoch
	DelegatorEpochs        []model.DelegatorEpoch
	Accounts               []model.Account
	Events                 []model.Event
}
