package commands

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/binary"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/credential"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/text"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	cryptorKey = "some key"
	hashKey    = "some key"

	cardNumber     = "1234 1234 1234 1234"
	cardExpiration = "12/23"
	cardCVV        = "123"
	cardHolder     = "holder"

	credLogin    = "login"
	credPassword = "password"

	textValue = "some text"
)

func getCommands(input string) (*Commands, *inmemory.Storage, *bytes.Buffer) {
	reader := bytes.NewReader([]byte(input))
	writer := bytes.NewBuffer(nil)
	userInteractor := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

	sm := statemachines.NewStateMachines(userInteractor)
	bankCard := bankcard.NewBankCard(userInteractor, sm)
	cred := credential.NewCredentials(userInteractor, sm)
	txt := text.NewText(userInteractor, sm)

	hash := hasher.CreateHasher(hashKey, hasher.TypeSHA256, logger.CreateMock())
	dataCryptor := ska.NewSKA(cryptorKey, ska.Key16)

	binaryConfig := binary.Config{
		Iactr:   userInteractor,
		Sm:      sm,
		Hash:    hash,
		Cryptor: dataCryptor,
	}
	bin := binary.NewBinary(&binaryConfig)

	inMemoryStorage := inmemory.NewStorage(dataCryptor)

	c := &Commands{
		iactr:  userInteractor,
		sm:     sm,
		bc:     bankCard,
		creds:  cred,
		text:   txt,
		binary: bin,
	}
	return c, inMemoryStorage, writer
}

func TestCommands_handleAdd(t *testing.T) {

	type input struct {
		command string
	}
	type args struct {
		recordType models.RecordType
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
			name: "base bankcard",
			input: input{
				testutils.AddNewRow(cardNumber) +
					testutils.AddNewRow(cardExpiration) +
					testutils.AddNewRow(cardCVV) +
					testutils.AddNewRow(cardHolder) +
					testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				recordType: models.TypeBankCard,
			},
			want: want{
				response: []byte("enter card number(XXXX XXXX XXXX XXXX): enter card expiration (XX/XX): enter card CVV (XXX or XXXX): enter card holder name: entered card models: {Number:1234 1234 1234 1234 Expiration:12/23 CVV:123 Name:holder}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\n"),
				record:   &models.Record{ID: -1, Data: models.Data{RecordType: models.TypeBankCard, BankCard: &models.BankCard{Number: cardNumber, Expiration: cardExpiration, CVV: cardCVV, Name: cardHolder}}},
				err:      assert.NoError,
			},
		},
		{
			name: "err bankcard",
			input: input{
				testutils.AddNewRow(cardNumber) +
					testutils.AddNewRow(cardExpiration) +
					testutils.AddNewRow(cardCVV) +
					testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordType: models.TypeBankCard,
			},
			want: want{
				response: []byte("enter card number(XXXX XXXX XXXX XXXX): enter card expiration (XX/XX): enter card CVV (XXX or XXXX): enter card holder name: "),
				err:      assert.Error,
			},
		},
		{
			name: "base creds",
			input: input{
				testutils.AddNewRow(credLogin) +
					testutils.AddNewRow(credPassword) +
					testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				recordType: models.TypeCredentials,
			},
			want: want{
				response: []byte("enter credential login: enter credential password: entered credential models: {Login:login Password:password}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\n"),
				record:   &models.Record{ID: -1, Data: models.Data{RecordType: models.TypeCredentials, Credentials: &models.Credential{Login: credLogin, Password: credPassword}}},
				err:      assert.NoError,
			},
		},
		{
			name: "error creds",
			input: input{
				testutils.AddNewRow(credLogin) +
					testutils.AddNewRow(credPassword) +
					utils.CommandSave,
			},
			args: args{
				recordType: models.TypeCredentials,
			},
			want: want{
				response: []byte("enter credential login: enter credential password: entered credential models: {Login:login Password:password}\nenter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "base text",
			input: input{
				testutils.AddNewRow(textValue) +
					testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				recordType: models.TypeText,
			},
			want: want{
				response: []byte("enter text to save: entered text models: {Data:some text}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\n"),
				record:   &models.Record{ID: -1, Data: models.Data{RecordType: models.TypeText, Text: &models.Text{Data: textValue}}},
				err:      assert.NoError,
			},
		},
		{
			name: "error text",
			input: input{
				testutils.AddNewRow(textValue) +
					testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordType: models.TypeText,
			},
			want: want{
				response: []byte("enter text to save: entered text models: {Data:some text}\nenter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "invalid type",
			input: input{
				command: "",
			},
			args: args{
				recordType: models.TypeUndefined,
			},
			want: want{
				response: nil,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			record, err := c.handleAdd(tt.args.recordType, inMemoryStorage)
			tt.want.err(t, err, "err assertion fail")

			if record != nil {
				assert.True(t, reflect.DeepEqual(record.Data, tt.want.record.Data), "record fail")
			}

			storageRecord, err := inMemoryStorage.GetRecord(-1)
			if err == nil {
				assert.True(t, reflect.DeepEqual(record.Data, storageRecord.Data), "record in storage fail")
			}

			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_handleCommandError(t *testing.T) {
	type args struct {
		err            error
		command        string
		supportedTypes []string
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
			name: "base",
			args: args{
				err:            nil,
				command:        utils.CommandAdd,
				supportedTypes: []string{models.StrBankCard, models.StrBinary},
			},
			want: want{
				response: []byte("request processing error: <nil>\n"),
			},
		},
		{
			name: "cancel by user",
			args: args{
				err:            errs.ErrInterruptedByUser,
				command:        utils.CommandAdd,
				supportedTypes: []string{models.StrBankCard, models.StrBinary},
			},
			want: want{
				response: []byte("'add' command was canceled by user\n"),
			},
		},
		{
			name: "incorrect record type",
			args: args{
				err:            errs.ErrIncorrectRecordType,
				command:        utils.CommandDelete,
				supportedTypes: []string{models.StrBankCard, models.StrBinary},
			},
			want: want{
				response: []byte("request processing error: incorrect record type. only ([card bin]) are supported\n"),
			},
		},
		{
			name: "incorrect server action type",
			args: args{
				err:            errs.ErrIncorrectServerActionType,
				command:        utils.CommandDelete,
				supportedTypes: []string{models.StrBinary},
			},
			want: want{
				response: []byte("request processing error: incorrect server action type. only ([bin]) are supported\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(nil)
			writer := bytes.NewBuffer(nil)
			userInteractor := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			c := &Commands{
				iactr: userInteractor,
			}

			c.handleCommandError(tt.args.err, tt.args.command, tt.args.supportedTypes)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_Add(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		parts []string
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
			name: "base text",
			input: input{
				command: testutils.AddNewRow(textValue) + testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				parts: []string{"add", models.StrText},
			},
			want: want{
				response: []byte("enter text to save: entered text models: {Data:some text}\nenter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\nrecord added: {ID: -1, Text: {Data:some text}, MetaData: map[]}\n"),
			},
		},
		{
			name: "error text",
			input: input{
				command: testutils.AddNewRow(textValue) + testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				parts: []string{"add", models.StrText},
			},
			want: want{
				response: []byte("enter text to save: entered text models: {Data:some text}\nenter meta models(format: 'key : value') or 'cancel' or 'save': 'add' command was canceled by user\n"),
			},
		},
		{
			name: "invalid parts count",
			input: input{
				command: testutils.AddNewRow(textValue) + testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				parts: []string{"add", models.StrText, models.StrBinary},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'add' and object type([creds card text bin])\n"),
			},
		},
		{
			name: "invalid second part",
			input: input{
				command: testutils.AddNewRow(textValue) + testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				parts: []string{"add", models.StrUndefined},
			},
			want: want{
				response: []byte("request processing error: process 'add' command: incorrect record type. only ([creds card text bin]) are supported\n"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			c.Add(tt.args.parts, inMemoryStorage)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
