package client

import (
	"context"
	"net/http"

	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
)

var (
	_ BaseClient = (*Default)(nil)
)

type Default struct {
	client *http.Client
	host   string
}

func NewDefault(host string) BaseClient {
	return &Default{
		client: http.DefaultClient,
		host:   host,
	}
}

func (d *Default) Login(ctx context.Context, creds *models.Credential) error {
	return nil
}

func (d *Default) Register(ctx context.Context, creds *models.Credential) error {
	return nil
}

func (d *Default) Push(ctx context.Context, storageRecords []localModels.StorageRecord) error {
	return nil
}

func (d *Default) Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error) {
	return nil, nil
}

func (d *Default) Close() error {
	return nil
}
