package client

import (
	"context"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
)

type BaseClient interface {
	Push(context.Context, []localModels.StorageRecord) error
	Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error)
	Close() error
}
