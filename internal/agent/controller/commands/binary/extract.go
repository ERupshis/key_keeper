package binary

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (b *Binary) ProcessExtractCommand(record *models.Record) error {
	cfg := statemachines.ExtractConfig{
		Record:   record,
		FileSave: b.fileSave,
	}

	return b.sm.Extract(cfg)
}

func (b *Binary) fileSave(record *models.Record, pathToFile string) error {
	errMsg := "decode and save file from local storage: %w"
	fileBytes, err := os.ReadFile(filepath.Join(b.storePath, record.Data.Binary.SecuredFileName))
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	decryptedFileBytes, err := b.decryptFileAndValidate(fileBytes, record.Data.Binary.SecuredFileName)
	if err != nil {
		return fmt.Errorf("parse protected file: %w", err)
	}

	err = os.WriteFile(filepath.Join(pathToFile, record.Data.Binary.Name), decryptedFileBytes, 0666)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	b.iactr.Printf("file extracted: %s\n", filepath.Join(pathToFile, record.Data.Binary.Name))
	return err
}

func (b *Binary) decryptFileAndValidate(fileBytes []byte, checkSum string) ([]byte, error) {
	decryptedFileBytes, err := b.cryptor.Decrypt(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("decrypt file data: %w", err)
	}

	if err = b.checkHashSum(decryptedFileBytes, checkSum); err != nil {
		return nil, err
	}

	return decryptedFileBytes, nil
}

func (b *Binary) checkHashSum(fileBytes []byte, checkSum string) error {
	hashSum, err := b.hash.HashMsg(fileBytes)
	if err != nil {
		return fmt.Errorf("calculate data hashsum: %w", err)
	}

	if hashSum != checkSum {
		return ErrHashSumInvalid
	}

	return nil
}
