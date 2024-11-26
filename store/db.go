package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBStore struct {
	*gorm.DB
}

func NewSQLiteBackedStore() (*DBStore, error) {
	db, err := gorm.Open(sqlite.Open("tables.db"))

	if err != nil {
		return nil, err
	}
	return &DBStore{db}, nil
}

func NewPostgresBackedStore(dsn string) (*DBStore, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &DBStore{db}, nil
}

func (s *DBStore) Migrate() error {
	return s.DB.AutoMigrate(&Item{})
}
