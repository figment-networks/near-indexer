package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/store/queries"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Blocks        BlocksStore
	Epochs        EpochsStore
	Accounts      AccountsStore
	Delegators    DelegatorsStore
	Validators    ValidatorsStore
	ValidatorAggs ValidatorAggsStore
	Transactions  TransactionsStore
	Stats         StatsStore
	Events        EventsStore
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
	rows, err := s.db.DB().Query(queries.UtilListTables)
	if err != nil {
		return err
	}
	tables := []string{}

	for {
		var tableName string

		if !rows.Next() {
			break
		}

		err := rows.Scan(&tableName)
		if err != nil {
			return err
		}

		tables = append(tables, tableName)
	}

	for _, table := range tables {
		log.Println("truncating table", table)

		q := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY;", table)
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

		Blocks:        BlocksStore{scoped(conn, model.Block{})},
		Epochs:        EpochsStore{scoped(conn, model.Epoch{})},
		Accounts:      AccountsStore{scoped(conn, model.Account{})},
		Delegators:    DelegatorsStore{scoped(conn, model.DelegatorEpoch{})},
		Validators:    ValidatorsStore{scoped(conn, model.Validator{})},
		ValidatorAggs: ValidatorAggsStore{scoped(conn, model.ValidatorAgg{})},
		Transactions:  TransactionsStore{scoped(conn, model.Transaction{})},
		Events:        EventsStore{scoped(conn, model.Event{})},
		Stats:         StatsStore{baseStore{db: conn}},
	}, nil
}
