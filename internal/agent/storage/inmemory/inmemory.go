package inmemory

import (
	"time"

	"github.com/erupshis/key_keeper/internal/common/data"
)

type Storage struct {
	records []data.Record
	// TODO: need to added last local id and refresh it during syncing.
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) AddRecord(record *data.Record) error {
	record.UpdatedAt = time.Now()
	s.records = append(s.records, *record)
	return nil
}

func (s *Storage) AddRecords(records []data.Record) error {
	for idx := range records {
		records[idx].UpdatedAt = time.Now()
	}

	s.records = append(s.records, records...)
	return nil
}

func (s *Storage) DeleteRecord(id int64) error {
	for idx, rec := range s.records {
		if rec.Id == id {
			s.records = append(s.records[:idx], s.records[idx+1:]...)
			return nil
		}
	}

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
	scanAllTypes := recordType == data.TypeAny

	var res []data.Record
	for _, rec := range s.records {
		if rec.RecordType != recordType && !scanAllTypes {
			continue
		}

		if len(filters) == 0 {
			res = append(res, rec)
			continue
		}

		if isRecordMatchToFilters(&rec, filters) {
			res = append(res, rec)
		}
	}

	return res, nil
}

func isRecordMatchToFilters(record *data.Record, filters map[string]string) bool {
	match := true
	for key, val := range filters {
		if key == data.StrAny {
			match = isSomeRecordMetaDataHasValue(record, val)
			if !match {
				break
			}

			continue
		}

		if metaValue, ok := record.MetaData[key]; !ok || val != metaValue {
			match = false
			break
		}
	}

	return match
}

func isSomeRecordMetaDataHasValue(record *data.Record, val string) bool {
	for _, metaVal := range record.MetaData {
		if val == metaVal {
			return true
		}
	}

	return false
}
