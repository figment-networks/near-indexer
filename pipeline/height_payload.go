package pipeline

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/near"
)

// HeightPayload contains all raw data for a single height
type HeightPayload struct {
	Height uint64
	Error  error
	Skip   bool

	Block        *near.Block
	Validators   []near.Validator
	Chunks       []near.ChunkDetails
	Transactions []near.TransactionDetails
	Delegations  []near.Delegation
	Accounts     []near.Account

	CurrentEpoch         bool
	PreviousValidators   []near.Validator
	PreviousEpochKickOut []near.ValidatorKickout
	PreviousBlock        *near.Block

	Parsed *ParsedPayload
}

// SkipWithError marks the payload as skipped
func (p *HeightPayload) SkipWithError(err error) {
	p.Error = err
	p.Skip = true
}

type ParsedPayload struct {
	Block           *model.Block
	Epoch           *model.Epoch
	Transactions    []model.Transaction
	Validators      []model.Validator
	ValidatorAggs   []model.ValidatorAgg
	ValidatorEpochs []model.ValidatorEpoch
	Accounts        []model.Account
	Events          []model.Event
}
