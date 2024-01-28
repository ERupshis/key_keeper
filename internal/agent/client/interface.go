package client

import (
	"context"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/models"
)

type BaseClient interface {
	Login(ctx context.Context, creds *models.Credential) error
	Register(ctx context.Context, creds *models.Credential) error

	Push(context.Context, []localModels.StorageRecord) error
	Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error)
	Close() error
}
