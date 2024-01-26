package commands

import (
	"context"

	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
)

func (c *Commands) Push(ctx context.Context, client client.BaseClient, inmemory *inmemory.Storage) {
	c.Pull(ctx, client, inmemory)

	storageRecords, err := inmemory.GetAllRecordsForServer()
	if err != nil {
		c.iactr.Printf("failed to extract records for push on server: %v", err)
	}
	if err = client.Push(ctx, storageRecords); err != nil {
		c.iactr.Printf("failed to push records on server: %v", err)
	}
}
