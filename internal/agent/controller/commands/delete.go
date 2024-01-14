package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func Delete(parts []string, storage *inmemory.Storage) {
	if len(parts) != 1 {
		fmt.Printf("incorrect request. should contain command '%s' only\n", utils.CommandDelete)
		return
	}

	err := handleDelete(storage)
	if err != nil {
		handleCommandError(err, utils.CommandDelete, nil)
		return
	}

	return
}

func handleDelete(storage *inmemory.Storage) error {
	id, err := statemachines.Delete()
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if id != nil {
		return findAndDeleteRecordByID(*id, storage)
	}

	return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, errs.ErrUnexpected)
}

func findAndDeleteRecordByID(id int64, storage *inmemory.Storage) error {
	records, err := getRecordByID(id, storage)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if len(records) != 1 {
		fmt.Printf("Record with id '%d' was not found\n", id)
		return nil
	}

	return confirmAndDeleteByID(&records[0], storage)
}

func confirmAndDeleteByID(record *data.Record, storage *inmemory.Storage) error {
	confirmed, err := statemachines.Confirm(record, utils.CommandDelete)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if confirmed {
		if err = storage.DeleteRecord(record.ID); err != nil {
			return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
		}
		fmt.Printf("Record sucessfully deleted\n")
	} else {
		fmt.Printf("Record deleting was interrupted by user\n")
	}

	return nil
}
