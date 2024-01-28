package server

import (
	"context"
	"fmt"

	"github.com/erupshis/key_keeper/internal/models"
)

func (s *Server) ProcessRegisterCommand(ctx context.Context) error {
	creds := models.Credential{}
	if err := s.collectCreds(&creds); err != nil {
		return fmt.Errorf("collect credentials: %w", err)
	}

	if err := s.client.Register(ctx, &creds); err != nil {
		return fmt.Errorf("login on server: %w", err)
	}

	return nil
}
