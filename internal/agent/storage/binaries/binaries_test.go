package binaries

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fileWExt     = "file.txt"
	fileWExtData = "text"

	fileWOExt     = "wo_ext"
	fileWOExtData = "without extension"

	fileWOExtRemove1 = "wo_ext_rem1"
	fileWOExtRemove2 = "wo_ext_rem2"

	fileMissing = "missing.txt"
)

func Test_checkFile(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	files := map[string]string{
		fileWExt:  fileWExtData,
		fileWOExt: fileWOExtData,
	}

	for name, data := range files {
		assert.NoError(t, os.WriteFile(name, []byte(data), 0o666))
	}

	defer func() {
		for name := range files {
			assert.NoError(t, os.Remove(name))
		}
	}()

	type args struct {
		path     string
		binFiles *[]string
	}
	type want struct {
		binFiles []string
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base not encrypted file",
			args: args{
				path:     filepath.Join(wd, fileWExt),
				binFiles: &[]string{},
			},
			want: want{
				binFiles: []string{},
				err:      assert.NoError,
			},
		},
		{
			name: "base encrypted file",
			args: args{
				path:     filepath.Join(wd, fileWOExt),
				binFiles: &[]string{},
			},
			want: want{
				binFiles: []string{filepath.Join(wd, fileWOExt)},
				err:      assert.NoError,
			},
		},
		{
			name: "missing",
			args: args{
				path:     filepath.Join(wd, fileMissing),
				binFiles: &[]string{},
			},
			want: want{
				binFiles: []string{},
				err:      assert.Error,
			},
		},
		{
			name: "folder",
			args: args{
				path:     wd,
				binFiles: &[]string{},
			},
			want: want{
				binFiles: []string{},
				err:      assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.err(t, checkFile(tt.args.path, tt.args.binFiles), "checkFile")
			assert.True(t, reflect.DeepEqual(tt.want.binFiles, *tt.args.binFiles), "check result")
		})
	}
}

func TestBinaryManager_SyncFiles(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	type fields struct {
		path string
	}
	type args struct {
		actualFiles map[string]struct{}
	}
	type want struct {
		filesInFolder []string
		err           assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base",
			fields: fields{
				path: filepath.Join(wd, "SyncFiles"),
			},
			args: args{
				actualFiles: map[string]struct{}{
					fileWOExt: {},
				},
			},
			want: want{
				filesInFolder: []string{
					filepath.Join(wd, "SyncFiles", fileWOExt),
				},
				err: assert.NoError,
			},
		},
		{
			name: "base 2 elems",
			fields: fields{
				path: filepath.Join(wd, "SyncFiles"),
			},
			args: args{
				actualFiles: map[string]struct{}{
					fileWOExt:        {},
					fileWOExtRemove1: {},
				},
			},
			want: want{
				filesInFolder: []string{
					filepath.Join(wd, "SyncFiles", fileWOExt),
					filepath.Join(wd, "SyncFiles", fileWOExtRemove1),
				},
				err: assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filesInFolder := map[string]string{
				fileWExt:         fileWExtData,
				fileWOExt:        fileWOExtData,
				fileWOExtRemove1: fileWOExtData,
				fileWOExtRemove2: fileWOExtData,
			}

			assert.NoError(t, os.Mkdir("SyncFiles", 0o666))
			for name, data := range filesInFolder {
				assert.NoError(t, os.WriteFile(filepath.Join(wd, "SyncFiles", name), []byte(data), 0o666))
			}

			defer func() {
				_ = os.RemoveAll(filepath.Join(wd, "SyncFiles"))
			}()

			bm := &BinaryManager{
				path: tt.fields.path,
			}
			tt.want.err(t, bm.SyncFiles(tt.args.actualFiles), fmt.Sprintf("SyncFiles(%v)", tt.args.actualFiles))

			var binFiles []string
			checkFileFunc := func(path string, info os.FileInfo, err error) error {
				return checkFile(path, &binFiles)
			}
			assert.NoError(t, filepath.Walk(bm.path, checkFileFunc), "check files in directory")
			assert.True(t, reflect.DeepEqual(tt.want.filesInFolder, binFiles))
		})
	}
}

func TestBinaryManager_GetFiles(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	type fields struct {
		path string
	}
	type args struct {
		binFilesList map[string]struct{}
	}
	type want struct {
		files map[string][]byte
		err   assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base",
			fields: fields{
				path: filepath.Join(wd, "GetFiles"),
			},
			args: args{
				binFilesList: map[string]struct{}{
					fileWOExt: {},
					fileWExt:  {},
				},
			},
			want: want{
				files: map[string][]byte{
					fileWOExt: []byte(fileWOExtData),
					fileWExt:  []byte(fileWExtData),
				},
				err: assert.NoError,
			},
		},
		{
			name: "base 2",
			fields: fields{
				path: filepath.Join(wd, "GetFiles"),
			},
			args: args{
				binFilesList: map[string]struct{}{
					fileWOExt:        {},
					fileWOExtRemove1: {},
				},
			},
			want: want{
				files: map[string][]byte{
					fileWOExt:        []byte(fileWOExtData),
					fileWOExtRemove1: []byte(fileWExtData),
				},
				err: assert.NoError,
			},
		},
		{
			name: "missing file",
			fields: fields{
				path: filepath.Join(wd, "GetFiles"),
			},
			args: args{
				binFilesList: map[string]struct{}{
					fileWOExt:   {},
					fileMissing: {},
				},
			},
			want: want{
				files: nil,
				err:   assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filesInFolder := map[string]string{
				fileWExt:         fileWExtData,
				fileWOExt:        fileWOExtData,
				fileWOExtRemove1: fileWExtData,
				fileWOExtRemove2: fileWOExtData,
			}

			assert.NoError(t, os.Mkdir("GetFiles", 0o666))
			for name, data := range filesInFolder {
				assert.NoError(t, os.WriteFile(filepath.Join(wd, "GetFiles", name), []byte(data), 0o666))
			}

			defer func() {
				_ = os.RemoveAll(filepath.Join(wd, "GetFiles"))
			}()

			bm := &BinaryManager{
				path: tt.fields.path,
			}
			files, err := bm.GetFiles(tt.args.binFilesList)
			if !tt.want.err(t, err, fmt.Sprintf("GetFiles(%v)", tt.args.binFilesList)) {
				return
			}
			assert.Truef(t, reflect.DeepEqual(tt.want.files, files), "GetFiles(%v)", tt.args.binFilesList)
		})
	}
}

func TestBinaryManager_SaveBinaries(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	type fields struct {
		path string
	}
	type args struct {
		binaries          map[string][]byte
		toRemoveAfterTest []string
	}
	type want struct {
		err assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base",
			fields: fields{
				path: filepath.Join(wd, "SaveBinaries"),
			},
			args: args{
				binaries: map[string][]byte{
					fileWExt:         []byte(fileWExtData),
					fileWOExt:        []byte(fileWOExtData),
					fileWOExtRemove1: []byte(fileWExtData),
				},
				toRemoveAfterTest: []string{
					fileWExt,
					fileWOExt,
					fileWOExtRemove1,
				},
			},
			want: want{
				err: assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &BinaryManager{
				path: tt.fields.path,
			}
			assert.NoError(t, os.Mkdir(tt.fields.path, 0o666))
			tt.want.err(t, bm.SaveBinaries(tt.args.binaries), fmt.Sprintf("SaveBinaries(%v)", tt.args.binaries))

			for _, name := range tt.args.toRemoveAfterTest {
				assert.NoError(t, os.Remove(filepath.Join(tt.fields.path, name)))
			}

			assert.NoError(t, os.Remove(tt.fields.path))
		})
	}
}
