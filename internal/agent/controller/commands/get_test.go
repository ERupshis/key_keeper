package commands

import (
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	metaKey    = "key"
	metaVal    = "val"
	metaKeyVal = "key : val"
)

func TestCommands_getRecordByFilters(t *testing.T) {

	type args struct {
		recordsInBase []models.Record
		recordType    models.RecordType
		filters       map[string]string
	}
	type want struct {
		response []byte
		records  []models.Record
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base empty filters",
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
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
						},
					},
				},
				recordType: models.TypeText,
				filters:    map[string]string{},
			},
			want: want{
				response: nil,
				records: []models.Record{
					{ID: -1},
				},
				err: assert.NoError,
			},
		},
		{
			name: "several records by filters",
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
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
						},
					},
				},
				recordType: models.TypeText,
				filters:    map[string]string{},
			},
			want: want{
				response: nil,
				records: []models.Record{
					{ID: -1},
				},
				err: assert.NoError,
			},
		},
		{
			name: "select any",
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
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
						},
					},
				},
				recordType: models.TypeAny,
				filters:    map[string]string{},
			},
			want: want{
				response: nil,
				records: []models.Record{
					{ID: -1},
					{ID: -2},
				},
				err: assert.NoError,
			},
		},
		{
			name: "by filters",
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
				recordType: models.TypeAny,
				filters:    map[string]string{metaKey: metaVal},
			},
			want: want{
				response: nil,
				records: []models.Record{
					{ID: -2},
				},
				err: assert.NoError,
			},
		},
		{
			name: "by filters and type",
			args: args{
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
				recordType: models.TypeText,
				filters:    map[string]string{metaKey: metaVal},
			},
			want: want{
				response: nil,
				records:  []models.Record{},
				err:      assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands("")

			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			records, err := c.getRecordByFilters(tt.args.recordType, tt.args.filters, inMemoryStorage)
			if !tt.want.err(t, err, fmt.Sprintf("getRecordByFilters(%v, %v, %v)", tt.args.recordType, tt.args.filters, inMemoryStorage)) {
				return
			}

			require.Equal(t, len(tt.want.records), len(records), "records count is not equal to expected")
			for idx := range tt.want.records {
				assert.Equal(t, tt.want.records[idx].ID, records[idx].ID)
			}

			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_getRecordByID(t *testing.T) {
	type args struct {
		id            int64
		recordsInBase []models.Record
	}
	type want struct {
		response []byte
		records  []models.Record
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base",
			args: args{
				id: -1,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: nil,
				records: []models.Record{
					{ID: -1},
				},
				err: assert.NoError,
			},
		},
		{
			name: "missing id",
			args: args{
				id: -3,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: nil,
				records:  []models.Record{},
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands("")

			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			records, err := c.getRecordByID(tt.args.id, inMemoryStorage)
			if !tt.want.err(t, err, fmt.Sprintf("getRecordByID(%v, %v)", tt.args.id, inMemoryStorage)) {
				return
			}

			require.Equal(t, len(tt.want.records), len(records), "records count is not equal to expected")
			for idx := range tt.want.records {
				assert.Equal(t, tt.want.records[idx].ID, records[idx].ID)
			}

			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_handleGet(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		recordType    models.RecordType
		recordsInBase []models.Record
	}
	type want struct {
		response []byte
		records  []models.Record
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		args  args
		want  want
	}{
		{
			name: "base by id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) + testutils.AddNewRow("-1"),
			},
			args: args{
				recordType: models.TypeText,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: "),
				records: []models.Record{
					{ID: -1},
				},
				err: assert.NoError,
			},
		},
		{
			name: "invalid id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) + testutils.AddNewRow("-3"),
			},
			args: args{
				recordType: models.TypeText,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: "),
				records:  []models.Record{},
				err:      assert.Error,
			},
		},
		{
			name: "base any",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll),
			},
			args: args{
				recordType: models.TypeCredentials,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): "),
				records: []models.Record{
					{ID: -2},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base by filter",
			input: input{
				command: testutils.AddNewRow(utils.CommandFilters) +
					testutils.AddNewRow(metaKeyVal) +
					testutils.AddNewRow(utils.CommandContinue),
			},
			args: args{
				recordType: models.TypeCredentials,
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': entered filters: map[key:val]\n"),
				records: []models.Record{
					{ID: -2},
				},
				err: assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			records, err := c.handleGet(tt.args.recordType, inMemoryStorage)
			if !tt.want.err(t, err, fmt.Sprintf("handleGet(%v, %v)", tt.args.recordType, inMemoryStorage)) {
				return
			}

			require.Equal(t, len(tt.want.records), len(records), "records count is not equal to expected")
			for idx := range tt.want.records {
				assert.Equal(t, tt.want.records[idx].ID, records[idx].ID)
			}

			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_writeGetResult(t *testing.T) {
	type args struct {
		records []models.Record
	}
	type want struct {
		response []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base empty",
			args: args{
				records: []models.Record{},
			},
			want: want{
				response: []byte("missing record(s)\n"),
			},
		},
		{
			name: "base nil",
			args: args{
				records: nil,
			},
			want: want{
				response: []byte("missing record(s)\n"),
			},
		},
		{
			name: "1 record",
			args: args{
				records: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("found '1' records:\n-----\n   0.  ID: -1  Text: {Data:some text}  MetaData: map[key:key]\n-----\n"),
			},
		},
		{
			name: "2 records",
			args: args{
				records: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("found '2' records:\n-----\n   0.  ID: -1  Text: {Data:some text}                       MetaData: map[key:key]\n   1.  ID: -2  Credential: {Login:login Password:password}  MetaData: map[key:val]\n-----\n"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, _, writer := getCommands("")
			c.writeGetResult(tt.args.records)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_Get(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		parts         []string
		recordsInBase []models.Record
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
			name: "base by id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) + testutils.AddNewRow("-1"),
			},
			args: args{
				parts: []string{utils.CommandGet, models.StrText},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte(`enter search method('id' or 'filters' or 'all'): enter record id: found '1' records:
-----
   0.  ID: -1  Text: {Data:some text}  MetaData: map[key:key]
-----` + "\n"),
			},
		},
		{
			name: "base all",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll),
			},
			args: args{
				parts: []string{utils.CommandGet, models.StrAny},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte(`enter search method('id' or 'filters' or 'all'): found '2' records:
-----
   0.  ID: -1  Text: {Data:some text}                       MetaData: map[key:key]
   1.  ID: -2  Credential: {Login:login Password:password}  MetaData: map[key:val]
-----` + "\n"),
			},
		},
		{
			name: "incorrect command elems count",
			args: args{
				parts: []string{utils.CommandGet},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'get' and object type([any creds card text bin])\n"),
			},
		},
		{
			name: "incorrect type",
			args: args{
				parts: []string{utils.CommandGet, "some incorrect type"},
				recordsInBase: []models.Record{
					{
						ID: -1,
						Data: models.Data{
							RecordType: models.TypeText,
							Text: &models.Text{
								Data: textValue,
							},
							MetaData: map[string]string{
								metaKey: metaKey,
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
								metaKey: metaVal,
							},
						},
					},
				},
			},
			want: want{
				response: []byte("request processing error: process 'get' command: incorrect record type. only ([any creds card text bin]) are supported\n"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			for _, rec := range tt.args.recordsInBase {
				assert.NoError(t, inMemoryStorage.AddRecord(&rec))
			}

			c.Get(tt.args.parts, inMemoryStorage)

			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
