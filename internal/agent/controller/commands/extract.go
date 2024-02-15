package commands

import (
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

func (c *Commands) Extract(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{models.StrBinary}
	if len(parts) != 2 || parts[1] != models.StrBinary {
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
			c.iactr.Printf("attempt to extract unsupported type '%s'\n", models.ConvertRecordTypeToString(records[0].Data.RecordType))
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
