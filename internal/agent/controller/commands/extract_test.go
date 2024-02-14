package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCommands_handleExtract(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	type input struct {
		command string
	}
	type args struct {
		recordsInBase []models.Record
		files         map[string]string
	}
	type want struct {
		response []byte
		fileData []byte
		err      assert.ErrorAssertionFunc
		errExist assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		args  args
		want  want
	}{
		{
			name: "base",
			input: input{
				command: testutils.AddNewRow("%s" + string(filepath.Separator)),
			},
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(fmt.Sprintf(`enter absolute path to file: file extracted: %s%ctest.txt%s`, wd, filepath.Separator, "\n")),
				fileData: []byte("test file"),
				err:      assert.NoError,
				errExist: assert.NoError,
			},
		},
		{
			name: "filtered several records",
			input: input{
				command: testutils.AddNewRow("%s" + string(filepath.Separator)),
			},
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
						},
					},
					{
						ID: -2,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(`need more detailed request. (Only one record should be selected)
found '2' records:
-----
   0.  ID: -1  binary: {Name:test.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}  MetaData: map[]
   1.  ID: -2  binary: {Name:test.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}  MetaData: map[]
-----` + "\n"),
				fileData: nil,
				err:      assert.NoError,
				errExist: assert.Error,
			},
		},
		{
			name: "cancel by user",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				fileData: nil,
				err:      assert.Error,
				errExist: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, inMemoryStorage, writer := getCommands(fmt.Sprintf(tt.input.command, wd))

			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			for name, data := range tt.args.files {
				assert.NoError(t, os.WriteFile(name, []byte(data), 0o666))
			}
			defer func() {
				for name := range tt.args.files {
					assert.NoError(t, os.Remove(name))
				}
			}()

			tt.want.err(t, c.handleExtract(tt.args.recordsInBase), fmt.Sprintf("handleExtract(%v)", tt.args.recordsInBase))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")

			filePath := filepath.Join(wd, tt.args.recordsInBase[0].Data.Binary.Name)
			_, err = os.Stat(filePath)
			tt.want.errExist(t, err)

			defer func() {
				_ = os.Remove(tt.args.recordsInBase[0].Data.Binary.Name)
			}()

			if len(tt.args.recordsInBase) != 1 {
				return
			}

			fileBytes, _ := os.ReadFile(filePath)
			assert.Equal(t, tt.want.fileData, fileBytes, "extracted data fail")

		})
	}
}

func TestCommands_Extract(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err, "define working directory")

	type input struct {
		command string
	}
	type args struct {
		recordsInBase []models.Record
		parts         []string
		files         map[string]string
	}
	type want struct {
		response []byte
	}
	tests := []struct {
		name  string
		input input
		args  args
		want  want
	}{
		{
			name: "base id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) +
					testutils.AddNewRow("-1") +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(fmt.Sprintf("enter search method('id' or 'filters' or 'all'): enter record id: enter absolute path to file: file extracted: %s%ctest.txt%s", wd, filepath.Separator, "\n")),
			},
		},
		{
			name: "base filters",
			input: input{
				command: testutils.AddNewRow(utils.CommandFilters) +
					testutils.AddNewRow("key : val") +
					testutils.AddNewRow(utils.CommandContinue) +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(fmt.Sprintf("enter search method('id' or 'filters' or 'all'): enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': entered filters: map[key:val]\nenter absolute path to file: file extracted: %s%ctest.txt%s", wd, filepath.Separator, "\n")),
			},
		},
		{
			name: "base all",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll) +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(fmt.Sprintf("enter search method('id' or 'filters' or 'all'): enter absolute path to file: file extracted: %s%ctest.txt%s", wd, filepath.Separator, "\n")),
			},
		},
		{
			name: "few records via all",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll) +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
					{
						ID: -2,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte(`enter search method('id' or 'filters' or 'all'): need more detailed request. (Only one record should be selected)
found '2' records:
-----
   0.  ID: -1  binary: {Name:test.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}  MetaData: map[key:val]
   1.  ID: -2  binary: {Name:test.txt SecuredFileName:9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714}  MetaData: map[key:val]
-----` + "\n"),
			},
		},
		{
			name: "incorrect type",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) +
					testutils.AddNewRow("-1") +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrText},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'extract' and object type([bin])\n"),
			},
		},
		{
			name: "incorrect count of elems in command",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) +
					testutils.AddNewRow("-1") +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'extract' and object type([bin])\n"),
			},
		},
		{
			name: "fail to find by id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) +
					testutils.AddNewRow("-2") +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: request processing error: process 'get' command: record not found\n"),
			},
		},
		{
			name: "try to extract incorrect type",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) +
					testutils.AddNewRow("-2") +
					testutils.AddNewRow("%s"+string(filepath.Separator)),
			},
			args: args{
				parts: []string{utils.CommandExtract, models.StrBinary},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeBinary,
							Binary: &models.Binary{
								Name:            "test.txt",
								SecuredFileName: "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714",
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
					{
						ID: -2,
						Data: models.Data{
							RecordType: models.TypeCredentials,
							Credentials: &models.Credential{
								Login:    credLogin,
								Password: credPassword,
							},
							MetaData: map[string]string{
								"key": "val",
							},
						},
					},
				},
				files: map[string]string{
					"9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714": "0sAy8Kq2v5AQtBaeGXahfBV1lPpzp3Ob4HITwxMARgs=",
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: attempt to extract unsupported type 'creds'\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, inMemoryStorage, writer := getCommands(fmt.Sprintf(tt.input.command, wd))
			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			for name, data := range tt.args.files {
				assert.NoError(t, os.WriteFile(name, []byte(data), 0o666))
			}
			defer func() {
				for name := range tt.args.files {
					assert.NoError(t, os.Remove(name))
				}
			}()

			defer func() {
				_ = os.Remove(tt.args.recordsInBase[0].Data.Binary.Name)
			}()

			c.Extract(tt.args.parts, inMemoryStorage)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
