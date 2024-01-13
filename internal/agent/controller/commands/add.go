package commands

import (
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func Add(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{data.StrCredentials, data.StrBankCard, data.StrText, data.StrBinary}
	if len(parts) != 2 {
		fmt.Printf("incorrect request. should contain command '%s' and object type(%v)\n", utils.CommandAdd, supportedTypes)
		return
	}

	record, err := processAddCommand(data.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		handleProcessError(err, utils.CommandAdd, supportedTypes)
		return
	}

	fmt.Printf("record added: %+v\n", record)
	return
}

func handleProcessError(err error, command string, supportedTypes []string) {
	if errors.Is(err, errs.ErrInterruptedByUser) {
		fmt.Printf("'%s' command was canceled by user\n", command)
		return
	}

	fmt.Printf("request parsing error: %v\n", err)
	if errors.Is(err, errs.ErrIncorrectRecordType) {
		fmt.Printf("only (%v) are supported\n", supportedTypes)
	}
}

func processAddCommand(recordType data.RecordType, storage *inmemory.Storage) (*data.Record, error) {
	newRecord := &data.Record{
		Id: -1,
	}

	var err error
	switch recordType {
	case data.TypeBankCard:
		err = bankcard.ProcessAddCommand(newRecord)
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
