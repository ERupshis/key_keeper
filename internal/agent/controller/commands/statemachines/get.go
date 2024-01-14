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

type stateGet int

const (
	getInitialState    = stateGet(0)
	getSearchByID      = stateGet(1)
	getSearchByFilters = stateGet(2)
	getFinishState     = stateGet(3)
)

func Get() (*int64, map[string]string, error) {
	currentState := getInitialState

	var id *int64
	var filters map[string]string
	for currentState != getFinishState {
		switch currentState {
		case getInitialState:
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
		case getSearchByID:
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
				currentState = getFinishState
			}
		case getSearchByFilters:
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
				currentState = getFinishState
			}
		}
	}

	return id, filters, nil
}

func getStateAccordingMethod(method string) stateGet {
	switch method {
	case utils.CommandID:
		return getSearchByID
	case utils.CommandFilters:
		return getSearchByFilters
	default:
		// shouldn't happen.
		return getInitialState
	}
}

// SEARCH METHOD STATE MACHINE.
type stateGetMethod int

const (
	getMethodInitialState   = stateGetMethod(0)
	getMethodSelectionState = stateGetMethod(1)
	getMethodFinishState    = stateGetMethod(2)
)

var (
	regexGetMethodData = regexp.MustCompile(`^(id|filters)$`)
)

func getMethod() (string, error) {
	currentState := getMethodInitialState

	var method string
	var err error
	for currentState != getMethodFinishState {
		switch currentState {
		case getMethodInitialState:
			currentState = stateGetMethodInitial()
		case getMethodSelectionState:
			{
				currentState, method, err = stateGetMethodData()
				if err != nil {
					return "", err
				}
			}
		}
	}

	return method, nil
}

func stateGetMethodInitial() stateGetMethod {
	fmt.Printf("insert search method('%s' or '%s'): ", utils.CommandID, utils.CommandFilters)
	return getMethodSelectionState
}

func stateGetMethodData() (stateGetMethod, string, error) {
	method, ok, err := utils.GetUserInputAndValidate(regexGetMethodData)

	if !ok {
		return getMethodSelectionState, "", err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return getMethodSelectionState, "", err
	}

	return getMethodFinishState, method, nil
}

// ID STATE MACHINE.
type stateGetID int

const (
	getIDInitialState = stateGetID(0)
	getIDValueState   = stateGetID(1)
	getIDFinishState  = stateGetID(2)
)

var (
	regexGetID = regexp.MustCompile(`^-?\d{1,10}$`)
)

func getID() (int64, error) {
	currentState := getIDInitialState

	var id int64
	var err error
	for currentState != getIDFinishState {
		switch currentState {
		case getIDInitialState:
			currentState = stateGetIDInitial()
		case getIDValueState:
			{
				currentState, id, err = stateGetIDValue()
				if err != nil {
					return 0, err
				}
			}
		}
	}

	return id, nil
}

func stateGetIDInitial() stateGetID {
	fmt.Printf("insert record %s: ", utils.CommandID)
	return getIDValueState
}

func stateGetIDValue() (stateGetID, int64, error) {
	idStr, ok, err := utils.GetUserInputAndValidate(regexGetID)

	if !ok {
		return getIDValueState, 0, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return getIDValueState, 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return getIDValueState, 0, fmt.Errorf("get %s state: %w", utils.CommandID, err)
	}

	return getIDFinishState, id, nil
}

// FILTERS STATE MACHINE.
type stateGetFilters int

const (
	getFiltersInitialState = stateGetFilters(0)
	getFiltersValueState   = stateGetFilters(1)
	getFiltersFinishState  = stateGetFilters(2)
)

var (
	regexGetFilters = regexp.MustCompile(`^(?:[a-zA-Z0-9]+ : .+|continue)$`)
)

func getFilters() (map[string]string, error) {
	currentState := getFiltersInitialState

	filters := make(map[string]string)
	var err error
	for currentState != getFiltersFinishState {
		switch currentState {
		case getFiltersInitialState:
			currentState = stateGetFiltersInitial()
		case getFiltersValueState:
			{
				currentState, err = stateGetFiltersValue(filters)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return filters, nil
}

func stateGetFiltersInitial() stateGetFilters {
	fmt.Printf(
		"insert filters through meta data(format: 'key%svalue') or '%s' or '%s': ",
		utils.MetaSeparator,
		utils.CommandCancel,
		utils.CommandContinue,
	)
	return getFiltersValueState
}

func stateGetFiltersValue(filters map[string]string) (stateGetFilters, error) {
	metaData, ok, err := utils.GetUserInputAndValidate(regexGetFilters)

	if metaData == utils.CommandContinue {
		fmt.Printf("inserted filters: %s\n", filters)
		return getFiltersFinishState, err
	}

	if !ok {
		return getFiltersValueState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return getFiltersValueState, err
	}

	parts := strings.Split(metaData, utils.MetaSeparator)
	filters[parts[0]] = parts[1]
	return getFiltersInitialState, nil
}
