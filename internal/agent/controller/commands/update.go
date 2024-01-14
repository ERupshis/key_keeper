package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func Update(parts []string, storage *inmemory.Storage) {
	if len(parts) != 1 {
		fmt.Printf("incorrect request. should contain command '%s' only\n", utils.CommandUpdate)
		return
	}

	err := handleUpdate(storage)
	if err != nil {
		handleCommandError(err, utils.CommandUpdate, nil)
		return
	}

	return
}

func handleUpdate(storage *inmemory.Storage) error {
	id, err := statemachines.Delete()
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	if id == nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, errs.ErrUnexpected)
	}

	return findAndUpdateRecordByID(*id, storage)
}

func findAndUpdateRecordByID(id int64, storage *inmemory.Storage) error {
	records, err := getRecordByID(id, storage)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if len(records) != 1 {
		fmt.Printf("Record with id '%d' was not found\n", id)
		return nil
	}

	tmpRecord := data.DeepCopyRecord(&records[0])
	tmpRecord.MetaData = make(data.MetaData)
	switch records[0].RecordType {
	case data.TypeBankCard:
		err = bankcard.ProcessUpdateCommand(tmpRecord)
	default:
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, errs.ErrIncorrectRecordType)
	}

	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	return confirmAndUpdateRecordByID(tmpRecord, storage)
}

func confirmAndUpdateRecordByID(record *data.Record, storage *inmemory.Storage) error {
	confirmed, err := statemachines.Confirm(record, utils.CommandUpdate)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	if confirmed {
		if err = storage.UpdateRecord(record); err != nil {
			return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
		}
		fmt.Printf("Record sucessfully updated\n")
	} else {
		fmt.Printf("Record updating was interrupted by user\n")
	}

	return nil
}
