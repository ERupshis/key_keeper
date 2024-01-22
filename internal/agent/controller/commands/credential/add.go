package credential

import (
	"errors"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (с *Credential) ProcessAddCommand(record *data.Record) error {
	record.Credentials = &data.Credentials{}
	record.RecordType = data.TypeCredentials

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: с.addMainData,
	}

	return с.sm.Add(cfg)
}

// MAIN DATA STATE MACHINE.
type addState int

const (
	addInitialState  = addState(0)
	addLoginState    = addState(1)
	addPasswordState = addState(2)
	addFinishState   = addState(3)
)

func (с *Credential) addMainData(record *data.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = с.stateInitial()
		case addLoginState:
			{
				currentState, err = с.stateLogin(record)
				if err != nil {
					return err
				}
			}
		case addPasswordState:
			{
				currentState, err = с.stateExpiration(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (с *Credential) stateInitial() addState {
	с.iactr.Printf("enter credential login: ")
	return addLoginState
}

func (с *Credential) stateLogin(record *data.Record) (addState, error) {
	credLogin, ok, err := с.iactr.GetUserInputAndValidate(nil)
	record.Credentials.Login = credLogin

	if !ok {
		return addLoginState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addLoginState, err
	}

	с.iactr.Printf("enter credential password: ")
	return addPasswordState, err

}

func (с *Credential) stateExpiration(record *data.Record) (addState, error) {
	credPassword, ok, err := с.iactr.GetUserInputAndValidate(nil)
	record.Credentials.Password = credPassword
	if !ok {
		return addPasswordState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addPasswordState, err
	}

	return addFinishState, err
}
