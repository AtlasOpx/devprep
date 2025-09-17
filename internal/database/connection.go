package database

import (
	"database/sql"
	"fmt"
	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
	Builder squirrel.StatementBuilderType
}

func Connect(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &DB{
		DB:      sqlDB,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(sqlDB),
	}, nil
}

func (db *DB) Select(columns ...string) squirrel.SelectBuilder {
	return db.Builder.Select(columns...)
}

func (db *DB) Insert(table string) squirrel.InsertBuilder {
	return db.Builder.Insert(table)
}

func (db *DB) Update(table string) squirrel.UpdateBuilder {
	return db.Builder.Update(table)
}

func (db *DB) Delete(table string) squirrel.DeleteBuilder {
	return db.Builder.Delete(table)
}
