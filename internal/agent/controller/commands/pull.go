package commands

import (
	"context"

	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
)

func (c *Commands) Pull(ctx context.Context, client client.BaseClient, inmemory *inmemory.Storage) {
	serverRecords, err := client.Pull(ctx)
	if err != nil {
		c.iactr.Printf("failed to pull records from server: %v", err)
	}

	if err = inmemory.Sync(serverRecords); err != nil {
		c.iactr.Printf("failed to pull server records: %v", err)
	}
}
