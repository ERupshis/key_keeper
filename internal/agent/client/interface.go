package client

import (
	"context"

	"github.com/erupshis/key_keeper/internal/common/models"
)

type BaseClient interface {
	Push(ctx context.Context, records []models.Record) error
	Pull(ctx context.Context, records []models.Record) error
}
