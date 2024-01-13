package bankcard

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func ProcessAddCommand(record *data.Record) error {
	record.BankCard = &data.BankCard{}
	record.RecordType = data.TypeBankCard

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: addMainData,
	}

	return statemachines.Add(cfg)
}

// MAIN DATA STATE MACHINE.
type addState int

const (
	addInitialState    = addState(0)
	addNumberState     = addState(1)
	addExpirationState = addState(2)
	addCVVState        = addState(3)
	addCardHolderState = addState(4)
	addFinishState     = addState(5)
)

var (
	regexNumber         = regexp.MustCompile(`^[0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4}$`)
	regexExpirationDate = regexp.MustCompile(`^(0[1-9]|1[0-2])\/[0-9]{2}$`)
	regexCVV            = regexp.MustCompile(`^[0-9]{3,4}$`)
	regexCardHolder     = regexp.MustCompile(`^\D+$`)
)

func addMainData(record *data.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = stateInitial()
		case addNumberState:
			{
				currentState, err = stateNumber(record)
				if err != nil {
					return err
				}
			}
		case addExpirationState:
			{
				currentState, err = stateExpiration(record)
				if err != nil {
					return err
				}
			}
		case addCVVState:
			{
				currentState, err = stateCVV(record)
				if err != nil {
					return err
				}
			}
		case addCardHolderState:
			{
				currentState, err = stateCardHolder(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func stateInitial() addState {
	fmt.Print("insert card number(XXXX XXXX XXXX XXXX): ")
	return addNumberState
}

func stateNumber(record *data.Record) (addState, error) {
	cardNumber, ok, err := utils.GetUserInputAndValidate(regexNumber)
	record.BankCard.Number = cardNumber

	if !ok {
		return addNumberState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addNumberState, err
	}

	fmt.Print("insert card expiration (XX/XX): ")
	return addExpirationState, err

}

func stateExpiration(record *data.Record) (addState, error) {
	cardExpiration, ok, err := utils.GetUserInputAndValidate(regexExpirationDate)
	record.BankCard.Expiration = cardExpiration
	if !ok {
		return addExpirationState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addExpirationState, err
	}

	fmt.Print("insert card CVV (XXX or XXXX): ")
	return addCVVState, err
}

func stateCVV(record *data.Record) (addState, error) {
	cardCVV, ok, err := utils.GetUserInputAndValidate(regexCVV)
	record.BankCard.CVV = cardCVV

	if !ok {
		return addCVVState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCVVState, err
	}

	fmt.Print("insert card holder name: ")
	return addCardHolderState, err
}

func stateCardHolder(record *data.Record) (addState, error) {
	cardHolder, ok, err := utils.GetUserInputAndValidate(regexCardHolder)
	record.BankCard.Name = cardHolder

	if !ok {
		return addCardHolderState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCardHolderState, err
	}

	fmt.Printf("inserted card data: %+v\n", *record.BankCard)
	return addFinishState, err
}
