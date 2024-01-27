package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	models2 "github.com/erupshis/key_keeper/internal/models"
)

func (c *Commands) Get(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{models2.StrAny, models2.StrCredentials, models2.StrBankCard, models2.StrText, models2.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandGet, supportedTypes)
		return
	}

	records, err := c.handleGet(models2.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandGet, supportedTypes)
		return
	}

	c.writeGetResult(records)
}

func (c *Commands) writeGetResult(records []models2.Record) {
	if len(records) == 0 {
		c.iactr.Printf("missing record(s)\n")
	} else {
		c.iactr.Printf("found '%d' records:\n", len(records))
		c.iactr.Printf("-----\n")

		w := tabwriter.NewWriter(c.iactr.Writer(), 0, 0, 2, ' ', 0)
		for idx, record := range records {
			_, _ = fmt.Fprint(w, "   ", idx, ".", record.TabString(), "\n")
		}
		_ = w.Flush()

		c.iactr.Printf("-----\n")
	}
}

func (c *Commands) handleGet(recordType models2.RecordType, storage *inmemory.Storage) ([]models2.Record, error) {
	if recordType == models2.TypeUndefined {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrIncorrectRecordType)
	}

	id, filters, err := c.sm.Get()
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if id != nil {
		return c.getRecordByID(*id, storage)
	}

	if filters != nil {
		return c.getRecordByFilters(recordType, filters, storage)
	}

	return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, errs.ErrUnexpected)
}

func (c *Commands) getRecordByID(id int64, storage *inmemory.Storage) ([]models2.Record, error) {
	record, err := storage.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if record == nil {
		return nil, nil
	}

	return []models2.Record{*record}, nil
}

func (c *Commands) getRecordByFilters(recordType models2.RecordType, filters map[string]string, storage *inmemory.Storage) ([]models2.Record, error) {
	records, err := storage.GetRecords(recordType, filters)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	return records, nil
}
