package local

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/models"
)

// fileScanner is responsible for scanning user models from the file.
type fileScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

// initScanner initializes the file scanner.
func (fm *FileManager) initScanner() error {
	file, err := os.OpenFile(fm.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("init scanner: %w", err)
	}

	fm.scanner = &fileScanner{file, bufio.NewScanner(file)}
	return nil
}

// ScanRecord scans and returns a user models record from the file.
func (fm *FileManager) ScanRecord() (*models.Record, error) {
	errMsg := "scan record: %w"

	if !fm.IsFileOpen() {
		return nil, fmt.Errorf(errMsg, ErrFileIsNotOpen)
	}

	if isScanOk, err := fm.scan(); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	} else if !isScanOk {
		return nil, nil
	}

	var storageRecord localModels.StorageRecord
	storageRecordBytes := fm.scannedBytes()
	if err := json.Unmarshal(storageRecordBytes, &storageRecord); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	storageRecordDataBytes, err := fm.cryptHasher.Decrypt(storageRecord.Data)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	record := models.Record{
		ID:        storageRecord.ID,
		Deleted:   storageRecord.Deleted,
		UpdatedAt: storageRecord.UpdatedAt,
	}

	if err = json.Unmarshal(storageRecordDataBytes, &record.Data); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return &record, nil
}

// scan scans the file for the next line.
func (fm *FileManager) scan() (bool, error) {
	if !fm.scanner.scanner.Scan() {
		if err := fm.scanner.scanner.Err(); err != nil {
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}

// scannedBytes returns the scanned bytes from the scanner.
func (fm *FileManager) scannedBytes() []byte {
	return fm.scanner.scanner.Bytes()
}
