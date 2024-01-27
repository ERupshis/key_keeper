package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/models"
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

func (s *Storage) getNextFreeIdx() int64 {
	s.freeIdx--
	return s.freeIdx
}
