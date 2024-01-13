package statemachines

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

const (
	GetInitialState    = state(0)
	GetSearchByID      = state(1)
	GetSearchByFilters = state(2)
	GetFinishState     = state(3)
)

func Get() (*int64, map[string]string, error) {
	currentState := GetInitialState

	var id *int64
	var filters map[string]string
	for currentState != GetFinishState {
		switch currentState {
		case GetInitialState:
			{
				method, err := getMethod()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return nil, nil, err
					} else {
						continue
					}
				}

				currentState = getStateAccordingMethod(method)
			}
		case GetSearchByID:
			{
				idTmp, err := getID()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return nil, nil, err
					} else {
						continue
					}
				}

				id = &idTmp
				currentState = GetFinishState
			}
		case GetSearchByFilters:
			{
				filtersTmp, err := getFilters()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return nil, nil, err
					} else {
						continue
					}
				}

				filters = filtersTmp
				currentState = GetFinishState
			}
		}
	}

	return id, filters, nil
}

func getStateAccordingMethod(method string) state {
	switch method {
	case utils.CommandID:
		return GetSearchByID
	case utils.CommandFilters:
		return GetSearchByFilters
	default:
		// shouldn't happen.
		return GetInitialState
	}
}

// SEARCH METHOD STATE MACHINE.
const (
	methodInitialState   = state(0)
	methodSelectionState = state(1)
	methodFinishState    = state(2)
)

var (
	regexGetMethodData = regexp.MustCompile(`^(id|filters)$`)
)

func getMethod() (string, error) {
	currentState := methodInitialState

	var method string
	var err error
	for currentState != methodFinishState {
		switch currentState {
		case methodInitialState:
			currentState = stateMethodInitial()
		case methodSelectionState:
			{
				currentState, method, err = stateMethodData()
				if err != nil {
					return "", err
				}
			}
		}
	}

	return method, nil
}

func stateMethodInitial() state {
	fmt.Printf("insert search method('%s' or '%s'): ", utils.CommandID, utils.CommandFilters)
	return methodSelectionState
}

func stateMethodData() (state, string, error) {
	method, ok, err := utils.GetUserInputAndValidate(regexGetMethodData)

	if !ok {
		return methodSelectionState, "", err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return methodSelectionState, "", err
	}

	return methodFinishState, method, nil
}

// ID STATE MACHINE.
const (
	IDInitialState = state(0)
	IDValueState   = state(1)
	IDFinishState  = state(2)
)

var (
	regexGetID = regexp.MustCompile(`^-?\d{1,10}$`)
)

func getID() (int64, error) {
	currentState := IDInitialState

	var id int64
	var err error
	for currentState != IDFinishState {
		switch currentState {
		case IDInitialState:
			currentState = stateIDInitial()
		case IDValueState:
			{
				currentState, id, err = stateIDValue()
				if err != nil {
					return 0, err
				}
			}
		}
	}

	return id, nil
}

func stateIDInitial() state {
	fmt.Printf("insert record %s: ", utils.CommandID)
	return IDValueState
}

func stateIDValue() (state, int64, error) {
	idStr, ok, err := utils.GetUserInputAndValidate(regexGetID)

	if !ok {
		return IDValueState, 0, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return IDValueState, 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return IDValueState, 0, fmt.Errorf("get %s state: %w", utils.CommandID, err)
	}

	return IDFinishState, id, nil
}

// FILTERS STATE MACHINE.
const (
	FiltersInitialState = state(0)
	FiltersValueState   = state(1)
	FiltersFinishState  = state(2)
)

var (
	regexGetFilters = regexp.MustCompile(`^(?:[a-zA-Z0-9]+ : .+|continue)$`)
)

func getFilters() (map[string]string, error) {
	currentState := FiltersInitialState

	var filters map[string]string
	var err error
	for currentState != FiltersFinishState {
		switch currentState {
		case FiltersInitialState:
			currentState = stateFiltersInitial()
		case FiltersValueState:
			{
				currentState, err = stateFiltersValue(filters)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return filters, nil
}

func stateFiltersInitial() state {
	fmt.Printf(
		"insert filters through meta data(format: 'key%svalue') or '%s' or '%s': ",
		utils.MetaSeparator,
		utils.CommandCancel,
		utils.CommandContinue,
	)
	return FiltersValueState
}

// TODO: remove duplicity?
func stateFiltersValue(filters map[string]string) (state, error) {
	metaData, ok, err := utils.GetUserInputAndValidate(regexGetFilters)

	if metaData == utils.CommandContinue {
		fmt.Printf("inserted filters: %v\n", filters)
		return FiltersInitialState, err
	}

	if !ok {
		return FiltersValueState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return FiltersValueState, err
	}

	parts := strings.Split(metaData, utils.MetaSeparator)

	if filters == nil {
		filters = make(map[string]string)
	}

	filters[parts[0]] = parts[1]
	return FiltersInitialState, nil
}
