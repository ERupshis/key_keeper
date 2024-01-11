package client

import (
	"context"

	"github.com/erupshis/key_keeper/internal/common/data"
)

type BaseClient interface {
	Push(ctx context.Context, records []data.Record) error
	Pull(ctx context.Context, records []data.Record) error
}
