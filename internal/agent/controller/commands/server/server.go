package server

import (
	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
)

type Config struct {
	Inmemory *inmemory.Storage
	Local    *local.FileManager

	Client client.BaseClient
	Iactr  *interactor.Interactor
}

type Server struct {
	inmemory *inmemory.Storage
	local    *local.FileManager

	client client.BaseClient

	iactr *interactor.Interactor
}

func NewServer(cfg *Config) *Server {
	return &Server{
		iactr:    cfg.Iactr,
		local:    cfg.Local,
		client:   cfg.Client,
		inmemory: cfg.Inmemory,
	}
}
