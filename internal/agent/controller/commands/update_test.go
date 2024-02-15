package commands

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/stretchr/testify/assert"
)

const (
	textValueUpdated    = "some updated text"
	credLoginUpdated    = "login new"
	credPasswordUpdated = "password new"

	cardNumberUpdated     = "8888 8888 8888 8888"
	cardExpirationUpdated = "12/88"
	cardCVVUpdated        = "8888"
	cardHolderUpdated     = "new holder"

	metaValUpdated = "value new"
)

func TestCommands_confirmAndUpdateRecordByID(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		recordInBase   *models.Record
		recordToUpdate *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		args  args
		want  want
	}{
		{
			name: "base confirmed",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				recordToUpdate: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
					},
				},
			},
			want: want{
				response: []byte("Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record successfully updated\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base not confirmed",
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				recordToUpdate: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
					},
				},
			},
			want: want{
				response: []byte("Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record updating was interrupted by user\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "no changes in update",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				recordToUpdate: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("Do you really want to update the record '{ID: -1, Text: {Data:some text}, MetaData: map[]}'(yes/no): Record successfully updated\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				recordToUpdate: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
					},
				},
			},
			want: want{
				response: []byte("Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): "),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
		{
			name: "invalid id",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				recordToUpdate: &models.Record{
					ID: -2,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
					},
				},
			},
			want: want{
				response: []byte("Do you really want to update the record '{ID: -2, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): "),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase), "add record in storage")
			tt.want.err(t, c.confirmAndUpdateRecordByID(tt.args.recordToUpdate, inMemoryStorage), fmt.Sprintf("confirmAndUpdateRecordByID(%v, %v)", tt.args.recordToUpdate, inMemoryStorage))

			updatedRecord, err := inMemoryStorage.GetRecord(tt.args.recordInBase.ID)
			assert.NoError(t, err, "get record from storage")
			assert.Equal(t, tt.want.record.Data, updatedRecord.Data, "update effect check")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_findAndUpdateRecordByID(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		id           int64
		recordInBase *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		args  args
		want  want
	}{
		{
			name: "base Text",
			input: input{
				command: testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("enter text to save: entered text models: {Data:some updated text}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\nDo you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record successfully updated\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "record not found",
			input: input{
				command: testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -2,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: nil,
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
		{
			name: "base creds",
			input: input{
				command: testutils.AddNewRow(credLoginUpdated) +
					testutils.AddNewRow(credPasswordUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeCredentials,
						Credentials: &models.Credential{
							Login:    credLogin,
							Password: credPassword,
						},
					},
				},
			},
			want: want{
				response: []byte("enter credential login: enter credential password: entered credential models: {Login:login new Password:password new}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\nDo you really want to update the record '{ID: -1, Credential: {Login:login new Password:password new}, MetaData: map[]}'(yes/no): Record successfully updated\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeCredentials,
						Credentials: &models.Credential{
							Login:    credLoginUpdated,
							Password: credPasswordUpdated,
						},
						MetaData: map[string]string{},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base bank card",
			input: input{
				command: testutils.AddNewRow(cardNumberUpdated) +
					testutils.AddNewRow(cardExpirationUpdated) +
					testutils.AddNewRow(cardCVVUpdated) +
					testutils.AddNewRow(cardHolderUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeBankCard,
						BankCard: &models.BankCard{
							Number:     cardNumber,
							Expiration: cardExpiration,
							CVV:        cardCVV,
							Name:       cardHolder,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter card number(1234 1234 1234 1234): enter card expiration (12/23): enter card CVV (123): enter card holder name(holder): entered card models: {Number:8888 8888 8888 8888 Expiration:12/88 CVV:8888 Name:new holder}
enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]
Do you really want to update the record '{ID: -1, BankCard: {Number:8888 8888 8888 8888 Expiration:12/88 CVV:8888 Name:new holder}, MetaData: map[]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeBankCard,
						BankCard: &models.BankCard{
							Number:     cardNumberUpdated,
							Expiration: cardExpirationUpdated,
							CVV:        cardCVVUpdated,
							Name:       cardHolderUpdated,
						},
						MetaData: map[string]string{},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base Text + new meta",
			input: input{
				command: testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(metaKeyVal) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[key:val]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[key:val]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{metaKey: metaVal},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base Text + update meta",
			input: input{
				command: testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(fmt.Sprintf("%s : %s", metaKey, metaValUpdated)) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
						MetaData: map[string]string{metaKey: metaVal},
					},
				},
			},
			want: want{
				response: []byte(`enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[key:value new]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[key:value new]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{metaKey: metaValUpdated},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "base Text + partial update meta",
			input: input{
				command: testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(fmt.Sprintf("%s : %s", metaVal, metaValUpdated)) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
						MetaData: map[string]string{
							metaKey: metaVal,
							metaVal: metaKey,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[val:value new]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[val:value new]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{
							metaVal: metaValUpdated,
						},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				id: -1,
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("enter text to save: "),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase), "add record in storage")
			tt.want.err(t, c.findAndUpdateRecordByID(tt.args.id, inMemoryStorage), fmt.Sprintf("findAndUpdateRecordByID(%v, %v)", tt.args.id, inMemoryStorage))
			updatedRecord, err := inMemoryStorage.GetRecord(tt.args.recordInBase.ID)
			assert.NoError(t, err, "get record from storage")
			assert.Equal(t, tt.want.record.Data.RecordType, updatedRecord.Data.RecordType, "update effect check Type")
			assert.Equal(t, tt.want.record.Data.Text, updatedRecord.Data.Text, "update effect check Text")
			assert.Equal(t, tt.want.record.Data.Credentials, updatedRecord.Data.Credentials, "update effect check Creds")
			assert.Equal(t, tt.want.record.Data.BankCard, updatedRecord.Data.BankCard, "update effect check Bank card")
			assert.True(t, reflect.DeepEqual(tt.want.record.Data.MetaData, updatedRecord.Data.MetaData), "update effect check meta")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_handleUpdate(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		recordInBase *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
		err      assert.ErrorAssertionFunc
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
				command: testutils.AddNewRow("-1") +
					testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter record id: enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "enter id from 2nd time",
			input: input{
				command: testutils.AddNewRow("should be id") +
					testutils.AddNewRow("-1") +
					testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter record id: incorrect input, try again or interrupt by 'cancel' command: enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{},
					},
				},
				err: assert.NoError,
			},
		},
		{
			name: "cancel on id",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("enter record id: "),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
		{
			name: "invalid id",
			input: input{
				command: testutils.AddNewRow("-2"),
			},
			args: args{
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("enter record id: "),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
				err: assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase), "add record in storage")
			tt.want.err(t, c.handleUpdate(inMemoryStorage), fmt.Sprintf("handleUpdate(%v)", inMemoryStorage))
			updatedRecord, err := inMemoryStorage.GetRecord(tt.args.recordInBase.ID)
			assert.NoError(t, err, "get record from storage")
			assert.Equal(t, tt.want.record.Data.RecordType, updatedRecord.Data.RecordType, "update effect check Type")
			assert.Equal(t, tt.want.record.Data.Text, updatedRecord.Data.Text, "update effect check Text")
			assert.Equal(t, tt.want.record.Data.Credentials, updatedRecord.Data.Credentials, "update effect check Creds")
			assert.Equal(t, tt.want.record.Data.BankCard, updatedRecord.Data.BankCard, "update effect check Bank card")
			assert.True(t, reflect.DeepEqual(tt.want.record.Data.MetaData, updatedRecord.Data.MetaData), "update effect check meta")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_Update(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		parts        []string
		recordInBase *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
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
				command: testutils.AddNewRow("-1") +
					testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				parts: []string{utils.CommandUpdate},
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte(`enter record id: enter text to save: entered text models: {Data:some updated text}
enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]
Do you really want to update the record '{ID: -1, Text: {Data:some updated text}, MetaData: map[]}'(yes/no): Record successfully updated
`),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValueUpdated,
						},
						MetaData: map[string]string{},
					},
				},
			},
		},
		{
			name: "invalid command parse",
			input: input{
				command: testutils.AddNewRow("-1") +
					testutils.AddNewRow(textValueUpdated) +
					testutils.AddNewRow(utils.CommandSave) +
					testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				parts: []string{utils.CommandUpdate, models.StrText},
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'update' only\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				parts: []string{utils.CommandUpdate},
				recordInBase: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
			want: want{
				response: []byte("enter record id: 'update' command was canceled by user\n"),
				record: &models.Record{
					ID: -1,
					Data: models.Data{
						RecordType: models.TypeText,
						Text: &models.Text{
							Data: textValue,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase), "add record in storage")
			c.Update(tt.args.parts, inMemoryStorage)
			updatedRecord, err := inMemoryStorage.GetRecord(tt.args.recordInBase.ID)
			assert.NoError(t, err, "get record from storage")
			assert.Equal(t, tt.want.record.Data.RecordType, updatedRecord.Data.RecordType, "update effect check Type")
			assert.Equal(t, tt.want.record.Data.Text, updatedRecord.Data.Text, "update effect check Text")
			assert.Equal(t, tt.want.record.Data.Credentials, updatedRecord.Data.Credentials, "update effect check Creds")
			assert.Equal(t, tt.want.record.Data.BankCard, updatedRecord.Data.BankCard, "update effect check Bank card")
			assert.True(t, reflect.DeepEqual(tt.want.record.Data.MetaData, updatedRecord.Data.MetaData), "update effect check meta")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
