package commands

import (
	"context"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
)

func (c *Commands) RestoreLocalStorage(ctx context.Context, inmemoryStorage *inmemory.Storage, localStorage *local.FileManager) error {
	exist, err := localStorage.IsFileExist()
	if err != nil {
		c.iactr.Printf("failed to check local storage existence")
	}

	if err = c.local.ProcessRestore(ctx, exist, inmemoryStorage, localStorage); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
