package commands

import (
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Commands) Add(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{data.StrCredentials, data.StrBankCard, data.StrText, data.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandAdd, supportedTypes)
		return
	}

	record, err := c.handleAdd(data.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandAdd, supportedTypes)
		return
	}

	c.iactr.Printf("record added: %s\n", record)
	return
}

func (c *Commands) handleCommandError(err error, command string, supportedTypes []string) {
	if errors.Is(err, errs.ErrInterruptedByUser) {
		c.iactr.Printf("'%s' command was canceled by user\n", command)
		return
	}

	c.iactr.Printf("request parsing error: %v", err)
	if errors.Is(err, errs.ErrIncorrectRecordType) {
		c.iactr.Printf(". only (%s) are supported", supportedTypes)
	}
	c.iactr.Printf("\n")
}

func (c *Commands) handleAdd(recordType data.RecordType, storage *inmemory.Storage) (*data.Record, error) {
	newRecord := &data.Record{
		ID: -1,
	}

	var err error
	switch recordType {
	case data.TypeBankCard:
		err = c.bc.ProcessAddCommand(newRecord)
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
