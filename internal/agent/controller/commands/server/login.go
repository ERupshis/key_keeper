package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/models"
)

func (s *Server) ProcessLoginCommand(ctx context.Context) error {
	creds := models.Credential{}
	if err := s.collectCreds(&creds); err != nil {
		return fmt.Errorf("collect credentials: %w", err)
	}

	if err := s.client.Login(ctx, &creds); err != nil {
		return fmt.Errorf("login on server: %w", err)
	}

	return nil
}

// MAIN DATA STATE MACHINE.
type loginState int

const (
	addInitialState  = loginState(0)
	addLoginState    = loginState(1)
	addPasswordState = loginState(2)
	addFinishState   = loginState(3)
)

func (s *Server) collectCreds(creds *models.Credential) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = s.stateInitial()
		case addLoginState:
			{
				currentState, err = s.stateLogin(creds)
				if err != nil {
					return err
				}
			}
		case addPasswordState:
			{
				currentState, err = s.statePassword(creds)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Server) stateInitial() loginState {
	s.iactr.Printf("enter login: ")
	return addLoginState
}

func (s *Server) stateLogin(creds *models.Credential) (loginState, error) {
	credLogin, ok, err := s.iactr.GetUserInputAndValidate(nil)
	creds.Login = credLogin

	if !ok {
		return addLoginState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addLoginState, err
	}

	s.iactr.Printf("enter password: ")
	return addPasswordState, err

}

func (s *Server) statePassword(creds *models.Credential) (loginState, error) {
	credPassword, ok, err := s.iactr.GetUserInputAndValidate(nil)
	creds.Password = credPassword
	if !ok {
		return addPasswordState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addPasswordState, err
	}

	s.iactr.Printf("entered credentials: %+v\n", *creds)
	return addFinishState, err
}
