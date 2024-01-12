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
	s.records = append(s.records, *record)
	return nil
}

func (s *Storage) AddRecords(record []data.Record) error {
	s.records = append(s.records, record...)
	return nil
}

func (s *Storage) GetRecord(id int64) (*data.Record, error) {
	for _, rec := range s.records {
		if rec.Id == id {
			return &rec, nil
		}
	}

	return nil, nil
}

func (s *Storage) GetRecords(recordType data.RecordType, filters map[string]string) ([]data.Record, error) {
	scanAllTypes := recordType == data.TypeUndefined

	var res []data.Record
	for _, rec := range s.records {
		if rec.RecordType != recordType && !scanAllTypes {
			continue
		}

		match := true
		for key, val := range filters {
			if metaValue, ok := rec.MetaData[key]; !ok || val != metaValue {
				match = false
				break
			}
		}

		if match {
			res = append(res, rec)
		}
	}

	return res, nil
}
