package statemachines

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

type AddState int

const (
	AddInitialState  = AddState(0)
	addMainDataState = AddState(1)
	AddFinishState   = AddState(2)
)

type AddConfig struct {
	Record   *data.Record
	MainData func(record *data.Record) error
}

func Add(cfg AddConfig) error {
	currentState := AddInitialState

	for currentState != AddFinishState {
		switch currentState {
		case AddInitialState:
			{
				if err := cfg.MainData(cfg.Record); err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return err
					} else {
						continue
					}
				}
				currentState = addMainDataState
			}
		case addMainDataState:
			{
				if err := addMetaData(cfg.Record); err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return err
					} else {
						continue
					}
				}
				currentState = AddFinishState
			}
		}
	}

	return nil
}

// MAIN DATA STATE MACHINE.
type state int

const (
	addInitialState  = state(0)
	addMetaDataState = state(1)
	addFinishState   = state(2)
)

var (
	regexMetaData = regexp.MustCompile(`^(?:[a-zA-Z0-9]+ : .+|save)$`)
)

func addMetaData(record *data.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = stateInitial()
		case addMetaDataState:
			{
				currentState, err = stateMetaData(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func stateInitial() state {
	fmt.Print("insert meta data(format: 'some_name : some_value' without quotes) or 'cancel' or 'save': ")
	return addMetaDataState
}

func stateMetaData(record *data.Record) (state, error) {
	metaData, ok, err := utils.GetUserInputAndValidate(regexMetaData)

	if metaData == utils.CommandSave {
		fmt.Printf("inserted metadata: %+v\n", record.MetaData)
		return addFinishState, err
	}

	if !ok {
		return addMetaDataState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addMetaDataState, err
	}

	parts := strings.Split(metaData, " : ")

	if record.MetaData == nil {
		record.MetaData = make(data.MetaData)
	}

	record.MetaData[parts[0]] = parts[1]
	return addInitialState, nil
}
