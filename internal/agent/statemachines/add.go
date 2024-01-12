package statemachines

import (
	"errors"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

type AddState int

const (
	AddInitialState  = AddState(0)
	addDataState     = AddState(1)
	AddMetaDataState = AddState(2)
	AddFinishState   = AddState(3)
)

type AddConfig struct {
	Record   *data.Record
	MainData func(record *data.Record) error
	MetaData func(record *data.Record) error // TODO: common state machine for all types.
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
				currentState = addDataState
			}
		case addDataState:
			{
				if err := cfg.MetaData(cfg.Record); err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return err
					} else {
						continue
					}
				}
				currentState = AddMetaDataState
			}
		case AddMetaDataState:
			{
				// TODO: some info for user?
				currentState = AddFinishState
			}
		}
	}

	return nil
}
