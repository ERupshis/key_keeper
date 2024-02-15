package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erupshis/key_keeper/internal/common/auth/models"
	"github.com/erupshis/key_keeper/internal/common/auth/storage"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/retrier"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

func (p *Postgres) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := p.createGetUserByLoginQueryFunc(ctx, login)

	rows, err := retrier.RetryCallWithTimeout(ctx, []int{1, 1, 3}, db.DatabaseErrorsToRetry, query)
	if err != nil {
		return nil, fmt.Errorf("select user '%s': %w", login, err)
	}

	defer deferutils.ExecWithLogError(rows.Close, p.logger)
	users, err := p.parseGetUserByLoginResult(rows)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if len(users) != 1 {
		return nil, storage.ErrUserNotFound
	}

	return &users[0], nil
}

func (p *Postgres) createGetUserByLoginQueryFunc(ctx context.Context, login string) func(context context.Context) (*sql.Rows, error) {
	return func(context context.Context) (*sql.Rows, error) {
		return p.DB.QueryContext(ctx,
			`SELECT 
    					id,
    					login,
    					password
       				FROM users WHERE login = $1;`,
			login,
		)
	}
}

func (p *Postgres) parseGetUserByLoginResult(rows *sql.Rows) ([]models.User, error) {
	var res []models.User
	for rows.Next() {
		var tmp models.User
		err := rows.Scan(
			&tmp.ID,
			&tmp.Login,
			&tmp.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("parse db result: %w", err)
		}

		res = append(res, tmp)
	}

	return res, nil
}
