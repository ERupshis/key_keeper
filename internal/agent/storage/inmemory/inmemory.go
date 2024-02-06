package inmemory

import (
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
)

type Storage struct {
	records []models.Record

	cryptHasher *ska.SKA
	freeIdx     int64
}

func NewStorage(cryptHasher *ska.SKA) *Storage {
	return &Storage{
		cryptHasher: cryptHasher,
	}
}
func (s *Storage) resetNextFreeIdx() {
	minIdx := int64(0)
	for _, record := range s.records {
		minIdx = min(minIdx, record.ID)
	}

	s.freeIdx = minIdx
}

func (s *Storage) getNextFreeIdx() int64 {
	s.freeIdx--
	return s.freeIdx
}
