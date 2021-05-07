package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/near-indexer/model"
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
	rr := "WITH epoch_stats AS (  SELECT    epoch,   MIN(time) AS start_time,    MIN(id) AS start_height,    MAX(time) AS end_time,    MAX(id) AS end_height,    COUNT(1) AS blocks_count,    COUNT(DISTINCT producer) AS validators_count  FROM    blocks  WHERE    epoch IN (?)  GROUP BY    epoch) UPDATE epochs SET  start_time         = epoch_stats.start_time,  start_height       = epoch_stats.start_height,  end_time           = epoch_stats.end_time,  end_height         = epoch_stats.end_height,  blocks_count       = epoch_stats.blocks_count,  validators_count   = epoch_stats.validators_count,  average_efficiency = (    SELECT ROUND(COALESCE(AVG(efficiency), 0), 4) FROM validator_epochs WHERE epoch = epoch_stats.epoch  ) FROM  epoch_stats WHERE  epochs.id = epoch_stats.epoch"
	return s.db.Exec(rr, ids).Error
}

// Import created epochs records in batch
func (s EpochsStore) Import(records []model.Epoch) error {
	rr:= "INSERT INTO epochs (  id,  start_height,  start_time,  end_height,  end_time,  blocks_count,  validators_count,  average_efficiency) VALUES @values  ON CONFLICT (id) DO UPDATE SET  end_height = excluded.end_height,  end_time   = excluded.end_time"
	return s.bulkImport(rr, len(records), func(i int) bulk.Row {
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
