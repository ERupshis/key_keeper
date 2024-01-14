package statemachines

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

type stateConfirm int

const (
	confirmInitialState = stateConfirm(0)
	confirmApproveState = stateConfirm(1)
	confirmFinishState  = stateConfirm(2)
)

var (
	regexConfirmApprove = regexp.MustCompile(`^(yes|no)$`)
)

func Confirm(record *data.Record, command string) (bool, error) {
	currentState := confirmInitialState

	var confirmed bool
	for currentState != confirmFinishState {
		switch currentState {
		case confirmInitialState:
			{
				currentState = stateConfirmInitial(record, command)
			}
		case confirmApproveState:
			{
				currentStateTmp, approve, err := stateConfirmApprove()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return false, err
					} else {
						continue
					}
				}

				confirmed = approve
				currentState = currentStateTmp
			}
		}
	}

	return confirmed, nil
}

func stateConfirmInitial(record *data.Record, command string) stateConfirm {
	switch command {
	case utils.CommandDelete:
		fmt.Printf("Do you really want to permanently delete the record '%s'(yes/no): ", record)
	case utils.CommandUpdate:
		fmt.Printf("Do you really want to update the record '%s'(yes/no): ", record)
	default:
		fmt.Printf("Do you really want to commit action with record '%s'(yes/no): ", record)
	}

	return confirmApproveState
}

func stateConfirmApprove() (stateConfirm, bool, error) {
	approve, ok, err := utils.GetUserInputAndValidate(regexConfirmApprove)

	if !ok {
		return confirmApproveState, false, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return confirmApproveState, false, err
	}

	return confirmFinishState, approve == utils.CommandYes, nil
}