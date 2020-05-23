package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound = errors.New("record not found")
)

// baseStore implements generic store operations
type baseStore struct {
	db    *gorm.DB
	model interface{}
}

// Conn returns an underlying database connection
func (s baseStore) Conn() *sql.DB {
	return s.db.DB()
}

// Create creates a new record. Must pass a pointer.
func (s baseStore) Create(record interface{}) error {
	err := s.db.Create(record).Error
	return checkErr(err)
}

// Update updates the existing record. Must pass a pointer.
func (s baseStore) Update(record interface{}) error {
	err := s.db.Save(record).Error
	return checkErr(err)
}

// Truncate removes all records from the table
func (s baseStore) Truncate() error {
	return s.db.Delete(s.model).Error
}

// DeleteByHeight removes all records associated with a height
func (s baseStore) DeleteByHeight(height uint64) error {
	return s.db.Delete(s.model, "height = ?", height).Error
}

// Import imports records in bulk
func (s baseStore) Import(query string, rows int, fn bulkRowFunc) error {
	return bulkInsert(s.db, query, rows, fn)
}

func scoped(conn *gorm.DB, m interface{}) baseStore {
	return baseStore{conn, m}
}

func isNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err) || err == ErrNotFound
}

func findBy(db *gorm.DB, dst interface{}, key string, value interface{}) error {
	return db.
		Model(dst).
		Where(fmt.Sprintf("%s = ?", key), value).
		First(dst).
		Error
}

func checkErr(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return ErrNotFound
	}
	return err
}
