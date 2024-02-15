package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (s *Storage) UpdateRecord(record *models.Record) error {
	var updated bool
	for idx := range s.records {
		if s.records[idx].ID == record.ID {
			record.UpdatedAt = time.Now()
			s.records[idx] = *record
			updated = true
			break
		}
	}

	if !updated {
		return ErrRecordNotFound
	}

	return nil
}
