package statemachines

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type stateDelete int

const (
	deleteInitialState = stateDelete(0)
	deleteIDState      = stateDelete(1)
	deleteFinishState  = stateDelete(2)
)

func (s *StateMachines) Delete() (*int64, error) {
	currentState := deleteInitialState

	var id *int64
	for currentState != deleteFinishState {
		switch currentState {
		case deleteInitialState:
			{
				currentState = s.stateDeleteIDInitial()
			}
		case deleteIDState:
			{
				currentStateTmp, idTmp, err := s.stateDeleteIDValue()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return nil, err
					} else {
						continue
					}
				}

				id = &idTmp
				currentState = currentStateTmp
			}
		}
	}

	return id, nil
}

func (s *StateMachines) stateDeleteIDInitial() stateDelete {
	fmt.Printf("insert record %s: ", utils.CommandID)
	return deleteIDState
}

func (s *StateMachines) stateDeleteIDValue() (stateDelete, int64, error) {
	idStr, ok, err := s.iactr.GetUserInputAndValidate(regexGetID)

	if !ok {
		return deleteIDState, 0, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return deleteIDState, 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return deleteIDState, 0, fmt.Errorf("get %s state: %w", utils.CommandID, err)
	}

	return deleteFinishState, id, nil
}
