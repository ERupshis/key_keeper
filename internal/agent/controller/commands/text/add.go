package text

import (
	"errors"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/models"
)

func (t *Text) ProcessAddCommand(record *models.Record) error {
	record.Data.Text = &models.Text{}
	record.Data.RecordType = models.TypeText

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: t.addMainData,
	}

	return t.sm.Add(cfg)
}

// MAIN DATA STATE MACHINE.
type addState int

const (
	addInitialState = addState(0)
	addDataState    = addState(1)
	addFinishState  = addState(2)
)

func (t *Text) addMainData(record *models.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = t.stateInitial()
		case addDataState:
			{
				currentState, err = t.stateData(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Text) stateInitial() addState {
	t.iactr.Printf("enter text to save: ")
	return addDataState
}

func (t *Text) stateData(record *models.Record) (addState, error) {
	text, ok, err := t.iactr.GetUserInputAndValidate(nil)
	record.Data.Text.Data = text

	if !ok {
		return addDataState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addDataState, err
	}

	t.iactr.Printf("entered credential models: %+v\n", *record.Data.Text)
	return addFinishState, err

}
