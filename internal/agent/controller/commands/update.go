package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Commands) Update(parts []string, storage *inmemory.Storage) {
	if len(parts) != 1 {
		c.iactr.Printf("incorrect request. should contain command '%s' only\n", utils.CommandUpdate)
		return
	}

	err := c.handleUpdate(storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandUpdate, nil)
		return
	}

	return
}

func (c *Commands) handleUpdate(storage *inmemory.Storage) error {
	id, err := c.sm.Delete()
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	if id == nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, errs.ErrUnexpected)
	}

	return c.findAndUpdateRecordByID(*id, storage)
}

func (c *Commands) findAndUpdateRecordByID(id int64, storage *inmemory.Storage) error {
	records, err := c.getRecordByID(id, storage)
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
		err = c.bc.ProcessUpdateCommand(tmpRecord)
	case data.TypeCredentials:
		err = c.bc.ProcessUpdateCommand(tmpRecord)
	default:
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, errs.ErrIncorrectRecordType)
	}

	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	return c.confirmAndUpdateRecordByID(tmpRecord, storage)
}

func (c *Commands) confirmAndUpdateRecordByID(record *data.Record, storage *inmemory.Storage) error {
	confirmed, err := c.sm.Confirm(record, utils.CommandUpdate)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
	}

	if confirmed {
		if err = storage.UpdateRecord(record); err != nil {
			return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandUpdate, err)
		}
		c.iactr.Printf("Record sucessfully updated\n")
	} else {
		c.iactr.Printf("Record updating was interrupted by user\n")
	}

	return nil
}
