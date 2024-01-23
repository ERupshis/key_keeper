package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Commands) Extract(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{data.StrText, data.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandExtract, supportedTypes)
		return
	}

	records, err := c.handleGet(data.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandExtract, supportedTypes)
		return
	}

	if err = c.handleExtract(records); err != nil {
		c.handleCommandError(err, utils.CommandExtract, supportedTypes)
		return
	}
}

func (c *Commands) handleExtract(records []data.Record) error {
	var err error
	if len(records) == 1 {
		switch records[0].RecordType {
		case data.TypeBinary:
			err = c.binary.ProcessExtractCommand(&records[0])
		case data.TypeText:
		default:
			c.iactr.Printf("unsupported type '%s' for extracting", data.ConvertRecordTypeToString(records[0].RecordType))
		}
	} else {
		c.iactr.Printf("need more detailed request. (Only one record should be selected)\n")
		c.writeGetResult(records)
	}

	if err != nil {
		return fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandExtract, err)
	}

	return nil
}
