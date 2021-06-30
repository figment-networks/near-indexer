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
		return errors.New("parameter is required for delegator search")
	}
	return nil
}
