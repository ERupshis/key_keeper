package binary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

func (b *Binary) ProcessAddCommand(record *models.Record) error {
	record.Data.Binary = &models.Binary{}
	record.Data.RecordType = models.TypeBinary

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

func (b *Binary) addMainData(record *models.Record) error {
	currentState := addInitialState

	currentLocalFile := record.Data.Binary.SecuredFileName
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

	if currentLocalFile == record.Data.Binary.SecuredFileName {
		return nil
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

func (b *Binary) stateFilePath(record *models.Record) (addState, error) {
	errMsg := "read and secure binary models: %w"

	pathToFile, err := b.getFileNameFromUserInput()
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}

	file, err := b.getFileStream(pathToFile)
	if err != nil {
		return addFilePathState, fmt.Errorf(errMsg, err)
	}
	defer deferutils.ExecSilent(file.Close)

	fileBytes, hashSum, err := b.getFileBytesAndHashSum(file)
	if err != nil {
		return addFilePathState, fmt.Errorf("process file data: %w", err)
	}

	record.Data.Binary.Name = filepath.Base(pathToFile)
	record.Data.Binary.SecuredFileName = hashSum

	if err = b.saveEncryptedFile(fileBytes, hashSum); err != nil {
		return addFilePathState, fmt.Errorf("handle file: %w", err)
	}

	b.iactr.Printf("file saved: %+v\n", *record.Data.Binary)
	return addFinishState, nil
}

func (b *Binary) getFileNameFromUserInput() (string, error) {
	pathToFile, ok, err := b.iactr.GetUserInputAndValidate(nil)
	if !ok {
		return "", err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return "", err
	}

	return pathToFile, nil
}

func (b *Binary) getFileStream(pathToFile string) (*os.File, error) {
	if !filepath.IsAbs(pathToFile) {
		b.iactr.Printf("entered local path. Try to set absolute path: ")
		return nil, nil
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("open file to parse: %w", err)
	}

	return file, nil
}

func (b *Binary) getFileBytesAndHashSum(file *os.File) ([]byte, string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("read file data %w", err)
	}

	hashSum, err := b.hash.HashMsg(fileBytes)
	if err != nil {
		return nil, "", fmt.Errorf("calculate file data hashsum %w", err)
	}

	return fileBytes, hashSum, nil
}

func (b *Binary) saveEncryptedFile(fileBytes []byte, hashSum string) error {
	encryptedFileBytes, err := b.cryptor.Encrypt(fileBytes)
	if err != nil {
		return fmt.Errorf("encrypt file data: %w", err)
	}

	err = os.WriteFile(filepath.Join(b.storePath, hashSum), encryptedFileBytes, 0666)
	if err != nil {
		return fmt.Errorf("write encypted data in storage file: %w", err)
	}

	return nil
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
