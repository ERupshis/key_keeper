package commands

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/local"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Commands struct {
	iactr *interactor.Interactor

	sm    *statemachines.StateMachines
	bc    *bankcard.BankCard
	local *local.Local
}

func NewCommands(iactr *interactor.Interactor, sm *statemachines.StateMachines, bc *bankcard.BankCard, local *local.Local) *Commands {
	return &Commands{
		iactr: iactr,
		sm:    sm,
		bc:    bc,
		local: local,
	}
}
