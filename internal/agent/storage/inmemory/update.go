package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/common/models"
)

func (s *Storage) UpdateRecord(record *models.Record) error {
	for idx := range s.records {
		if s.records[idx].ID == record.ID {
			record.UpdatedAt = time.Now()
			s.records[idx] = *record
			break
		}
	}

	return nil
}
