package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/common/data"
)

func (s *Storage) AddRecord(record *data.Record) error {
	record.Id = s.getNextFreeIdx()
	record.UpdatedAt = time.Now()
	s.records = append(s.records, *record)
	return nil
}

func (s *Storage) AddRecords(records []data.Record) error {
	for idx := range records {
		records[idx].Id = s.getNextFreeIdx()
		records[idx].UpdatedAt = time.Now()
	}

	s.records = append(s.records, records...)
	return nil
}
