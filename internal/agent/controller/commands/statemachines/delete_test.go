package statemachines

import (
	"bytes"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	positiveInt = 9
	negativeInt = -9

	positive = "9"
	negative = "-9"
	str      = "text"
	empty    = ""
)

func TestStateMachines_stateDeleteIDValue(t *testing.T) {
	type input struct {
		data string
	}
	type want struct {
		response []byte
		state    stateDelete
		id       int64
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "base",
			input: input{
				data: testutils.AddNewRow(positive),
			},
			want: want{
				response: nil,
				state:    deleteFinishState,
				id:       positiveInt,
				err:      assert.NoError,
			},
		},
		{
			name: "base negative",
			input: input{
				data: testutils.AddNewRow(negative),
			},
			want: want{
				response: nil,
				state:    deleteFinishState,
				id:       negativeInt,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid, str",
			input: input{
				data: testutils.AddNewRow(str),
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    deleteIDState,
				id:       0,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				data: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: nil,
				state:    deleteIDState,
				id:       0,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				data: utils.CommandCancel,
			},
			want: want{
				response: nil,
				state:    deleteIDState,
				id:       0,
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				data: testutils.AddNewRow(empty),
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    deleteIDState,
				id:       0,
				err:      assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte(tt.input.data))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}

			state, id, err := s.stateDeleteIDValue()
			tt.want.err(t, err, "stateDeleteIDValue()")
			assert.Equalf(t, tt.want.state, state, "state fail")
			assert.Equalf(t, tt.want.id, id, "id fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateDeleteIDInitial(t *testing.T) {
	type want struct {
		response []byte
		state    stateDelete
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				response: []byte("enter record id: "),
				state:    deleteIDState,
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

			assert.Equalf(t, tt.want.state, s.stateDeleteIDInitial(), "state fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_Delete(t *testing.T) {
	type input struct {
		data string
	}
	type want struct {
		response []byte
		nil      bool
		id       int64
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "base",
			input: input{
				data: testutils.AddNewRow(positive),
			},
			want: want{
				response: []byte("enter record id: "),
				nil:      false,
				id:       positiveInt,
				err:      assert.NoError,
			},
		},
		{
			name: "base negative",
			input: input{
				data: testutils.AddNewRow(negative),
			},
			want: want{
				response: []byte("enter record id: "),
				nil:      false,
				id:       negativeInt,
				err:      assert.NoError,
			},
		},
		{
			name: "not numeric",
			input: input{
				data: testutils.AddNewRow(str),
			},
			want: want{
				response: []byte("enter record id: incorrect input, try again or interrupt by 'cancel' command: "),
				nil:      true,
				id:       positiveInt,
				err:      assert.Error,
			},
		},
		{
			name: "cancel",
			input: input{
				data: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter record id: "),
				nil:      true,
				id:       positiveInt,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				data: utils.CommandCancel,
			},
			want: want{
				response: []byte("enter record id: "),
				nil:      true,
				id:       positiveInt,
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				data: testutils.AddNewRow(empty),
			},
			want: want{
				response: []byte("enter record id: incorrect input, try again or interrupt by 'cancel' command: "),
				nil:      true,
				id:       positiveInt,
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte(tt.input.data))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}

			got, err := s.Delete()
			tt.want.err(t, err, "Delete()")
			if tt.want.nil {
				assert.Nil(t, got, "id is not nil, but has to be")

			} else {
				assert.Equalf(t, tt.want.id, *got, "id fail")
			}
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
