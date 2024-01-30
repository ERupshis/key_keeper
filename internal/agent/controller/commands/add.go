package commands

import (
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	models2 "github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

func (c *Commands) Add(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{models2.StrCredentials, models2.StrBankCard, models2.StrText, models2.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandAdd, supportedTypes)
		return
	}

	record, err := c.handleAdd(models2.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandAdd, supportedTypes)
		return
	}

	c.iactr.Printf("record added: %s\n", record)
}

func (c *Commands) handleCommandError(err error, command string, supportedTypes []string) {
	if errors.Is(err, errs.ErrInterruptedByUser) {
		c.iactr.Printf("'%s' command was canceled by user\n", command)
		return
	}

	c.iactr.Printf("request processing error: %v", err)
	if errors.Is(err, errs.ErrIncorrectRecordType) || errors.Is(err, errs.ErrIncorrectServerActionType) {
		c.iactr.Printf(". only (%s) are supported", supportedTypes)
	}

	c.iactr.Printf("\n")
}

func (c *Commands) handleAdd(recordType models2.RecordType, storage *inmemory.Storage) (*models2.Record, error) {
	newRecord := &models2.Record{
		ID: -1,
	}

	var err error
	switch recordType {
	case models2.TypeBankCard:
		err = c.bc.ProcessAddCommand(newRecord)
	case models2.TypeCredentials:
		err = c.creds.ProcessAddCommand(newRecord)
	case models2.TypeText:
		err = c.text.ProcessAddCommand(newRecord)
	case models2.TypeBinary:
		err = c.binary.ProcessAddCommand(newRecord)
	default:
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandAdd, errs.ErrIncorrectRecordType)
	}

	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandAdd, err)
	}

	if err = storage.AddRecord(newRecord); err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandAdd, err)
	}

	return newRecord, nil
}
