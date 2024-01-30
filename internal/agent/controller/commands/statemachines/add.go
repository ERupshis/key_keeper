package statemachines

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type stateAdd int

const (
	addInitialState  = stateAdd(0)
	addMainDataState = stateAdd(1)
	addFinishState   = stateAdd(2)
)

type AddConfig struct {
	Record   *models.Record
	MainData func(record *models.Record) error
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

func (s *StateMachines) addMetaData(record *models.Record) error {
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
		"enter meta models(format: 'key%svalue') or '%s' or '%s': ",
		utils.MetaSeparator,
		utils.CommandCancel,
		utils.CommandSave,
	)
	return addMetaDataState
}

func (s *StateMachines) stateMetaData(record *models.Record) (stateAddMeta, error) {
	metaData, ok, err := s.iactr.GetUserInputAndValidate(regexMetaData)

	if metaData == utils.CommandSave {
		fmt.Printf("entered metadata: %s\n", record.Data.MetaData)
		return addMetaFinishState, err
	}

	if !ok {
		return addMetaDataState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addMetaDataState, err
	}

	parts := strings.Split(metaData, utils.MetaSeparator)

	if record.Data.MetaData == nil {
		record.Data.MetaData = make(models.MetaData)
	}

	record.Data.MetaData[parts[0]] = parts[1]
	return addMetaInitialState, nil
}
