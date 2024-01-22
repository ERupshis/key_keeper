package binary

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Binary struct {
	iactr *interactor.Interactor
	sm    *statemachines.StateMachines
}

func NewBinary(iactr *interactor.Interactor, machines *statemachines.StateMachines) *Binary {
	return &Binary{
		iactr: iactr,
		sm:    machines,
	}
}
