package credential

import (
	"errors"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Credential) ProcessAddCommand(record *data.Record) error {
	record.Credentials = &data.Credential{}
	record.RecordType = data.TypeCredentials

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: c.addMainData,
	}

	return c.sm.Add(cfg)
}

// MAIN DATA STATE MACHINE.
type addState int

const (
	addInitialState  = addState(0)
	addLoginState    = addState(1)
	addPasswordState = addState(2)
	addFinishState   = addState(3)
)

func (c *Credential) addMainData(record *data.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = c.stateInitial()
		case addLoginState:
			{
				currentState, err = c.stateLogin(record)
				if err != nil {
					return err
				}
			}
		case addPasswordState:
			{
				currentState, err = c.stateExpiration(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Credential) stateInitial() addState {
	c.iactr.Printf("enter credential login: ")
	return addLoginState
}

func (c *Credential) stateLogin(record *data.Record) (addState, error) {
	credLogin, ok, err := c.iactr.GetUserInputAndValidate(nil)
	record.Credentials.Login = credLogin

	if !ok {
		return addLoginState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addLoginState, err
	}

	c.iactr.Printf("enter credential password: ")
	return addPasswordState, err

}

func (c *Credential) stateExpiration(record *data.Record) (addState, error) {
	credPassword, ok, err := c.iactr.GetUserInputAndValidate(nil)
	record.Credentials.Password = credPassword
	if !ok {
		return addPasswordState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addPasswordState, err
	}

	c.iactr.Printf("entered credential data: %+v\n", *record.Credentials)
	return addFinishState, err
}
