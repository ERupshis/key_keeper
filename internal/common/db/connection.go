package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/common/retrier"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DSN              string
	MigrationsFolder string
}

type Connection struct {
	*sql.DB
}

func NewConnection(ctx context.Context, cfg Config) (*Connection, error) {
	errMsg := "create db: %w"
	database, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsFolder, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf(errMsg, err)
	}

	manager := &Connection{
		DB: database,
	}

	if _, err = manager.CheckConnection(ctx); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return manager, nil
}

func (c *Connection) CheckConnection(ctx context.Context) (bool, error) {
	exec := func(context context.Context) (any, error) {
		return nil, c.PingContext(context)
	}
	_, err := retrier.RetryCallWithTimeout(ctx, nil, DatabaseErrorsToRetry, exec)
	if err != nil {
		return false, fmt.Errorf("check connection: %w", err)
	}
	return true, nil
}

func (c *Connection) Close() error {
	return c.DB.Close()
}
