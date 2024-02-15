package inmemory

import (
	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (s *Storage) RestoreRecords(records []models.Record) error {
	s.records = append(s.records, records...)
	s.resetNextFreeIdx()
	return nil
}
