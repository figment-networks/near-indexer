package store

import "github.com/figment-networks/near-indexer/model"

// HeightsStore handles operations on heights
type HeightsStore struct {
	baseStore
}

// Last returns a last height record
func (s HeightsStore) Last() (*model.Height, error) {
	result := &model.Height{}

	err := s.db.
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}
