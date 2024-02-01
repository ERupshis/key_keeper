package server

import (
	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/binaries"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
)

type Config struct {
	Inmemory *inmemory.Storage
	Local    *local.FileManager
	Binary   *binaries.BinaryManager

	Client client.BaseClient
	Iactr  *interactor.Interactor
}

type Server struct {
	inmemory *inmemory.Storage
	local    *local.FileManager
	binary   *binaries.BinaryManager

	client client.BaseClient

	iactr *interactor.Interactor
}

func NewServer(cfg *Config) *Server {
	return &Server{
		iactr:    cfg.Iactr,
		local:    cfg.Local,
		client:   cfg.Client,
		inmemory: cfg.Inmemory,
		binary:   cfg.Binary,
	}
}
