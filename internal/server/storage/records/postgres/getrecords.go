package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/retrier"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

func (p *Postgres) GetRecords(ctx context.Context, userID int64) ([]models.StorageRecord, error) {
	query := p.createGetRecordsQueryFunc(ctx, userID)

	rows, err := retrier.RetryCallWithTimeout(ctx, []int{1, 1, 3}, db.DatabaseErrorsToRetry, query)
	if err != nil {
		return nil, fmt.Errorf("select records with user_id '%d': %w", userID, err)
	}

	defer deferutils.ExecWithLogError(rows.Close, p.logger)
	return p.parseGetRecordsResult(rows)
}

func (p *Postgres) createGetRecordsQueryFunc(ctx context.Context, userID int64) func(context context.Context) (*sql.Rows, error) {
	return func(context context.Context) (*sql.Rows, error) {
		return p.DB.QueryContext(ctx,
			`SELECT 
    					id,
    					data,
    					deleted,
    					updated_at
       				FROM records WHERE user_id = $1 AND deleted = false;`,
			userID,
		)
	}
}

func (p *Postgres) parseGetRecordsResult(rows *sql.Rows) ([]models.StorageRecord, error) {
	var res []models.StorageRecord
	for rows.Next() {
		var tmp models.StorageRecord
		err := rows.Scan(
			&tmp.ID,
			&tmp.Data,
			&tmp.Deleted,
			&tmp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("parse db result: %w", err)
		}

		res = append(res, tmp)
	}

	return res, nil
}
