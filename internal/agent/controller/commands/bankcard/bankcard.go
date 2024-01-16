package bankcard

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type BankCard struct {
	iactr *interactor.Interactor
	sm    *statemachines.StateMachines
}

func NewBankCard(iactr *interactor.Interactor, machines *statemachines.StateMachines) *BankCard {
	return &BankCard{
		iactr: iactr,
		sm:    machines,
	}
}
