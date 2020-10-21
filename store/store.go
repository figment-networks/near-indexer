package store

import (
	"database/sql"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/near-indexer/model"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Runs          RunsStore
	Heights       HeightsStore
	Blocks        BlocksStore
	Accounts      AccountsStore
	Validators    ValidatorsStore
	ValidatorAggs ValidatorAggsStore
	Transactions  TransactionsStore
	Stats         StatsStore
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// Conn returns an underlying database connection
func (s *Store) Conn() *sql.DB {
	return s.db.DB()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// ResetAll performs a full database reset without dropping any objects
func (s *Store) ResetAll() error {
	queries := []string{
		"TRUNCATE TABLE blocks RESTART IDENTITY",
		"TRUNCATE TABLE validators RESTART IDENTITY",
		"TRUNCATE TABLE validator_counts RESTART IDENTITY",
		"TRUNCATE TABLE validator_epochs RESTART IDENTITY",
		"TRUNCATE TABLE validator_aggregates RESTART IDENTITY",
		"TRUNCATE TABLE runs RESTART IDENTITY",
		"TRUNCATE TABLE heights RESTART IDENTITY",
		"TRUNCATE TABLE transactions RESTART IDENTITY",
	}

	for _, q := range queries {
		log.Println("executing", q)
		if err := s.db.Exec(q).Error; err != nil {
			return err
		}
	}

	return nil
}

// New returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,

		Heights:       HeightsStore{scoped(conn, model.Height{})},
		Runs:          RunsStore{scoped(conn, model.Run{})},
		Blocks:        BlocksStore{scoped(conn, model.Block{})},
		Accounts:      AccountsStore{scoped(conn, model.Account{})},
		Validators:    ValidatorsStore{scoped(conn, model.Validator{})},
		ValidatorAggs: ValidatorAggsStore{scoped(conn, model.ValidatorAgg{})},
		Transactions:  TransactionsStore{scoped(conn, model.Transaction{})},
		Stats:         StatsStore{baseStore{db: conn}},
	}, nil
}
