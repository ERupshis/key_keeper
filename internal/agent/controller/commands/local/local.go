package local

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Local struct {
	iactr *interactor.Interactor
	sm    *statemachines.StateMachines
}

func NewLocal(iactr *interactor.Interactor) *Local {
	return &Local{
		iactr: iactr,
	}
}
