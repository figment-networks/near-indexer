package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
	"github.com/lib/pq"
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
		Take(epoch, "uuid = ?", id).
		Error

	return epoch, checkErr(err)
}

// Recent returns a set of recent epochs
func (s EpochsStore) Recent(limit int) ([]model.Epoch, error) {
	epochs := []model.Epoch{}

	err := s.db.
		Model(&model.Epoch{}).
		Order("id DESC").
		Limit(limit).
		Find(&epochs).
		Error

	return epochs, checkErr(err)
}

func (s EpochsStore) UpdateCounters(ids []string) error {
	return s.db.Exec(queries.EpochsUpdateCounts, pq.Array(ids)).Error
}

func (s EpochsStore) Import(records []model.Epoch) error {
	return s.bulkImport(queries.EpochsImport, len(records), func(i int) bulk.Row {
		r := records[i]

		return bulk.Row{
			r.UUID,
			r.StartHeight,
			r.StartTime,
			r.EndHeight,
			r.EndTime,
		}
	})
}
