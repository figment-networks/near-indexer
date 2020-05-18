package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/near-indexer/model"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Runs       RunsStore
	Syncables  SyncablesStore
	Blocks     BlocksStore
	Validators ValidatorsStore
}

func (s *Store) Automigrate() error {
	return s.db.AutoMigrate(
		&model.Run{},
		&model.Syncable{},
		&model.Block{},
		&model.Validator{},
	).Error
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// New returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,

		Runs:       NewRunsStore(conn),
		Syncables:  NewSyncablesStore(conn),
		Blocks:     NewBlocksStore(conn),
		Validators: NewValidatorsStore(conn),
	}, nil
}

func NewSyncablesStore(db *gorm.DB) SyncablesStore {
	return SyncablesStore{scoped(db, model.Syncable{})}
}

func NewRunsStore(db *gorm.DB) RunsStore {
	return RunsStore{scoped(db, model.Run{})}
}

func NewBlocksStore(db *gorm.DB) BlocksStore {
	return BlocksStore{scoped(db, model.Block{})}
}

func NewValidatorsStore(db *gorm.DB) ValidatorsStore {
	return ValidatorsStore{scoped(db, model.Validator{})}
}
