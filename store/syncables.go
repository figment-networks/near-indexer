package store

import (
	"time"

	"github.com/figment-networks/near-indexer/model"
)

// SyncablesStore handles operations on syncables
type SyncablesStore struct {
	baseStore
}

// Exists returns true if a syncable of a given kind exists at give height
func (s SyncablesStore) Exists(kind string, height int64) (exists bool, err error) {
	result := &model.Syncable{}

	err = s.db.
		Where("processed_at IS NOT NULL").
		Where("type = ? AND height = ?", kind, height).
		First(result).
		Error

	exists = result.ID > 0

	if err != nil && isNotFound(err) {
		err = nil
	}
	return
}

// Count returns the total number of syncables
func (s SyncablesStore) Count(kind string) (int, error) {
	var result int

	err := s.db.
		Where("type = ?", kind).
		Count(&result).
		Error

	return result, checkErr(err)
}

// MarkProcessed updates the processed timestamp and saves the changes
func (s SyncablesStore) MarkProcessed(syncable *model.Syncable) error {
	now := time.Now()
	syncable.ProcessedAt = &now

	return s.Update(syncable)
}

// FindMostRecent returns the most recent processed syncable for type
func (s SyncablesStore) FindMostRecent(kind string) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Where("processed_at IS NOT NULL").
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}

// Height returns the lowest most recent processed height
func (s SyncablesStore) Height() (uint64, error) {
	result := struct {
		Height uint64
	}{}

	err := s.db.
		Raw(syncablesRecentHeight).
		Scan(&result).
		Error

	return result.Height, checkErr(err)
}
