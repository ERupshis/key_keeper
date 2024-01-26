package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/common/models"
)

func (s *Storage) AddRecord(record *models.Record) error {
	record.ID = s.getNextFreeIdx()
	record.UpdatedAt = time.Now()

	s.records = append(s.records, *record)
	return nil
}

func (s *Storage) AddRecords(records []models.Record) error {
	for idx := range records {
		records[idx].ID = s.getNextFreeIdx()
		records[idx].UpdatedAt = time.Now()
	}

	s.records = append(s.records, records...)
	return nil
}
