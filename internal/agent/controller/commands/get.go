package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func Get(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{data.StrAny, data.StrCredentials, data.StrBankCard, data.StrText, data.StrBinary}
	if len(parts) != 2 {
		fmt.Printf("incorrect request. should contain command '%s' and object type(%v)\n", utils.CommandGet, supportedTypes)
		return
	}

	records, err := handleGet(data.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		handleCommandError(err, utils.CommandGet, supportedTypes)
		return
	}

	writeGetResult(records)
	return
}

func writeGetResult(records []data.Record) {
	if len(records) == 0 {
		fmt.Printf("missing record(s)\n")
	} else {
		fmt.Printf("found '%d' records:\n", len(records))
		for idx, record := range records {
			fmt.Printf("   %d. %+v\n", idx, record)
		}
	}
}

func handleGet(recordType data.RecordType, storage *inmemory.Storage) ([]data.Record, error) {
	if recordType == data.TypeUndefined {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrIncorrectRecordType)
	}

	id, filters, err := statemachines.Get()
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if id != nil {
		return getRecordByID(*id, storage)
	}

	if filters != nil {
		return getRecordByFilters(recordType, filters, storage)
	}

	return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrUnexpected)
}

func getRecordByID(id int64, storage *inmemory.Storage) ([]data.Record, error) {
	record, err := storage.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if record == nil {
		return nil, nil
	}

	return []data.Record{*record}, nil
}

func getRecordByFilters(recordType data.RecordType, filters map[string]string, storage *inmemory.Storage) ([]data.Record, error) {
	records, err := storage.GetRecords(recordType, filters)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	return records, nil
}
