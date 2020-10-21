package store

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

// HeightsStore handles operations on heights
type HeightsStore struct {
	baseStore
}

// Last returns a last height record
func (s HeightsStore) Last() (*model.Height, error) {
	result := &model.Height{}

	err := s.db.
		Order("height DESC").
		Limit(1).
		Take(&result).
		Error

	return result, checkErr(err)
}

// LastSuccessful returns the last successful height record
func (s HeightsStore) LastSuccessful() (*model.Height, error) {
	result := &model.Height{}

	err := s.db.
		Where("status = ?", model.HeightStatusOK).
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}

// StatusCounts returns height sync statuses with counts
func (s HeightsStore) StatusCounts() ([]model.HeightStatusCount, error) {
	result := []model.HeightStatusCount{}

	err := s.db.
		Raw(queries.HeightsReport).
		Scan(&result).
		Error

	return result, err
}
