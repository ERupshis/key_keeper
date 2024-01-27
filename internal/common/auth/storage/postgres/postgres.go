package postgres

import (
	"github.com/erupshis/key_keeper/internal/common/auth/storage"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/logger"
)

var (
	_ storage.BaseAuthStorage = (*Postgres)(nil)
)

type Postgres struct {
	*db.Connection

	logger logger.BaseLogger
}

// NewPostgres creates postgresql implementation. Supports migrations and check connection to database.
func NewPostgres(connection *db.Connection, logger logger.BaseLogger) storage.BaseAuthStorage {
	return &Postgres{
		Connection: connection,
		logger:     logger,
	}
}
