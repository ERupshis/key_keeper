package binary

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestBinary_checkHashSum(t *testing.T) {
	type fields struct {
		iactr     *interactor.Interactor
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		storePath string
	}
	type args struct {
		fileBytes []byte
		checkSum  string
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
				hash: hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("test file"),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				err: assert.NoError,
			},
		},
		{
			name: "invalid sum",
			fields: fields{
				hash: hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("test file"),
				checkSum:  "1111",
			},
			want: want{
				err: assert.Error,
			},
		},
		{
			name: "incorrect hasher key",
			fields: fields{
				hash: hasher.CreateHasher("wrong key", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("test file"),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &Binary{
				iactr:     tt.fields.iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}
			tt.want.err(t, b.checkHashSum(tt.args.fileBytes, tt.args.checkSum), fmt.Sprintf("checkHashSum(%v, %v)", tt.args.fileBytes, tt.args.checkSum))
		})
	}
}

func TestBinary_decryptFileAndValidate(t *testing.T) {
	type fields struct {
		iactr     *interactor.Interactor
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		storePath string
	}
	type args struct {
		fileBytes []byte
		checkSum  string
	}
	type want struct {
		fileBytes []byte
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
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				fileBytes: []byte("test file"),
				err:       assert.NoError,
			},
		},
		{
			name: "invalid cryptor key",
			fields: fields{
				cryptor: ska.NewSKA("wrong", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				fileBytes: nil,
				err:       assert.Error,
			},
		},
		{
			name: "invalid hasher key",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("wrong", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				fileBytes: nil,
				err:       assert.Error,
			},
		},
		{
			name: "invalid check sum",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
				checkSum:  "1111",
			},
			want: want{
				fileBytes: nil,
				err:       assert.Error,
			},
		},
		{
			name: "invalid protected file bytes",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
			},
			args: args{
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxRgs="),
				checkSum:  "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
			},
			want: want{
				fileBytes: nil,
				err:       assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &Binary{
				iactr:     tt.fields.iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}
			got, err := b.decryptFileAndValidate(tt.args.fileBytes, tt.args.checkSum)
			if !tt.want.err(t, err, fmt.Sprintf("decryptFileAndValidate(%v, %v)", tt.args.fileBytes, tt.args.checkSum)) {
				return
			}
			assert.Equalf(t, tt.want.fileBytes, got, "decryptFileAndValidate(%v, %v)", tt.args.fileBytes, tt.args.checkSum)
		})
	}
}

func TestBinary_saveFile(t *testing.T) {
	type fields struct {
		sm        *statemachines.StateMachines
		hash      *hasher.Hasher
		cryptor   *ska.SKA
		rd        *bytes.Reader
		wr        *bytes.Buffer
		storePath string
	}
	type args struct {
		record     *models.Record
		fileBytes  []byte
		pathToFile string
	}
	type want struct {
		response  []byte
		fileBytes []byte
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
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				rd:      bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr:      bytes.NewBuffer(nil),
			},
			args: args{
				record:    &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "saveFile_base.txt", SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"}}},
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
			},
			want: want{
				response:  []byte("file extracted: saveFile_base.txt\n"),
				fileBytes: []byte("test file"),
				err:       assert.NoError,
			},
		},
		{
			name: "invalid cryptor key",
			fields: fields{
				cryptor: ska.NewSKA("wrong", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				rd:      bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr:      bytes.NewBuffer(nil),
			},
			args: args{
				record:    &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "saveFile_invalid_cryptor_key.txt", SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"}}},
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
			},
			want: want{
				response:  nil,
				fileBytes: nil,
				err:       assert.Error,
			},
		},
		{
			name: "invalid hasher key",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("wrong", hasher.TypeSHA256, logger.CreateMock()),
				rd:      bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr:      bytes.NewBuffer(nil),
			},
			args: args{
				record:    &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "saveFile_invalid_hasher_key.txt", SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"}}},
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs="),
			},
			want: want{
				response:  nil,
				fileBytes: nil,
				err:       assert.Error,
			},
		},
		{
			name: "invalid protected file bytes",
			fields: fields{
				cryptor: ska.NewSKA("pass", ska.Key16),
				hash:    hasher.CreateHasher("", hasher.TypeSHA256, logger.CreateMock()),
				rd:      bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr:      bytes.NewBuffer(nil),
			},
			args: args{
				record:    &models.Record{Data: models.Data{RecordType: models.TypeBinary, Binary: &models.Binary{Name: "saveFile_invalid_hasher_key.txt", SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"}}},
				fileBytes: []byte("0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxRgs="),
			},
			want: want{
				response:  nil,
				fileBytes: nil,
				err:       assert.Error,
			},
		},
	}
	for _, tt := range tests {
		// TODO: make parallel (different files to use).
		t.Run(tt.name, func(t *testing.T) {
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			b := &Binary{
				iactr:     iactr,
				sm:        tt.fields.sm,
				hash:      tt.fields.hash,
				cryptor:   tt.fields.cryptor,
				storePath: tt.fields.storePath,
			}

			assert.NoError(t, os.WriteFile(tt.args.record.Data.Binary.SecuredFileName, tt.args.fileBytes, 0o666))
			defer func() { // TODO: bad idea to call in cycle.
				assert.NoError(t, os.Remove(tt.args.record.Data.Binary.SecuredFileName))
			}()

			err := b.saveFile(tt.args.record, tt.args.pathToFile)
			tt.want.err(t, err)
			if err != nil {
				return
			}

			writtenBytes, err := os.ReadFile(tt.args.record.Data.Binary.Name)
			assert.NoError(t, err)
			defer func() { // TODO: bad idea to call in cycle.
				assert.NoError(t, os.Remove(tt.args.record.Data.Binary.Name))
			}()
			assert.Equal(t, tt.want.fileBytes, writtenBytes)
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}
