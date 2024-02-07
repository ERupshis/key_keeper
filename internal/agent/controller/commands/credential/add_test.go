package credential

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
	login = "login"
	pwd   = "pwd"
)

func TestCredential_stateInitial(t *testing.T) {
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
				response: []byte("enter credential login: "),
				state:    addLoginState,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			cred := &Credential{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			assert.Equalf(t, tt.want.state, cred.stateInitial(), "stateInitial()")
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestCredential_stateLogin(t *testing.T) {
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(login))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential password: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Login: login, Password: ""}}},
				state:    addPasswordState,
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential password: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Login: ""}}},
				state:    addPasswordState,
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				state:    addLoginState,
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				state:    addLoginState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			cred := &Credential{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := cred.stateLogin(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateLogin(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateLogin(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestCredential_statePassword(t *testing.T) {
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(pwd))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("entered credential models: {Login: Password:pwd}\n"),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Password: pwd}}},
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("entered credential models: {Login: Password:}\n"),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				state:    addFinishState,
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				state:    addPasswordState,
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
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				state:    addPasswordState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			cred := &Credential{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := cred.statePassword(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("statePassword(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "statePassword(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestCredential_addMainData(t *testing.T) {
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
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{{
		name: "base",
		fields: fields{
			rd: bytes.NewReader([]byte(testutils.AddNewRow(login) + testutils.AddNewRow(pwd))),
			wr: bytes.NewBuffer(nil),
			sm: nil,
		},
		args: args{
			record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
		},
		want: want{
			response: []byte("enter credential login: enter credential password: entered credential models: {Login:login Password:pwd}\n"),
			record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Login: login, Password: pwd}}},
			err:      assert.NoError,
		},
	},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow("") + testutils.AddNewRow(""))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential login: enter credential password: entered credential models: {Login: Password:}\n"),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				err:      assert.NoError,
			},
		},
		{
			name: "cancel on login",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential login: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				err:      assert.Error,
			},
		},
		{
			name: "cancel on pwd",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(login) + testutils.AddNewRow(utils.CommandCancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential login: enter credential password: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Login: login}}},
				err:      assert.Error,
			},
		},
		{
			name: "eof on login",
			fields: fields{
				rd: bytes.NewReader([]byte{}),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential login: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
				err:      assert.Error,
			},
		},
		{
			name: "eof on pwd",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(login))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Credentials: &models.Credential{}}},
			},
			want: want{
				response: []byte("enter credential login: enter credential password: "),
				record:   &models.Record{Data: models.Data{Credentials: &models.Credential{Login: login}}},
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			cred := &Credential{
				iactr: iactr,
				sm:    tt.fields.sm,
			}

			if !tt.want.err(t, cred.addMainData(tt.args.record), fmt.Sprintf("addMainData(%v)", tt.args.record)) {
				return
			}

			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}
