package mapper

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// ValidatorAddEvent generates an event for validator addition to set
func ValidatorAddEvent(block *near.Block, validator *near.Validator) (*model.Event, error) {
	event := &model.Event{
		Scope:       model.ScopeStaking,
		Action:      model.ActionValidatorAdded,
		BlockHeight: block.Header.Height,
		BlockTime:   util.ParseTime(block.Header.Timestamp),
		Epoch:       block.Header.EpochID,
		ItemID:      validator.AccountID,
		ItemType:    model.ItemTypeValidator,
		Metadata: types.Map{
			"stake": validator.Stake,
		},
		CreatedAt: time.Now(),
	}

	return event, event.Validate()
}

// ValidatorKickout generates an event from validator kickout record
func ValidatorKickout(block *near.Block, kick *near.ValidatorKickout) (*model.Event, error) {
	metadata := types.NewMap()

	switch kick.Reason.Name {
	case near.Unstaked:
		metadata["reason"] = "unstaked"
	case near.Slashed:
		metadata["reason"] = "slashed"
	case near.DidNotGetASeat:
		metadata["reason"] = "no_seat"
	case near.NotEnoughBlocks:
		metadata = kick.Reason.Data
		metadata["reason"] = "not_enough_blocks"
	case near.NotEnoughChunks:
		metadata = kick.Reason.Data
		metadata["reason"] = "not_enough_chunks"
	case near.NotEnoughStake:
		metadata["reason"] = "not_enough_stake"
		metadata["stake"] = kick.Reason.Data["stake_u128"]
		metadata["threshold"] = kick.Reason.Data["threshold_u128"]
	}

	event := &model.Event{
		Scope:       model.ScopeStaking,
		Action:      model.ActionValidatorRemoved,
		BlockHeight: block.Header.Height,
		BlockTime:   util.ParseTime(block.Header.Timestamp),
		Epoch:       block.Header.EpochID,
		ItemID:      kick.Account,
		ItemType:    model.ItemTypeValidator,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
	}

	return event, event.Validate()
}
