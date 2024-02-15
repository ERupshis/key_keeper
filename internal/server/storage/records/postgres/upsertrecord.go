package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/retrier"
)

func (p *Postgres) UpsertRecord(ctx context.Context, userID int64, record *models.StorageRecord) error {
	exec := p.createUpdateRecordExecFunc(ctx, userID, record)

	result, err := retrier.RetryCallWithTimeout(ctx, []int{1, 1, 3}, db.DatabaseErrorsToRetry, exec)
	if err != nil {
		return fmt.Errorf("upser record with id '%d': %w", record.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (p *Postgres) createUpdateRecordExecFunc(ctx context.Context, userID int64, record *models.StorageRecord) func(context context.Context) (sql.Result, error) {
	return func(context context.Context) (sql.Result, error) {
		return p.DB.ExecContext(ctx,
			`INSERT INTO records (id, data, deleted, updated_at, user_id)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (id) DO UPDATE SET
					  data = excluded.data,
					  deleted = excluded.deleted,
					  updated_at = excluded.updated_at,
					  user_id = excluded.user_id;`,
			record.ID,
			record.Data,
			record.Deleted,
			record.UpdatedAt,
			userID,
		)
	}
}
