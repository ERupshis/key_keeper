package records

import (
	"context"

	"github.com/erupshis/key_keeper/internal/agent/storage/models"
)

type BaseStorage interface {
	UpsertRecord(ctx context.Context, userID int64, record *models.StorageRecord) error
	GetRecords(ctx context.Context, userID int64) ([]models.StorageRecord, error)
}
