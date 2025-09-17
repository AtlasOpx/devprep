package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type TestDatabase struct {
	DB       *database.DB
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func SetupTestDatabase() (*TestDatabase, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=testpassword",
			"POSTGRES_USER=testuser",
			"POSTGRES_DB=devprep_test",
			"listen_addresses = '*'",
		},
	}

	resource, err := pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://testuser:testpassword@%s/devprep_test?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	resource.Expire(120)

	pool.MaxWait = 120 * time.Second
	var db *sql.DB
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	cfg := &config.Config{
		DatabaseURL: databaseURL,
	}

	dbWrapper, err := database.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not connect to test database: %s", err)
	}

	err = runMigrations(dbWrapper.DB)
	if err != nil {
		return nil, fmt.Errorf("could not run migrations: %s", err)
	}

	return &TestDatabase{
		DB:       dbWrapper,
		Pool:     pool,
		Resource: resource,
	}, nil
}

func (td *TestDatabase) Cleanup() error {
	if td.DB != nil {
		td.DB.Close()
	}
	if td.Pool != nil && td.Resource != nil {
		return td.Pool.Purge(td.Resource)
	}
	return nil
}

func (td *TestDatabase) CleanupData() error {
	if td.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	queries := []string{
		"DELETE FROM sessions",
		"DELETE FROM users",
	}

	for _, query := range queries {
		_, err := td.DB.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to cleanup data: %s", err)
		}
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			user_agent TEXT,
			ip_address INET,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %s", err)
		}
	}

	return nil
}
