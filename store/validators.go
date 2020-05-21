package store

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
)

// ValidatorsStore handles operations on blocks
type ValidatorsStore struct {
	baseStore
}

// ByHeight returns all transactions for a given height
func (s ValidatorsStore) ByHeight(height types.Height) ([]model.Validator, error) {
	result := []model.Validator{}

	err := s.db.
		Where("height = ?", height).
		Order("id DESC").
		Find(&result).
		Error

	return result, err
}
