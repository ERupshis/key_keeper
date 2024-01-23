package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Commands) Get(parts []string, storage *inmemory.Storage) {
	supportedTypes := []string{data.StrAny, data.StrCredentials, data.StrBankCard, data.StrText, data.StrBinary}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and object type(%s)\n", utils.CommandGet, supportedTypes)
		return
	}

	records, err := c.handleGet(data.ConvertStringToRecordType(parts[1]), storage)
	if err != nil {
		c.handleCommandError(err, utils.CommandGet, supportedTypes)
		return
	}

	c.writeGetResult(records)
}

func (c *Commands) writeGetResult(records []data.Record) {
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

func (c *Commands) handleGet(recordType data.RecordType, storage *inmemory.Storage) ([]data.Record, error) {
	if recordType == data.TypeUndefined {
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

func (c *Commands) getRecordByID(id int64, storage *inmemory.Storage) ([]data.Record, error) {
	record, err := storage.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	if record == nil {
		return nil, nil
	}

	return []data.Record{*record}, nil
}

func (c *Commands) getRecordByFilters(recordType data.RecordType, filters map[string]string, storage *inmemory.Storage) ([]data.Record, error) {
	records, err := storage.GetRecords(recordType, filters)
	if err != nil {
		return nil, fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandGet, err)
	}

	return records, nil
}
