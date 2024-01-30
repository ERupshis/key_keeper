package inmemory

import (
	"encoding/json"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
)

func (s *Storage) GetAllRecordsForServer() ([]localModels.StorageRecord, error) {
	var res []localModels.StorageRecord
	for idx := range s.records {
		recordDataBytes, err := json.Marshal(s.records[idx].Data)
		if err != nil {
			return nil, fmt.Errorf("marshal record data: %w", err)
		}

		encryptedDataRecord, err := s.cryptHasher.Encrypt(recordDataBytes)
		if err != nil {
			return nil, fmt.Errorf("convert inmem record to storage rec: %w", err)
		}

		storageRecord := localModels.StorageRecord{
			ID:        s.records[idx].ID,
			Data:      encryptedDataRecord,
			Deleted:   s.records[idx].Deleted,
			UpdatedAt: s.records[idx].UpdatedAt,
		}

		res = append(res, storageRecord)
	}

	return res, nil
}

func (s *Storage) Sync(serverRecords map[int64]localModels.StorageRecord) error {
	syncedRecordsIdxs, err := s.syncLocalRecords(serverRecords)
	if err != nil {
		return err
	}

	return s.addMissingServerRecords(serverRecords, syncedRecordsIdxs)
}

func (s *Storage) syncLocalRecords(serverRecords map[int64]localModels.StorageRecord) (map[int64]struct{}, error) {
	syncedRecordsIdxs := map[int64]struct{}{}
	for idx := range s.records {
		if serverRecord, ok := serverRecords[s.records[idx].ID]; ok {
			if serverRecord.UpdatedAt.After(s.records[idx].UpdatedAt) {
				data, err := s.parseRecordData(&serverRecord)
				if err != nil {
					return nil, fmt.Errorf("sync local and server data: %w", err)
				}

				s.records[idx].Data = *data
				s.records[idx].UpdatedAt = serverRecord.UpdatedAt
				s.records[idx].Deleted = serverRecord.Deleted
			}

			syncedRecordsIdxs[serverRecord.ID] = struct{}{}
		}
	}

	return syncedRecordsIdxs, nil
}

func (s *Storage) addMissingServerRecords(serverRecords map[int64]localModels.StorageRecord, syncedRecordsIdxs map[int64]struct{}) error {
	for ID, val := range serverRecords {
		if _, ok := syncedRecordsIdxs[ID]; !ok {
			record := models.Record{
				ID:        val.ID,
				Deleted:   val.Deleted,
				UpdatedAt: val.UpdatedAt,
			}

			data, err := s.parseRecordData(&val)
			if err != nil {
				return fmt.Errorf("sync misssing server records: %w", err)
			}

			record.Data = *data

			if err = s.AddRecord(&record); err != nil {
				return fmt.Errorf("sync misssing server records: %w", err)
			}
		}
	}

	return nil
}

func (s *Storage) parseRecordData(serverRecord *localModels.StorageRecord) (*models.Data, error) {
	serverRecordDataBytes, err := s.cryptHasher.Decrypt(serverRecord.Data)
	if err != nil {
		return nil, fmt.Errorf("decrypt server record '%d': %w", serverRecord.ID, err)
	}

	data := models.Data{}
	if err = json.Unmarshal(serverRecordDataBytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshal server record data '%d': %w", serverRecord.ID, err)
	}

	return &data, nil
}
