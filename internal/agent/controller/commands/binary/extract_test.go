package binary

import (
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
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
