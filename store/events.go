package store

import "github.com/figment-networks/near-indexer/model"

// EventsStore manages events records
type EventsStore struct {
	baseStore
}

type EventsSearch struct {
	Pagination

	ItemID   string `form:"item_id"`
	ItemType string `form:"item_type"`
	Height   uint64 `form:"height"`
}

// FindByID returns an event for a given ID
func (s EventsStore) FindByID(id int) (*model.Event, error) {
	event := &model.Event{}

	err := s.db.
		Model(event).
		Take(event, "id = ?", id).
		Error

	return event, checkErr(err)
}

// Search performs an event search and returns matching records
func (s EventsStore) Search(search EventsSearch) (*PaginatedResult, error) {
	if err := search.Validate(); err != nil {
		return nil, err
	}

	scope := s.db.
		Model(&model.Event{}).
		Order("id DESC").
		Limit(search.Limit)

	if search.ItemID != "" && search.ItemType != "" {
		scope = scope.Where("item_id = ? AND item_type = ?", search.ItemID, search.ItemType)
	}
	if search.Height > 0 {
		scope = scope.Where("block_height = ?", search.Height)
	}

	var count uint
	if err := scope.Count(&count).Error; err != nil {
		return nil, err
	}

	events := []model.Event{}

	err := scope.
		Offset((search.Page - 1) * search.Limit).
		Limit(search.Limit).
		Find(&events).
		Error

	if err != nil {
		return nil, err
	}

	result := &PaginatedResult{
		Page:    search.Page,
		Limit:   search.Limit,
		Count:   count,
		Records: events,
	}

	return result.Update(), nil
}
