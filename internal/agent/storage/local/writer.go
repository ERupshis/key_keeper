package local

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/models"
)

// fileWriter is responsible for writing user models into the file.
type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// initWriter initializes the file writer.
func (fm *FileManager) initWriter(withTrunc bool) error {
	var flag int
	flag = os.O_WRONLY | os.O_CREATE
	if withTrunc {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(fm.path, flag, 0666)
	if err != nil {
		return fmt.Errorf("init writer: %w", err)
	}

	fm.writer = &fileWriter{file, bufio.NewWriter(file)}
	return nil
}

// WriteRecord writes a user record into the file.
func (fm *FileManager) WriteRecord(record *models.Record) error {
	errMsg := "write record: %w"

	if !fm.IsFileOpen() {
		return fmt.Errorf(errMsg, ErrFileIsNotOpen)
	}

	recordDataBytes, err := json.Marshal(record.Data)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	encryptedDataRecord, err := fm.cryptHasher.Encrypt(recordDataBytes)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	storageRecord := localModels.StorageRecord{
		ID:        record.ID,
		Data:      encryptedDataRecord,
		Deleted:   record.Deleted,
		UpdatedAt: record.UpdatedAt,
	}

	storageRecordBytes, err := json.Marshal(storageRecord)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	storageRecordBytes = append(storageRecordBytes, '\n')
	if _, err = fm.write(storageRecordBytes); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	err = fm.flushWriter()
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

// write writes models to the file.
func (fm *FileManager) write(data []byte) (int, error) {
	return fm.writer.writer.Write(data)
}

// flushWriter flushes the writer buffer.
func (fm *FileManager) flushWriter() error {
	return fm.writer.writer.Flush()
}
