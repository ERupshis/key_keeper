package inmemory

import (
	"github.com/erupshis/key_keeper/internal/common/models"
)

func (s *Storage) GetRecord(id int64) (*models.Record, error) {
	for _, rec := range s.records {
		if rec.ID == id {
			return &rec, nil
		}
	}

	return nil, nil
}

func (s *Storage) GetAllRecords() ([]models.Record, error) {
	return s.records, nil
}

func (s *Storage) GetRecords(recordType models.RecordType, filters map[string]string) ([]models.Record, error) {
	var res []models.Record
	for idx := range s.records {
		if !canRecordBeReturned(&s.records[idx], recordType) {
			continue
		}

		if isRecordMatchToFilters(&s.records[idx], filters) {
			res = append(res, s.records[idx])
		}
	}

	return res, nil
}

func canRecordBeReturned(record *models.Record, recordType models.RecordType) bool {
	if record.Deleted {
		return false
	}

	if record.Data.RecordType != recordType && !(recordType == models.TypeAny) {
		return false
	}

	return true
}

func isRecordMatchToFilters(record *models.Record, filters map[string]string) bool {
	match := true
	if len(filters) == 0 {
		return true
	}

	for key, val := range filters {
		if key == models.StrAny {
			match = isSomeRecordMetaDataHasValue(record, val)
			if !match {
				break
			}

			continue
		}

		if metaValue, ok := record.Data.MetaData[key]; !ok || val != metaValue {
			match = false
			break
		}
	}

	return match
}

func isSomeRecordMetaDataHasValue(record *models.Record, val string) bool {
	for _, metaVal := range record.Data.MetaData {
		if val == metaVal {
			return true
		}
	}

	return false
}
