package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/data"
)

type Storage struct {
	records []data.Record
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) AddRecord(record *data.Record) error {
	return nil
}

func (s *Storage) AddRecords(record []data.Record) error {
	return nil
}

func (s *Storage) Record() (*data.Record, error) {
	return nil, nil
}

func (s *Storage) Records() ([]data.Record, error) {
	return nil, nil
}
