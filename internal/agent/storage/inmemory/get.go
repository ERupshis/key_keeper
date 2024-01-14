package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (s *Storage) GetRecord(id int64) (*data.Record, error) {
	for _, rec := range s.records {
		if rec.ID == id {
			return &rec, nil
		}
	}

	return nil, nil
}

func (s *Storage) GetRecords(recordType data.RecordType, filters map[string]string) ([]data.Record, error) {
	var res []data.Record
	for idx := range s.records {
		if !canRecordBeReturned(&s.records[idx], recordType, filters) {
			continue
		}

		if isRecordMatchToFilters(&s.records[idx], filters) {
			res = append(res, s.records[idx])
		}
	}

	return res, nil
}

func canRecordBeReturned(record *data.Record, recordType data.RecordType, filters map[string]string) bool {
	if record.Deleted {
		return false
	}

	if record.RecordType != recordType && !(recordType == data.TypeAny) {
		return false
	}

	return true
}

func isRecordMatchToFilters(record *data.Record, filters map[string]string) bool {
	match := true
	if len(filters) == 0 {
		return true
	}

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
