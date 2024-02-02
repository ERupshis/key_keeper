package binaries

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BinaryManager struct {
	path string
}

func NewBinaryManager(path string) *BinaryManager {
	return &BinaryManager{path: path}
}

func (bm *BinaryManager) SetPath(newPath string) {
	bm.path = newPath
}

func (bm *BinaryManager) SaveBinaries(binaries map[string][]byte) error { // TODO:need to add goroutine and return channel.
	for k, v := range binaries {
		err := os.WriteFile(filepath.Join(bm.path, k), v, 0666)
		if err != nil {
			return fmt.Errorf("save binary file '%s' locally: %w", k, err)
		}
	}

	return nil
}

func (bm *BinaryManager) GetFiles(binFilesList map[string]struct{}) (map[string][]byte, error) {
	res := make(map[string][]byte)

	for k, _ := range binFilesList {
		fileBytes, err := os.ReadFile(filepath.Join(bm.path, k))
		if err != nil {
			return nil, fmt.Errorf("read binary file: %w", err)
		}

		res[k] = fileBytes
	}

	return res, nil
}

func (bm *BinaryManager) SyncFiles(actualFiles map[string]struct{}) error {
	var binFiles []string

	checkFileFunc := func(path string, info os.FileInfo, err error) error {
		return checkFile(path, &binFiles)
	}

	if err := filepath.Walk(bm.path, checkFileFunc); err != nil {
		return fmt.Errorf("create list of bin files: %w", err)
	}

	for _, fileName := range binFiles {
		if _, ok := actualFiles[filepath.Base(fileName)]; !ok {
			if err := os.Remove(fileName); err != nil {
				return fmt.Errorf("remove unused bin file: %w", err)
			}
		}
	}

	return nil
}

func checkFile(path string, binFiles *[]string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("get file meta: %w", err)
	}

	if !info.IsDir() {
		if !strings.Contains(info.Name(), ".") {
			*binFiles = append(*binFiles, path)
		}
	}

	return nil
}
