package client

import (
	"context"

	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
)

type BaseClient interface {
	Login(ctx context.Context, creds *models.Credential) error
	Register(ctx context.Context, creds *models.Credential) error

	Push(ctx context.Context, records []localModels.StorageRecord) error
	Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error)
	PushBinary(ctx context.Context, binaries map[string]string) error
	PullBinary(ctx context.Context) error

	Close() error
}
