package binary

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFile        = "test_file.txt"
	anotherTestFile = "another_test_file.txt"
)

func TestBinary_removeOldSecuredFile(t *testing.T) {
	type fields struct {
		fileName string
	}
	type args struct {
		fileName string
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
				fileName: "test_removeOldSecuredFile.txt",
			},
			args: args{
				fileName: "test_removeOldSecuredFile.txt",
			},
			want: want{
				err: assert.NoError,
			},
		},
		{
			name: "empty file name",
			fields: fields{
				fileName: "test_removeOldSecuredFile.txt",
			},
			args: args{
				fileName: "",
			},
			want: want{
				err: assert.NoError,
			},
		},
		{
			name: "missing file",
			fields: fields{
				fileName: "test_removeOldSecuredFile.txt",
			},
			args: args{
				fileName: "missing_file",
			},
			want: want{
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binary{}
			assert.NoError(t, os.WriteFile(tt.fields.fileName, []byte("some text"), 0o666))
			defer func() {
				_ = os.Remove(tt.fields.fileName)
			}()

			tt.want.err(t, b.removeOldSecuredFile(tt.args.fileName), fmt.Sprintf("removeOldSecuredFile(%v)", tt.args.fileName))
		})
	}
}

func TestBinary_saveEncryptedFile(t *testing.T) {
	type fields struct {
		iactr     *interactor.Interactor
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		storePath string
	}
	type args struct {
		fileBytes []byte
		hashSum   string
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
				cryptor: ska.NewSKA("pass", ska.Key16),
			},
			args: args{
				fileBytes: []byte("test_file"),
				hashSum:   "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				err: assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binary{
				iactr:     tt.fields.iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}
			tt.want.err(t, b.saveEncryptedFile(tt.args.fileBytes, tt.args.hashSum), fmt.Sprintf("saveEncryptedFile(%v, %v)", tt.args.fileBytes, tt.args.hashSum))
			assert.NoError(t, os.Remove(tt.args.hashSum))
		})
	}
}

func TestBinary_getFileBytesAndHashSum(t *testing.T) {
	type fields struct {
		iactr     *interactor.Interactor
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		storePath string
	}
	type args struct {
		file      string
		fileBytes []byte
	}
	type want struct {
		fileBytes []byte
		hashSum   string
		err       assert.ErrorAssertionFunc
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
				hash: hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				file:      "test_file.txt",
				fileBytes: []byte("test file"),
			},
			want: want{
				fileBytes: []byte("test file"),
				hashSum:   "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
				err:       assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binary{
				iactr:     tt.fields.iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}

			assert.NoError(t, os.WriteFile(tt.args.file, tt.args.fileBytes, 0o666))
			defer func() {
				assert.NoError(t, os.Remove(tt.args.file))
			}()

			fileStream, err := os.Open(tt.args.file)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, fileStream.Close())
			}()

			fileBytes, hashSum, err := b.getFileBytesAndHashSum(fileStream)
			if !tt.want.err(t, err, fmt.Sprintf("getFileBytesAndHashSum(%v)", tt.args.file)) {
				return
			}
			assert.Equalf(t, tt.want.fileBytes, fileBytes, "getFileBytesAndHashSum(%v)", tt.args.file)
			assert.Equalf(t, tt.want.hashSum, hashSum, "getFileBytesAndHashSum(%v)", tt.args.file)
		})
	}
}

func TestBinary_getFileStream(t *testing.T) {
	type fields struct {
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		rd        *bytes.Reader
		wr        *bytes.Buffer
		storePath string
	}
	type args struct {
		absPath  bool
		fileName string
		fileData []byte
	}
	type want struct {
		response   []byte
		shouldNill bool
		err        assert.ErrorAssertionFunc
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
			},
			args: args{
				absPath:  true,
				fileName: "test_file.txt",
				fileData: []byte("test file"),
			},
			want: want{
				response:   nil,
				shouldNill: false,
				err:        assert.NoError,
			},
		},
		{
			name: "not abs path",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
			},
			args: args{
				absPath:  false,
				fileName: "test_file.txt",
				fileData: []byte("test file"),
			},
			want: want{
				response:   []byte("entered local path. Try to set absolute path: "),
				shouldNill: true,
				err:        assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			b := &Binary{
				iactr:     iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}

			pathToFile := tt.args.fileName
			if tt.args.absPath {
				wd, err := os.Getwd()
				assert.NoError(t, err)
				pathToFile = filepath.Join(wd, tt.args.fileName)
			}

			assert.NoError(t, os.WriteFile(pathToFile, tt.args.fileData, 0o666))
			defer func() {
				assert.NoError(t, os.Remove(pathToFile))
			}()

			got, err := b.getFileStream(pathToFile)
			if !tt.want.err(t, err, fmt.Sprintf("getFileStream(%v)", pathToFile)) {
				return
			}

			assert.Equal(t, tt.want.shouldNill, got == nil)
			if got != nil {
				assert.NoError(t, got.Close())
			}

			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBinary_getFileNameFromUserInput(t *testing.T) {
	type fields struct {
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		rd        *bytes.Reader
		wr        *bytes.Buffer
		storePath string
	}
	type want struct {
		response []byte
		input    string
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "base",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(testFile))),
				wr: bytes.NewBuffer(nil),
			},
			want: want{
				response: nil,
				input:    testFile,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel by user",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
			},
			want: want{
				response: nil,
				input:    "",
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte(testFile)),
				wr: bytes.NewBuffer(nil),
			},
			want: want{
				response: nil,
				input:    "",
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			b := &Binary{
				iactr:     iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}
			got, err := b.getFileNameFromUserInput()
			if !tt.want.err(t, err, "getFileNameFromUserInput()") {
				return
			}
			assert.Equalf(t, tt.want.input, got, "getFileNameFromUserInput()")
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBinary_stateFilePath(t *testing.T) {
	type fields struct {
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		rd        *bytes.Reader
		wr        *bytes.Buffer
		storePath string
	}
	type args struct {
		input    string
		record   *models.Record
		fileName string
		fileData []byte
	}
	type want struct {
		response []byte
		state    addState
		err      assert.ErrorAssertionFunc
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
				rd:      nil,
				wr:      bytes.NewBuffer(nil),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				cryptor: ska.NewSKA("pass", ska.Key16),
			},
			args: args{
				input:    testutils.AddNewRow(testFile),
				record:   &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "", SecuredFileName: ""}}},
				fileName: testFile,
				fileData: []byte("test file"),
			},
			want: want{
				response: []byte("file saved: {Name:test_file.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}\n"),
				state:    addFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel by user",
			fields: fields{
				rd:      nil,
				wr:      bytes.NewBuffer(nil),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				cryptor: ska.NewSKA("pass", ska.Key16),
			},
			args: args{
				input:    testutils.AddNewRow(utils.CommandCancel),
				record:   &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "", SecuredFileName: ""}}},
				fileName: testFile,
				fileData: []byte("test file"),
			},
			want: want{
				response: nil,
				state:    addFilePathState,
				err:      assert.Error,
			},
		},
		{
			name: "oef",
			fields: fields{
				rd:      nil,
				wr:      bytes.NewBuffer(nil),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				cryptor: ska.NewSKA("pass", ska.Key16),
			},
			args: args{
				input:    testFile,
				record:   &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "", SecuredFileName: ""}}},
				fileName: testFile,
				fileData: []byte("test file"),
			},
			want: want{
				response: nil,
				state:    addFilePathState,
				err:      assert.Error,
			},
		},
		{
			name: "missing file",
			fields: fields{
				rd:      nil,
				wr:      bytes.NewBuffer(nil),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				cryptor: ska.NewSKA("pass", ska.Key16),
			},
			args: args{
				input:    testutils.AddNewRow(testFile),
				record:   &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "", SecuredFileName: ""}}},
				fileName: anotherTestFile,
				fileData: []byte("test file"),
			},
			want: want{
				response: nil,
				state:    addFilePathState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, err := os.Getwd()
			require.NoError(t, err)
			reader := bytes.NewReader([]byte(filepath.Join(wd, tt.args.input)))
			iactr := testutils.CreateUserInteractor(reader, tt.fields.wr, logger.CreateMock())

			b := &Binary{
				iactr:     iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}

			assert.NoError(t, os.WriteFile(tt.args.fileName, tt.args.fileData, 0o666))
			defer func() {
				assert.NoError(t, os.Remove(tt.args.fileName))
			}()

			got, err := b.stateFilePath(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateFilePath(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateFilePath(%v)", tt.args.record)
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
			_ = os.Remove(tt.args.record.Data.Binary.SecuredFileName)
		})
	}
}

func TestBinary_stateInitial(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type want struct {
		response []byte
		state    addState
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "base",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer([]byte{}),
				sm: nil,
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				state:    addFilePathState,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			cred := &Binary{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			assert.Equalf(t, tt.want.state, cred.stateInitial(), "stateInitial()")
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBinary_addMainData(t *testing.T) {
	type fields struct {
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		storePath string
	}
	type args struct {
		input  string
		record *models.Record
	}
	type rawFile struct {
		name string
		data []byte
	}
	type existingFile struct {
		name string
		data []byte
	}
	type want struct {
		response []byte
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		rawFile      rawFile
		existingFile existingFile
		want         want
	}{
		{
			name: "base",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				input:  testutils.AddNewRow(testFile),
				record: &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{}}},
			},
			rawFile: rawFile{
				name: testFile,
				data: []byte("test file"),
			},
			want: want{
				response: []byte("enter absolute path to file: file saved: {Name:test_file.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "file updated with the same file hashSum",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				input:  testutils.AddNewRow(testFile),
				record: &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"}}},
			},
			rawFile: rawFile{
				name: testFile,
				data: []byte("test file"),
			},
			existingFile: existingFile{
				name: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
				data: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs"),
			},
			want: want{
				response: []byte("enter absolute path to file: file saved: {Name:test_file.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "remove old file assigned to record",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				input:  testutils.AddNewRow(testFile),
				record: &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{SecuredFileName: "old_file_hashSum"}}},
			},
			rawFile: rawFile{
				name: testFile,
				data: []byte("test file"),
			},
			existingFile: existingFile{
				name: "old_file_hashSum",
				data: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs"),
			},
			want: want{
				response: []byte("enter absolute path to file: file saved: {Name:test_file.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				input:  testutils.AddNewRow(utils.CommandCancel),
				record: &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{}}},
			},
			rawFile: rawFile{
				name: testFile,
				data: []byte("test file"),
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, err := os.Getwd()
			require.NoError(t, err)
			reader := bytes.NewReader([]byte(filepath.Join(wd, tt.args.input)))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			b := &Binary{
				iactr:     iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}

			assert.NoError(t, os.WriteFile(tt.rawFile.name, tt.rawFile.data, 0o666))
			defer func() {
				assert.NoError(t, os.Remove(tt.rawFile.name))
			}()

			if len(tt.existingFile.name) != 0 {
				assert.NoError(t, os.WriteFile(tt.existingFile.name, tt.existingFile.data, 0o666))
			}

			tt.want.err(t, b.addMainData(tt.args.record), fmt.Sprintf("addMainData(%v)", tt.args.record))
			assert.Equal(t, tt.want.response, writer.Bytes())
			_ = os.Remove(tt.args.record.Data.Binary.SecuredFileName)
		})
	}
}
