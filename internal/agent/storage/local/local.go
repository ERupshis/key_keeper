package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/common/data"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

// FileManager provides functionality to manage user data locally in the file.
type FileManager struct {
	path    string
	logs    logger.BaseLogger
	writer  *fileWriter
	scanner *fileScanner
}

// NewFileManager creates a new instance of FileManager with the specified data path and logger.
func NewFileManager(dataPath string, logger logger.BaseLogger) *FileManager {
	return &FileManager{
		path: dataPath,
		logs: logger,
	}
}

// Close closes the underlying file if open.
func (fm *FileManager) Close() error {
	return nil
}

// CheckConnection checks the file connection status and always returns true for FileManager.
func (fm *FileManager) CheckConnection(_ context.Context) (bool, error) {
	return true, nil
}

// SaveUserData saves user data in the file.
func (fm *FileManager) SaveUserData(_ context.Context, records []data.Record) error {
	if !fm.IsFileOpen() {
		if err := fm.OpenFile(fm.path, true); err != nil {
			return fmt.Errorf("cannot open file '%s' to save user data: %w", fm.path, err)
		}
		defer deferutils.ExecWithLogError(fm.CloseFile, fm.logs)
	}

	var errs []error
	for _, record := range records {
		errs = append(errs, fm.WriteRecord(&record))
	}
	return errors.Join(errs...)
}

// RestoreUserData reads user data from the file and restores it.
func (fm *FileManager) RestoreUserData(_ context.Context) ([]data.Record, error) {
	if !fm.IsFileOpen() {
		if err := fm.OpenFile(fm.path, false); err != nil {
			return nil, fmt.Errorf("cannot open file '%s' to read user data: %w", fm.path, err)
		}
		defer deferutils.ExecWithLogError(fm.CloseFile, fm.logs)
	}

	var res []data.Record
	record, err := fm.ScanRecord()
	for record != nil {
		if err != nil {
			fm.logs.Infof("failed to scan record from file '%s'", fm.path)
		} else {
			res = append(res, *record)
		}

		record, err = fm.ScanRecord()
	}

	return res, nil
}

// IsFileOpen checks if the file is open.
func (fm *FileManager) IsFileOpen() bool {
	return fm.writer != nil && fm.scanner != nil
}

// OpenFile opens or creates a file for writing or reading metrics.
func (fm *FileManager) OpenFile(path string, withTrunc bool) error {
	errMsg := "open file: %w"

	fm.path = path
	if err := os.MkdirAll(filepath.Dir(fm.path), 0755); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if err := fm.initWriter(withTrunc); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if err := fm.initScanner(); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

// CloseFile closes the file if open.
func (fm *FileManager) CloseFile() error {
	if !fm.IsFileOpen() {
		return nil
	}

	var errs []error
	errs = append(errs, fm.writer.file.Close())
	errs = append(errs, fm.scanner.file.Close())

	fm.writer = nil
	fm.scanner = nil

	return errors.Join(errs...)
}
