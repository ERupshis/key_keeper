package bankcard

import (
	"errors"
	"regexp"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (b *BankCard) ProcessAddCommand(record *models.Record) error {
	record.Data.BankCard = getBankCardDataTemplate()
	record.Data.RecordType = models.TypeBankCard

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

func (b *BankCard) addMainData(record *models.Record) error {
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

func (b *BankCard) stateInitial(record *models.Record) addState {
	b.iactr.Printf("enter card number(%s): ", record.Data.BankCard.Number)
	return addNumberState
}

func (b *BankCard) stateNumber(record *models.Record) (addState, error) {
	cardNumber, ok, err := b.iactr.GetUserInputAndValidate(regexNumber)
	if !ok {
		return addNumberState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addNumberState, err
	}

	record.Data.BankCard.Number = cardNumber
	b.iactr.Printf("enter card expiration (%s): ", record.Data.BankCard.Expiration)
	return addExpirationState, err

}

func (b *BankCard) stateExpiration(record *models.Record) (addState, error) {
	cardExpiration, ok, err := b.iactr.GetUserInputAndValidate(regexExpirationDate)
	if !ok {
		return addExpirationState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addExpirationState, err
	}

	record.Data.BankCard.Expiration = cardExpiration
	b.iactr.Printf("enter card CVV (%s): ", record.Data.BankCard.CVV)
	return addCVVState, err
}

func (b *BankCard) stateCVV(record *models.Record) (addState, error) {
	cardCVV, ok, err := b.iactr.GetUserInputAndValidate(regexCVV)
	if !ok {
		return addCVVState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCVVState, err
	}

	record.Data.BankCard.CVV = cardCVV
	if record.Data.BankCard.Name == "" {
		b.iactr.Printf("enter card holder name: ")
	} else {
		b.iactr.Printf("enter card holder name(%s): ", record.Data.BankCard.Name)
	}
	return addCardHolderState, err
}

func (b *BankCard) stateCardHolder(record *models.Record) (addState, error) {
	cardHolder, ok, err := b.iactr.GetUserInputAndValidate(regexCardHolder)
	if !ok {
		return addCardHolderState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addCardHolderState, err
	}

	record.Data.BankCard.Name = cardHolder
	b.iactr.Printf("entered card models: %+v\n", *record.Data.BankCard)
	return addFinishState, err
}

func getBankCardDataTemplate() *models.BankCard {
	return &models.BankCard{
		Number:     "XXXX XXXX XXXX XXXX",
		Expiration: "XX/XX",
		CVV:        "XXX or XXXX",
		Name:       "",
	}
}
