package commands

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/binary"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/credential"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/local"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/server"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/text"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Config struct {
	StateMachines   *statemachines.StateMachines
	LocalStorageCmd *local.Local

	BankCard   *bankcard.BankCard
	Credential *credential.Credential
	Text       *text.Text
	Binary     *binary.Binary

	Server *server.Server
}

type Commands struct {
	iactr *interactor.Interactor
	local *local.Local

	sm *statemachines.StateMachines

	bc     *bankcard.BankCard
	creds  *credential.Credential
	text   *text.Text
	binary *binary.Binary

	server *server.Server
}

func NewCommands(iactr *interactor.Interactor, cfg *Config) *Commands {
	return &Commands{
		iactr:  iactr,
		local:  cfg.LocalStorageCmd,
		sm:     cfg.StateMachines,
		bc:     cfg.BankCard,
		creds:  cfg.Credential,
		text:   cfg.Text,
		binary: cfg.Binary,
		server: cfg.Server,
	}
}
