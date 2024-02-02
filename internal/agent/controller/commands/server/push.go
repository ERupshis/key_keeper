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

	if err = s.pushRecordsToServer(ctx); err != nil {
		return fmt.Errorf("server push command: %w", err)
	}

	return s.pushBinariesToServer(ctx)
}

func (s *Server) pushRecordsToServer(ctx context.Context) error {
	storageRecords, err := s.inmemory.GetAllRecordsForServer()
	if err != nil {
		return fmt.Errorf("extract records for push on server: %w", err)
	}
	if err = s.client.Push(ctx, storageRecords); err != nil {
		return fmt.Errorf("push records on server: %w", err)
	}

	if err = s.inmemory.RemoveLocalRecords(); err != nil {
		return fmt.Errorf("delete local records error: %w", err)
	}

	return nil
}

func (s *Server) pushBinariesToServer(ctx context.Context) error {
	binFilesList := s.inmemory.GetBinFilesList()
	binFiles, err := s.binary.GetFiles(binFilesList)
	if err != nil {
		return fmt.Errorf("read local bin files: %w", err)
	}

	if err = s.client.PushBinary(ctx, binFiles); err != nil {
		return fmt.Errorf("send bin files to server: %w", err)
	}

	return nil
}
