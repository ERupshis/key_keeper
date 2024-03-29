package postgres

import (
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/server/storage/records"
)

var (
	_ records.BaseStorage = (*Postgres)(nil)
)

type Postgres struct {
	*db.Connection

	logger logger.BaseLogger
}

// NewPostgres creates postgresql implementation. Supports migrations and check connection to database.
func NewPostgres(connection *db.Connection, logger logger.BaseLogger) records.BaseStorage {
	return &Postgres{
		Connection: connection,
		logger:     logger,
	}
}
