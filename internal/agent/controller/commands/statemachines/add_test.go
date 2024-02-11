package statemachines

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	metaKeyVal = "key : val"
	metaKey    = "key"
	metaVal    = "val"
)

func TestStateMachines_stateMetaData(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
		state    stateAddMeta
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
				command: testutils.AddNewRow(metaKeyVal),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: nil,
				state:    addMetaInitialState,
				err:      assert.NoError,
			},
		},
		{
			name: "save",
			input: input{
				command: testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("entered metadata: map[]\n"),
				state:    addMetaFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "save with pair",
			input: input{
				command: testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: []byte("entered metadata: map[key:val]\n"),
				state:    addMetaFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid command with pair",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    addMetaDataState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid key/val input",
			input: input{
				command: testutils.AddNewRow(metaKey + ":" + metaVal),
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    addMetaDataState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: nil,
				state:    addMetaDataState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandCancel,
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: nil,
				state:    addMetaDataState,
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				command: testutils.AddNewRow(""),
			},
			args: args{
				record: &models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    addMetaDataState,
				err:      assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte(tt.input.command))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}
			got, err := s.stateMetaData(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateMetaData(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateMetaData(%v)", tt.args.record)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateMetaInitial(t *testing.T) {
	type want struct {
		response []byte
		state    stateAddMeta
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': "),
				state:    addMetaDataState,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(nil)
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}
			assert.Equalf(t, tt.want.state, s.stateMetaInitial(), "stateMetaInitial()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_addMetaData(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		record *models.Record
	}
	type want struct {
		response []byte
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
				command: testutils.AddNewRow(metaKeyVal) + testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[key:val]\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "base save empty",
			input: input{
				command: testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandCancel,
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				command: testutils.AddNewRow(""),
			},
			args: args{
				record: &models.Record{},
			},
			want: want{
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': incorrect input, try again or interrupt by 'cancel' command: "),
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte(tt.input.command))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}

			tt.want.err(t, s.addMetaData(tt.args.record), fmt.Sprintf("addMetaData(%v)", tt.args.record))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_Add(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		cfg AddConfig
	}
	type want struct {
		record   models.Record
		response []byte
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
				command: testutils.AddNewRow(metaKeyVal) + testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				cfg: AddConfig{
					Record: &models.Record{},
					MainData: func(record *models.Record) error {
						return nil
					},
				},
			},
			want: want{
				record:   models.Record{Data: models.Data{MetaData: map[string]string{metaKey: metaVal}}},
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[key:val]\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "base without meta",
			input: input{
				command: testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				cfg: AddConfig{
					Record: &models.Record{},
					MainData: func(record *models.Record) error {
						return nil
					},
				},
			},
			want: want{
				record:   models.Record{},
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': entered metadata: map[]\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "cancel meta",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				cfg: AddConfig{
					Record: &models.Record{},
					MainData: func(record *models.Record) error {
						return nil
					},
				},
			},
			want: want{
				record:   models.Record{},
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "eof meta",
			input: input{
				command: "",
			},
			args: args{
				cfg: AddConfig{
					Record: &models.Record{},
					MainData: func(record *models.Record) error {
						return nil
					},
				},
			},
			want: want{
				record:   models.Record{},
				response: []byte("enter meta models(format: 'key : value') or 'cancel' or 'save': "),
				err:      assert.Error,
			},
		},
		{
			name: "main interrupted",
			input: input{
				command: testutils.AddNewRow(utils.CommandSave),
			},
			args: args{
				cfg: AddConfig{
					Record: &models.Record{},
					MainData: func(record *models.Record) error {
						return errs.ErrInterruptedByUser
					},
				},
			},
			want: want{
				record:   models.Record{},
				response: nil,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte(tt.input.command))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}

			tt.want.err(t, s.Add(tt.args.cfg), fmt.Sprintf("Add(%v)", tt.args.cfg))
			assert.True(t, reflect.DeepEqual(tt.want.record, *tt.args.cfg.Record), "record fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
