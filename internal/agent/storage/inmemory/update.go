package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/common/data"
)

func (s *Storage) UpdateRecord(record *data.Record) error {
	for idx := range s.records {
		if s.records[idx].Id == record.Id {
			record.UpdatedAt = time.Now()
			s.records[idx] = *record
			break
		}
	}

	return nil
}
