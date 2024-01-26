package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/models"
)

type Storage struct {
	records []models.Record

	freeIdx int64
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) getNextFreeIdx() int64 {
	s.freeIdx--
	return s.freeIdx
}
