package server

import (
	"context"
	"fmt"
)

func (s *Server) ProcessPushCommand(ctx context.Context) error {
	err := s.ProcessPullCommand(ctx)
	if err != nil {
		return fmt.Errorf("server push command: %w", err)
	}

	storageRecords, err := s.inmemory.GetAllRecordsForServer()
	if err != nil {
		return fmt.Errorf("extract records for push on server: %w", err)
	}
	if err = s.client.Push(ctx, storageRecords); err != nil {
		return fmt.Errorf("push records on server: %w", err)
	}

	return nil
}
