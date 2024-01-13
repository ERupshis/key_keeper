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

type state int

const (
	addInitialState  = state(0)
	addMainDataState = state(1)
	addFinishState   = state(2)
)

type AddConfig struct {
	Record   *data.Record
	MainData func(record *data.Record) error
}

func Add(cfg AddConfig) error {
	currentState := addInitialState

	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
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
				currentState = addFinishState
			}
		}
	}

	return nil
}

// MAIN DATA STATE MACHINE.

const (
	addMetaInitialState = state(0)
	addMetaDataState    = state(1)
	addMetaFinishState  = state(2)
)

var (
	regexMetaData = regexp.MustCompile(`^(?:[a-zA-Z0-9]+ : .+|save)$`)
)

func addMetaData(record *data.Record) error {
	currentState := addMetaInitialState

	var err error
	for currentState != addMetaFinishState {
		switch currentState {
		case addMetaInitialState:
			currentState = stateMetaInitial()
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

func stateMetaInitial() state {
	fmt.Printf(
		"insert meta data(format: 'key%svalue') or '%s' or '%s': ",
		utils.MetaSeparator,
		utils.CommandCancel,
		utils.CommandSave,
	)
	return addMetaDataState
}

func stateMetaData(record *data.Record) (state, error) {
	metaData, ok, err := utils.GetUserInputAndValidate(regexMetaData)

	if metaData == utils.CommandSave {
		fmt.Printf("inserted metadata: %v\n", record.MetaData)
		return addMetaFinishState, err
	}

	if !ok {
		return addMetaDataState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addMetaDataState, err
	}

	parts := strings.Split(metaData, utils.MetaSeparator)

	if record.MetaData == nil {
		record.MetaData = make(data.MetaData)
	}

	record.MetaData[parts[0]] = parts[1]
	return addMetaInitialState, nil
}
