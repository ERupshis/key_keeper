package bankcard

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	numberCorrect = "1234 1234 1234 1234"
	numberWrong1  = "asdf 1234 2345 1111"
	numberWrong2  = "1111123423451111"

	expirationCorrect = "12/12"
	expirationWrong1  = "13/12"
	expirationWrong2  = "asdf"
	expirationWrong3  = "1212"

	cvvCorrect  = "123"
	cvvCorrect2 = "1234"
	cvvWrong1   = "asd"
	cvvWrong2   = "asdf"
	cvvWrong3   = "12"
	cvvWrong4   = "12345"

	cardHolderCorrect = "Card Holder"
	customCardHolder  = "Custom Name"

	tmplNumber     = "XXXX XXXX XXXX XXXX"
	tmplExpiration = "XX/XX"
	tmplCVV        = "XXX or XXXX"
)

func TestCredential_stateInitial(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		state    addState
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
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer([]byte{}),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: &models.BankCard{}}},
			},
			want: want{
				response: []byte("enter card number(): "),
				state:    addNumberState,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			card := &BankCard{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			assert.Equalf(t, tt.want.state, card.stateInitial(tt.args.record), "stateInitial()")
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBankCard_stateNumber(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(numberCorrect))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: getBankCardDataTemplate()}},
			},
			want: want{
				response: []byte("enter card expiration (XX/XX): "),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: numberCorrect, Expiration: tmplExpiration, CVV: tmplCVV}}},
				state:    addExpirationState,
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addNumberState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid number",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(numberWrong1))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addNumberState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid number 2",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(numberWrong2))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addNumberState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addNumberState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addNumberState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			card := &BankCard{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := card.stateNumber(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateNumber(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateNumber(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBankCard_stateExpiration(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(expirationCorrect))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: getBankCardDataTemplate()}},
			},
			want: want{
				response: []byte("enter card CVV (XXX or XXXX): "),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: expirationCorrect, CVV: tmplCVV}}},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid expiration",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(expirationWrong1))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid expiration 2",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(expirationWrong2))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid expiration 3",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(expirationWrong3))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addExpirationState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			card := &BankCard{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := card.stateExpiration(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateExpiration(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateExpiration(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBankCard_stateCVV(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvCorrect))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: getBankCardDataTemplate()}},
			},
			want: want{
				response: []byte("enter card holder name: "),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: tmplExpiration, CVV: cvvCorrect}}},
				state:    addCardHolderState,
				err:      assert.NoError,
			},
		},
		{
			name: "base 2",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvCorrect2))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: getBankCardDataTemplate()}},
			},
			want: want{
				response: []byte("enter card holder name: "),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: tmplExpiration, CVV: cvvCorrect2}}},
				state:    addCardHolderState,
				err:      assert.NoError,
			},
		},
		{
			name: "base not empty name",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvCorrect))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: tmplExpiration, CVV: tmplCVV, Name: customCardHolder}}},
			},
			want: want{
				response: []byte("enter card holder name(Custom Name): "),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: tmplExpiration, CVV: cvvCorrect, Name: customCardHolder}}},
				state:    addCardHolderState,
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid cvv",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvWrong1))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid cvv 2",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvWrong2))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid cvv 3",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvWrong3))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid cvv 4",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cvvWrong4))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addCVVState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			card := &BankCard{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := card.stateCVV(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateCVV(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateCVV(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestBankCard_stateCardHolder(t *testing.T) {
	type fields struct {
		rd *bytes.Reader
		wr *bytes.Buffer
		sm *statemachines.StateMachines
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		record   *models.Record
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cardHolderCorrect))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{BankCard: getBankCardDataTemplate()}},
			},
			want: want{
				response: []byte("entered card models: {Number:XXXX XXXX XXXX XXXX Expiration:XX/XX CVV:XXX or XXXX Name:Card Holder}\n"),
				record:   &models.Record{Data: models.Data{BankCard: &models.BankCard{Number: tmplNumber, Expiration: tmplExpiration, CVV: tmplCVV, Name: cardHolderCorrect}}},
				state:    addFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				record:   &models.Record{},
				state:    addCardHolderState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addCardHolderState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				record:   &models.Record{},
				state:    addCardHolderState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			card := &BankCard{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := card.stateCardHolder(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateCardHolder(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateCardHolder(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

// TODO: addMainData.
