package store

import (
	"github.com/figment-networks/near-indexer/model"
)

type RunsStore struct {
	baseStore
}

// Cleanup removes any runs with a height lower than the provided one
func (s RunsStore) Cleanup(maxHeight uint64) error {
	return s.db.Delete(s.model, "height < ?", maxHeight).Error
}

// Last returns the last run
func (s RunsStore) Last() (*model.Run, error) {
	result := &model.Run{}

	err := s.db.
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}
