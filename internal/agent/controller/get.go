package controller

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/controller/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Controller) handleGetCommand(parts []string) (bool, bool) {
	supportedTypes := []string{data.StrAny, data.StrCredentials, data.StrBankCard, data.StrText, data.StrBinary}
	if len(parts) != 2 {
		fmt.Printf("incorrect request. should contain command '%s' and object type(%v)\n", utils.CommandGet, supportedTypes)
		return true, false
	}

	records, err := c.processGetCommand(data.ConvertStringToRecordType(parts[1]))
	if err != nil {
		c.handleProcessError(err, utils.CommandGet, supportedTypes)
		return false, false
	}

	c.writeGetResult(records)
	return true, true
}

func (c *Controller) writeGetResult(records []data.Record) {
	if len(records) == 0 {
		fmt.Printf("missing record with given id\n")
	} else {
		fmt.Printf("found '%d' records:\n", len(records))
		for idx, record := range records {
			fmt.Printf("   %d. %+v\n", idx, record)
		}
	}
}

func (c *Controller) processGetCommand(recordType data.RecordType) ([]data.Record, error) {
	if recordType == data.TypeUndefined {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrIncorrectRecordType)
	}

	id, filters, err := statemachines.Get()
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if id != nil {
		return c.getRecordByID(*id)
	}

	if filters != nil {
		return c.getRecordByFilters(recordType, filters)
	}

	return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrUnexpected)
}

func (c *Controller) getRecordByID(id int64) ([]data.Record, error) {
	record, err := c.inmemory.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if record == nil {
		return nil, nil
	}

	return []data.Record{*record}, nil
}

func (c *Controller) getRecordByFilters(recordType data.RecordType, filters map[string]string) ([]data.Record, error) {
	records, err := c.inmemory.GetRecords(recordType, filters)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	return records, nil
}
