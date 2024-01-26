package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/models"
)

func (c *Commands) Extract(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{models.StrText, models.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandExtract, supportedTypes)
		return
	}

	records, err := c.handleGet(models.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandExtract, supportedTypes)
		return
	}

	if err = c.handleExtract(records); err != nil {
		c.handleCommandError(err, utils.CommandExtract, supportedTypes)
		return
	}
}

func (c *Commands) handleExtract(records []models.Record) error {
	var err error
	if len(records) == 1 {
		switch records[0].Data.RecordType {
		case models.TypeBinary:
			err = c.binary.ProcessExtractCommand(&records[0])
		case models.TypeText:
		default:
			c.iactr.Printf("unsupported type '%s' for extracting", models.ConvertRecordTypeToString(records[0].Data.RecordType))
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
