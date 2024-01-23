package binary

import (
	"errors"
	"fmt"
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
			currentState = b.stateInitial()
		case addFilePathState:
			{
				currentState, err = b.stateFilePath(record)
				if err != nil {
					return err
				}
			}
		}
	}

	if err = b.removeOldSecuredFile(currentLocalFile); err != nil {
		return err
	}

	return nil
}

func (b *Binary) stateInitial() addState {
	b.iactr.Printf("enter absolute path to file: ")
	return addFilePathState
}

func (b *Binary) stateFilePath(record *data.Record) (addState, error) {
	errMsg := "read anf secure binary data: %w"
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

	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}

	hashSum, err := b.hash.HashMsg(fileBytes)
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}

	record.Binary.Name = filepath.Base(pathToFile)
	record.Binary.SecuredFileName = hashSum

	encryptedFileBytes, err := b.cryptor.Encrypt(fileBytes)
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}

	err = os.WriteFile(filepath.Join(b.storePath, hashSum), encryptedFileBytes, 0666)
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}

	b.iactr.Printf("file saved: %+v\n", *record.Binary)
	return addFinishState, err
}

func (b *Binary) removeOldSecuredFile(fileName string) error {
	if fileName == "" {
		return nil
	}

	err := os.Remove(filepath.Join(b.storePath, fileName))
	if err != nil {
		return fmt.Errorf("remove old secured file: %w", err)
	}

	return nil
}
