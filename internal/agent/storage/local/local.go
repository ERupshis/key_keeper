package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/storage/binaries"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/ticker"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

const (
	KeyStorageName = "key_keeper_strg"
)

// AutoSaveConfig auto save settings.
type AutoSaveConfig struct {
	SaveInterval    time.Duration
	InMemoryStorage *inmemory.Storage
	BinaryManager   *binaries.BinaryManager
	Logs            logger.BaseLogger
}

// FileManager provides functionality to manage user models locally in the file.
type FileManager struct {
	path       string
	passPhrase string

	iactr *interactor.Interactor

	logs    logger.BaseLogger
	writer  *fileWriter
	scanner *fileScanner

	cryptHasher *ska.SKA
	autoSaveCfg *AutoSaveConfig
}

// NewFileManager creates a new instance of FileManager with the specified models path and logger.
func NewFileManager(dataPath string, logger logger.BaseLogger, iactr *interactor.Interactor, autoSaveCfg *AutoSaveConfig, cryptHasher *ska.SKA) *FileManager {
	return &FileManager{
		path:        dataPath + KeyStorageName,
		logs:        logger,
		iactr:       iactr,
		autoSaveCfg: autoSaveCfg,
		cryptHasher: cryptHasher,
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

// SaveUserData saves user models in the file.
func (fm *FileManager) SaveUserData(records []models.Record) error {
	if !fm.IsFileOpen() {
		if err := fm.OpenFile(fm.path, true); err != nil {
			return fmt.Errorf("cannot open file '%s' to save user models: %w", fm.path, err)
		}
		defer deferutils.ExecWithLogError(fm.CloseFile, fm.logs)
	}

	var errs []error
	for _, record := range records {
		errs = append(errs, fm.WriteRecord(&record))
	}
	return errors.Join(errs...)
}

// RestoreUserData reads user models from the file and restores it.
func (fm *FileManager) RestoreUserData(ctx context.Context) ([]models.Record, error) {
	if !fm.IsFileOpen() {
		if err := fm.OpenFile(fm.path, false); err != nil {
			return nil, fmt.Errorf("cannot open file '%s' to read user models: %w", fm.path, err)
		}
		defer deferutils.ExecWithLogError(fm.CloseFile, fm.logs)
	}

	errMsg := "restore storage: %w"
	var res []models.Record
	record, err := fm.ScanRecord()
	for record != nil {
		if err != nil {
			fm.logs.Infof("failed to scan record from file '%s'", fm.path)
		} else {
			res = append(res, *record)
		}

		record, err = fm.ScanRecord()
	}

	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	fm.RunAutoSave(ctx)
	return res, nil
}

// IsFileOpen checks if the file is open.
func (fm *FileManager) IsFileOpen() bool {
	return fm.writer != nil && fm.scanner != nil
}

func (fm *FileManager) IsFileExist() (bool, error) {
	fileStats, err := os.Stat(fm.path)

	if err == nil && fileStats != nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("check local storage file existense: %w", err)
	}
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

func (fm *FileManager) SetPassPhrase(newPassPhrase string) {
	fm.passPhrase = newPassPhrase
	fm.cryptHasher.SetAESKey(newPassPhrase, ska.Key16)
}

func (fm *FileManager) Path() string {
	return fm.path
}

func (fm *FileManager) SetPath(newPath string) {
	fm.path = newPath + KeyStorageName
}

func (fm *FileManager) RunAutoSave(ctx context.Context) {
	storeTicker := time.NewTicker(fm.autoSaveCfg.SaveInterval)
	go ticker.Run(storeTicker, ctx, func() {
		select {
		case <-ctx.Done():
			storeTicker.Stop()
			return
		default:
			fm.saveRecords()
		}

	})
}

func (fm *FileManager) saveRecords() {
	records, err := fm.autoSaveCfg.InMemoryStorage.GetAllRecords()
	if err != nil {
		fm.autoSaveCfg.Logs.Infof("failed to extract inmemory models, error: %v", err)
	}

	if err = fm.SaveUserData(records); err != nil {
		fm.autoSaveCfg.Logs.Infof("failed to save models in local storage, error: %v", err)
	}
}

func (fm *FileManager) SyncBinaries() {
	fm.autoSaveCfg.BinaryManager.SetPath(filepath.Dir(fm.path))
	actualFiles := fm.autoSaveCfg.InMemoryStorage.GetBinFilesList()

	actualFiles[KeyStorageName] = struct{}{}
	if err := fm.autoSaveCfg.BinaryManager.SyncFiles(actualFiles); err != nil {
		fm.autoSaveCfg.Logs.Infof("sync actual binaries, error: %v", err)
	}
}
