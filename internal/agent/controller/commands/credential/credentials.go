package credential

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Credential struct {
	iactr *interactor.Interactor
	sm    *statemachines.StateMachines
}

func NewCredentials(iactr *interactor.Interactor, machines *statemachines.StateMachines) *Credential {
	return &Credential{
		iactr: iactr,
		sm:    machines,
	}
}
