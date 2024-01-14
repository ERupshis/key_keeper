package inmemory

import (
	"time"
)

func (s *Storage) DeleteRecord(id int64) error {
	for idx, rec := range s.records {
		if rec.ID == id {
			if id < 0 {
				s.records = append(s.records[:idx], s.records[idx+1:]...)
			} else {
				s.records[idx].Deleted = true
				s.records[idx].UpdatedAt = time.Now()
			}
			return nil
		}
	}

	return nil
}
