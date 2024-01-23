package binary

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (b *Binary) ProcessExtractCommand(record *data.Record) error {
	cfg := statemachines.ExtractConfig{
		Record:   record,
		FileSave: b.fileSave,
	}

	return b.sm.Extract(cfg)
}

func (b *Binary) fileSave(record *data.Record, pathToFile string) error {
	errMsg := "decode and save file from local storage: %w"
	fileBytes, err := os.ReadFile(filepath.Join(b.storePath, record.Binary.SecuredFileName))
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	decryptedFileBytes, err := b.cryptor.Decrypt(fileBytes)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	hashSum, err := b.hash.HashMsg(decryptedFileBytes)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if hashSum != record.Binary.SecuredFileName {
		return fmt.Errorf("hash sum is not equal")
	}

	err = os.WriteFile(filepath.Join(pathToFile, record.Binary.Name), decryptedFileBytes, 0666)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	b.iactr.Printf("file extracted: %s\n", filepath.Join(pathToFile, record.Binary.Name))
	return err
}
