package text

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type Text struct {
	iactr *interactor.Interactor
	sm    *statemachines.StateMachines
}

func NewText(iactr *interactor.Interactor, machines *statemachines.StateMachines) *Text {
	return &Text{
		iactr: iactr,
		sm:    machines,
	}
}
