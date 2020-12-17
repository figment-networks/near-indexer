package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

// EpochsStore manages epochs records
type EpochsStore struct {
	baseStore
}

// FindByID returns an epoch for a given ID
func (s EpochsStore) FindByID(id string) (*model.Epoch, error) {
	epoch := &model.Epoch{}

	err := s.db.
		Model(epoch).
		Take(epoch, "id = ?", id).
		Error

	return epoch, checkErr(err)
}

// Recent returns a set of recent epochs
func (s EpochsStore) Recent(limit int) ([]model.Epoch, error) {
	epochs := []model.Epoch{}

	err := s.db.
		Model(&model.Epoch{}).
		Order("start_height DESC").
		Limit(limit).
		Find(&epochs).
		Error

	return epochs, checkErr(err)
}

// UpdateCounts updates epoch counts for a given set of epoch IDs
func (s EpochsStore) UpdateCounts(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return s.db.Exec(queries.EpochsUpdateCounts, ids).Error
}

// Import created epochs records in batch
func (s EpochsStore) Import(records []model.Epoch) error {
	return s.bulkImport(queries.EpochsImport, len(records), func(i int) bulk.Row {
		r := records[i]

		return bulk.Row{
			r.ID,
			r.StartHeight,
			r.StartTime,
			r.EndHeight,
			r.EndTime,
			r.BlocksCount,
			r.ValidatorsCount,
			r.AverageEfficiency,
		}
	})
}
