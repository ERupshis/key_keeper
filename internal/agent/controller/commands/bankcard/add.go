package bankcard

import (
	"errors"
	"regexp"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (b *BankCard) ProcessAddCommand(record *data.Record) error {
	record.BankCard = &data.BankCard{
		Number:     "XXXX XXXX XXXX XXXX",
		Expiration: "XX/XX",
		CVV:        "XXX or XXXX",
		Name:       "",
	}
	record.RecordType = data.TypeBankCard

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: b.addMainData,
	}

	return b.sm.Add(cfg)
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

func (b *BankCard) addMainData(record *data.Record) error {
	currentState := addInitialState

	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = b.stateInitial(record)
		case addNumberState:
			{
				currentState, err = b.stateNumber(record)
				if err != nil {
					return err
				}
			}
		case addExpirationState:
			{
				currentState, err = b.stateExpiration(record)
				if err != nil {
					return err
				}
			}
		case addCVVState:
			{
				currentState, err = b.stateCVV(record)
				if err != nil {
					return err
				}
			}
		case addCardHolderState:
			{
				currentState, err = b.stateCardHolder(record)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (b *BankCard) stateInitial(record *data.Record) addState {
	b.iactr.Printf("insert card number(%s): ", record.BankCard.Number)
	return addNumberState
}

func (b *BankCard) stateNumber(record *data.Record) (addState, error) {
	cardNumber, ok, err := b.iactr.GetUserInputAndValidate(regexNumber)
	record.BankCard.Number = cardNumber

	if !ok {
		return addNumberState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addNumberState, err
	}

	b.iactr.Printf("insert card expiration (%s): ", record.BankCard.Expiration)
	return addExpirationState, err

}

func (b *BankCard) stateExpiration(record *data.Record) (addState, error) {
	cardExpiration, ok, err := b.iactr.GetUserInputAndValidate(regexExpirationDate)
	record.BankCard.Expiration = cardExpiration
	if !ok {
		return addExpirationState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addExpirationState, err
	}

	b.iactr.Printf("insert card CVV (%s): ", record.BankCard.CVV)
	return addCVVState, err
}

func (b *BankCard) stateCVV(record *data.Record) (addState, error) {
	cardCVV, ok, err := b.iactr.GetUserInputAndValidate(regexCVV)
	record.BankCard.CVV = cardCVV

	if !ok {
		return addCVVState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCVVState, err
	}

	if record.BankCard.Name == "" {
		b.iactr.Printf("insert card holder name: ")
	} else {
		b.iactr.Printf("insert card holder name(%s): ", record.BankCard.Name)
	}
	return addCardHolderState, err
}

func (b *BankCard) stateCardHolder(record *data.Record) (addState, error) {
	cardHolder, ok, err := b.iactr.GetUserInputAndValidate(regexCardHolder)
	record.BankCard.Name = cardHolder

	if !ok {
		return addCardHolderState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCardHolderState, err
	}

	b.iactr.Printf("inserted card data: %+v\n", *record.BankCard)
	return addFinishState, err
}
