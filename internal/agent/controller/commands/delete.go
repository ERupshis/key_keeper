package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/models"
)

func (c *Commands) Delete(parts []string, storage *inmemory.Storage) {
	if len(parts) != 1 {
		c.iactr.Printf("incorrect request. should contain command '%s' only\n", utils.CommandDelete)
		return
	}

	err := c.handleDelete(storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandDelete, nil)
		return
	}
}

func (c *Commands) handleDelete(storage *inmemory.Storage) error {
	id, err := c.sm.Delete()
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if id != nil {
		return c.findAndDeleteRecordByID(*id, storage)
	}

	return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, errs.ErrUnexpected)
}

func (c *Commands) findAndDeleteRecordByID(id int64, storage *inmemory.Storage) error {
	records, err := c.getRecordByID(id, storage)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if len(records) != 1 {
		c.iactr.Printf("Record with id '%d' was not found\n", id)
		return nil
	}

	return c.confirmAndDeleteByID(&records[0], storage)
}

func (c *Commands) confirmAndDeleteByID(record *models.Record, storage *inmemory.Storage) error {
	confirmed, err := c.sm.Confirm(record, utils.CommandDelete)
	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
	}

	if confirmed {
		if err = storage.DeleteRecord(record.ID); err != nil {
			return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandDelete, err)
		}
		c.iactr.Printf("Record sucessfully deleted\n")
	} else {
		c.iactr.Printf("Record deleting was interrupted by user\n")
	}

	return nil
}
