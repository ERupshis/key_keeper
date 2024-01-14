package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/data"
)

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
