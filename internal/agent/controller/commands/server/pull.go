package server

import (
	"context"
	"fmt"
)

func (s *Server) ProcessPullCommand(ctx context.Context) error {
	serverRecords, err := s.client.Pull(ctx)
	if err != nil {
		return fmt.Errorf("pull records from server: %w", err)
	}

	if err = s.inmemory.Sync(serverRecords); err != nil {
		return fmt.Errorf("pull server records: %w", err)
	}

	binaries, err := s.client.PullBinary(ctx)
	if err != nil {
		return fmt.Errorf("pull server binaries: %w", err)
	}

	if err = s.binary.SaveBinaries(binaries); err != nil {
		return fmt.Errorf("save server binaries: %w", err)
	}

	return nil
}
