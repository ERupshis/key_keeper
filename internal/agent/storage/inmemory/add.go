package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (s *Storage) AddRecord(record *models.Record) error {
	if record.ID <= 0 {
		record.ID = s.getNextFreeIdx()
	}

	if record.UpdatedAt == time.UnixMilli(0) {
		record.UpdatedAt = time.Now()
	}

	s.records = append(s.records, *record)
	return nil
}
