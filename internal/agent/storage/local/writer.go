package local

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/erupshis/key_keeper/internal/common/data"
)

// fileWriter is responsible for writing user data into the file.
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
func (fm *FileManager) WriteRecord(record *data.Record) error {
	errMsg := "write record: %w"

	if !fm.IsFileOpen() {
		return fmt.Errorf(errMsg, ErrFileIsNotOpen)
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	encryptedRecord, err := fm.cryptHasher.Encrypt(recordBytes)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	encryptedRecord = append(encryptedRecord, '\n')
	if _, err = fm.write(encryptedRecord); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	err = fm.flushWriter()
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

// write writes data to the file.
func (fm *FileManager) write(data []byte) (int, error) {
	return fm.writer.writer.Write(data)
}

// flushWriter flushes the writer buffer.
func (fm *FileManager) flushWriter() error {
	return fm.writer.writer.Flush()
}
