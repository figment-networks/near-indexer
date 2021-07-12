package store

import (
	"errors"
)

type DelegatorEpochsSearch struct {
	AccountID   string `json:"account_id"`
	ValidatorID string `json:"validator_id"`
	Epoch       string `json:"epoch"`
}

func (s *DelegatorEpochsSearch) Validate() error {
	if s.Epoch == "" && s.AccountID == "" && s.ValidatorID == "" {
		return errors.New("at least a parameter is required for delegator search (epoch, account id or validator id)")
	}
	return nil
}
