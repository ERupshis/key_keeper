package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
)

var (
	zeroTime = time.Time{}
)

func (s *Storage) AddRecord(record *models.Record) error {
	if record.ID <= 0 {
		record.ID = s.getNextFreeIdx()
	}

	if record.UpdatedAt == zeroTime {
		record.UpdatedAt = time.Now()
	}

	s.records = append(s.records, *record)
	return nil
}
