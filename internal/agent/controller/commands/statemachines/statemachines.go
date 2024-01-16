package statemachines

import (
	"github.com/erupshis/key_keeper/internal/agent/interactor"
)

type StateMachines struct {
	iactr *interactor.Interactor
}

func NewStateMachines(iactr *interactor.Interactor) *StateMachines {
	return &StateMachines{iactr: iactr}
}
