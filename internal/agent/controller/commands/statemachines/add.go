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

type stateAdd int

const (
	addInitialState  = stateAdd(0)
	addMainDataState = stateAdd(1)
	addFinishState   = stateAdd(2)
)

type AddConfig struct {
	Record   *data.Record
	MainData func(record *data.Record) error
}

func (s *StateMachines) Add(cfg AddConfig) error {
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
				if err := s.addMetaData(cfg.Record); err != nil {
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
type stateAddMeta int

const (
	addMetaInitialState = stateAddMeta(0)
	addMetaDataState    = stateAddMeta(1)
	addMetaFinishState  = stateAddMeta(2)
)

var (
	regexMetaData = regexp.MustCompile(`^(?:[a-zA-Z0-9]+ : .+|save)$`)
)

func (s *StateMachines) addMetaData(record *data.Record) error {
	currentState := addMetaInitialState

	var err error
	for currentState != addMetaFinishState {
		switch currentState {
		case addMetaInitialState:
			currentState = s.stateMetaInitial()
		case addMetaDataState:
			{
				currentState, err = s.stateMetaData(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *StateMachines) stateMetaInitial() stateAddMeta {
	fmt.Printf(
		"insert meta data(format: 'key%svalue') or '%s' or '%s': ",
		utils.MetaSeparator,
		utils.CommandCancel,
		utils.CommandSave,
	)
	return addMetaDataState
}

func (s *StateMachines) stateMetaData(record *data.Record) (stateAddMeta, error) {
	metaData, ok, err := s.iactr.GetUserInputAndValidate(regexMetaData)

	if metaData == utils.CommandSave {
		fmt.Printf("inserted metadata: %s\n", record.MetaData)
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
