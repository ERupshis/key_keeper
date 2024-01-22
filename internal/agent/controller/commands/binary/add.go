package binary

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (b *Binary) ProcessAddCommand(record *data.Record) error {
	record.Binary = &data.Binary{}
	record.RecordType = data.TypeBinary

	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: b.addMainData,
	}

	return b.sm.Add(cfg)
}

// MAIN DATA STATE MACHINE.
type addState int

const (
	addInitialState  = addState(0)
	addFilePathState = addState(1)
	addFinishState   = addState(2)
)

func (b *Binary) addMainData(record *data.Record) error {
	currentState := addInitialState

	currentLocalFile := record.Binary.SecuredFileName
	var err error
	for currentState != addFinishState {
		switch currentState {
		case addInitialState:
			currentState = b.stateInitial(record)
		case addFilePathState:
			{
				currentState, err = b.stateFilePath(record)
				if err != nil {
					return err
				}
			}
		}
	}

	if currentLocalFile != "" {
		// TODO: remove old file from storage.
	}
	return nil
}

func (b *Binary) stateInitial(record *data.Record) addState {
	b.iactr.Printf("enter absolute path to file: ")
	return addFilePathState
}

func (b *Binary) stateFilePath(record *data.Record) (addState, error) {
	pathToFile, ok, err := b.iactr.GetUserInputAndValidate(nil)
	if !ok {
		return addFilePathState, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return addFilePathState, err
	}

	if !filepath.IsAbs(pathToFile) {
		b.iactr.Printf("entered local path. Try to set absolute path: ")
		return addFilePathState, nil
	}

	record.Binary.Name = filepath.Base(pathToFile)
	fileBytes, err := os.ReadFile(pathToFile)
	// TODO: add hash - it will be the name of secured file.
	// TODO: add ska securing for file data.
	// TODO: secured data need to store in the same folder as other entities data.
	//
	b.iactr.Printf("%s", string(fileBytes))
	return addFinishState, err

}
