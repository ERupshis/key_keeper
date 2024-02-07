package binaries

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
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

func (bm *BinaryManager) SaveBinaries(binaries map[string][]byte) error {
	g := errgroup.Group{}
	for k, v := range binaries {
		k, v := k, v
		g.Go(func() error {
			if err := os.WriteFile(filepath.Join(bm.path, k), v, 0666); err != nil {
				return fmt.Errorf("save binary file '%s' locally: %w", k, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save binary file locally: %w", err)
	}

	return nil
}

func (bm *BinaryManager) GetFiles(binFilesList map[string]struct{}) (map[string][]byte, error) {
	g := errgroup.Group{}
	res := make(map[string][]byte)
	for k := range binFilesList {
		k := k
		g.Go(func() error {
			fileBytes, err := os.ReadFile(filepath.Join(bm.path, k))
			if err != nil {
				return fmt.Errorf("read binary file: %w", err)
			}

			res[k] = fileBytes
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
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

	g := errgroup.Group{}
	for _, fileName := range binFiles {
		fileName := fileName
		if _, ok := actualFiles[filepath.Base(fileName)]; !ok {
			g.Go(func() error {
				if err := os.Remove(fileName); err != nil {
					return fmt.Errorf("remove unused bin file: %w", err)
				}

				return nil
			})
		}
	}

	return g.Wait()
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
