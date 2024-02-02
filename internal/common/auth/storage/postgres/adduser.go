package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erupshis/key_keeper/internal/common/auth/models"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/retrier"
)

func (p *Postgres) AddUser(ctx context.Context, user *models.User) error {
	exec := p.createAddUserExecFunc(ctx, user)

	result, err := retrier.RetryCallWithTimeout(ctx, []int{1, 1, 3}, db.DatabaseErrorsToRetry, exec)
	if err != nil {
		return fmt.Errorf("add user  '%d': %w", user.Login, err)
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

func (p *Postgres) createAddUserExecFunc(ctx context.Context, user *models.User) func(context context.Context) (sql.Result, error) {
	return func(context context.Context) (sql.Result, error) {
		return p.DB.ExecContext(ctx,
			`INSERT INTO users (login, password)
					VALUES ($1, $2);`,
			user.Login,
			user.Password,
		)
	}
}
