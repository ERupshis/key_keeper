package text

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	someText = "some text"
	empty    = ""
	cancel   = "cancel"
)

func TestText_stateInitial(t *testing.T) {
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
				response: []byte("enter text to save: "),
				state:    addDataState,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			text := &Text{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			assert.Equalf(t, tt.want.state, text.stateInitial(), "stateInitial()")
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestText_stateData(t *testing.T) {
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
				rd: bytes.NewReader([]byte(testutils.AddNewRow(someText))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("entered credential models: {Data:some text}\n"),
				record:   &models.Record{Data: models.Data{Text: &models.Text{Data: someText}}},
				state:    addFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(empty))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("entered credential models: {Data:}\n"),
				record:   &models.Record{Data: models.Data{Text: &models.Text{Data: empty}}},
				state:    addFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Text: &models.Text{}}},
				state:    addDataState,
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
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: nil,
				record:   &models.Record{Data: models.Data{Text: &models.Text{}}},
				state:    addDataState,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			text := &Text{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			got, err := text.stateData(tt.args.record)
			if !tt.want.err(t, err, fmt.Sprintf("stateData(%v)", tt.args.record)) {
				return
			}
			assert.Equalf(t, tt.want.state, got, "stateData(%v)", tt.args.record)
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}

func TestText_addMainData(t *testing.T) {
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
	}{
		{
			name: "base",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(someText))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("enter text to save: entered credential models: {Data:some text}\n"),
				record:   &models.Record{Data: models.Data{Text: &models.Text{Data: someText}}},
				err:      assert.NoError,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(empty))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("enter text to save: entered credential models: {Data:}\n"),
				record:   &models.Record{Data: models.Data{Text: &models.Text{Data: empty}}},
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd: bytes.NewReader([]byte(testutils.AddNewRow(cancel))),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("enter text to save: "),
				record:   &models.Record{Data: models.Data{Text: &models.Text{}}},
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			fields: fields{
				rd: bytes.NewReader([]byte("")),
				wr: bytes.NewBuffer(nil),
				sm: nil,
			},
			args: args{
				record: &models.Record{Data: models.Data{Text: &models.Text{}}},
			},
			want: want{
				response: []byte("enter text to save: "),
				record:   &models.Record{Data: models.Data{Text: &models.Text{}}},
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			iactr := testutils.CreateUserInteractor(tt.fields.rd, tt.fields.wr, logger.CreateMock())
			text := &Text{
				iactr: iactr,
				sm:    tt.fields.sm,
			}
			tt.want.err(t, text.addMainData(tt.args.record), fmt.Sprintf("addMainData(%v)", tt.args.record))
			assert.True(t, reflect.DeepEqual(tt.want.record, tt.args.record))
			assert.Equal(t, tt.want.response, tt.fields.wr.Bytes())
		})
	}
}
