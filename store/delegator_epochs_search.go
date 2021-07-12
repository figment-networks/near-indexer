package store

import (
	"errors"
)

type DelegatorEpochsSearch struct {
	AccountID   string `form:"account_id"`
	ValidatorID string `form:"validator_id"`
	Epoch       string `form:"epoch"`
}

func (s *DelegatorEpochsSearch) Validate() error {
	if s.Epoch == "" && s.AccountID == "" && s.ValidatorID == "" {
		return errors.New("at least a parameter is required for delegator search (epoch, account id or validator id)")
	}
	return nil
}
