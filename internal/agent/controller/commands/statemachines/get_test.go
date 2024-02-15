package statemachines

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestStateMachines_stateGetFiltersValue(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		filters map[string]string
	}
	type want struct {
		filters  map[string]string
		response []byte
		state    stateGetFilters
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
				filters: map[string]string{},
			},
			want: want{
				filters:  map[string]string{metaKey: metaVal},
				response: nil,
				state:    getFiltersInitialState,
				err:      assert.NoError,
			},
		},
		{
			name: "continue",
			input: input{
				command: testutils.AddNewRow(utils.CommandContinue),
			},
			args: args{
				filters: map[string]string{metaKey: metaVal},
			},
			want: want{
				filters:  map[string]string{metaKey: metaVal},
				response: []byte(fmt.Sprintf("entered filters: map[%s:%s]\n", metaKey, metaVal)),
				state:    getFiltersFinishState,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid input",
			input: input{
				command: testutils.AddNewRow(metaKey),
			},
			args: args{
				filters: map[string]string{},
			},
			want: want{
				filters:  map[string]string{},
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    getFiltersValueState,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				filters: map[string]string{},
			},
			want: want{
				filters:  map[string]string{},
				response: nil,
				state:    getFiltersValueState,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandCancel,
			},
			args: args{
				filters: map[string]string{},
			},
			want: want{
				filters:  map[string]string{},
				response: nil,
				state:    getFiltersValueState,
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

			state, err := s.stateGetFiltersValue(tt.args.filters)
			if !tt.want.err(t, err, fmt.Sprintf("stateGetFiltersValue(%v)", tt.args.filters)) {
				return
			}
			assert.Equalf(t, tt.want.state, state, "stateGetFiltersValue(%v)", tt.args.filters)
			assert.True(t, reflect.DeepEqual(tt.want.filters, tt.args.filters), "filters fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateGetFiltersInitial(t *testing.T) {
	type want struct {
		response []byte
		state    stateGetFilters
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				response: []byte("enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': "),
				state:    getFiltersValueState,
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

			assert.Equalf(t, tt.want.state, s.stateGetFiltersInitial(), "stateGetFiltersInitial()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_getFilters(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		filters  map[string]string
		response []byte
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
				command: testutils.AddNewRow(metaKeyVal) + testutils.AddNewRow(utils.CommandContinue),
			},
			want: want{
				filters:  map[string]string{metaKey: metaVal},
				response: []byte("enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': entered filters: map[key:val]\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				filters:  nil,
				response: []byte("enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': "),
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

			filters, err := s.getFilters()
			tt.want.err(t, err, "getFilters()")
			assert.Equalf(t, tt.want.filters, filters, "getFilters()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateGetIDValue(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		state    stateGetID
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
				command: testutils.AddNewRow(positive),
			},
			want: want{
				response: nil,
				state:    getIDFinishState,
				id:       positiveInt,
				err:      assert.NoError,
			},
		},
		{
			name: "base negative",
			input: input{
				command: testutils.AddNewRow(negative),
			},
			want: want{
				response: nil,
				state:    getIDFinishState,
				id:       negativeInt,
				err:      assert.NoError,
			},
		},
		{
			name: "str",
			input: input{
				command: testutils.AddNewRow(str),
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    getIDValueState,
				id:       0,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: nil,
				state:    getIDValueState,
				id:       0,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandCancel,
			},
			want: want{
				response: nil,
				state:    getIDValueState,
				id:       0,
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				command: testutils.AddNewRow(empty),
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    getIDValueState,
				id:       0,
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

			state, id, err := s.stateGetIDValue()
			tt.want.err(t, err, "stateGetIDValue()")
			assert.Equalf(t, tt.want.state, state, "state fail")
			assert.Equalf(t, tt.want.id, id, "id fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateGetIDInitial(t *testing.T) {
	type want struct {
		response []byte
		state    stateGetID
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				state:    getIDValueState,
				response: []byte("enter record id: "),
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

			assert.Equalf(t, tt.want.state, s.stateGetIDInitial(), "stateGetIDInitial()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_getID(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
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
				command: testutils.AddNewRow(positive),
			},
			want: want{
				response: []byte("enter record id: "),
				id:       positiveInt,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter record id: "),
				id:       0,
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

			id, err := s.getID()
			tt.want.err(t, err, "getID()")
			assert.Equalf(t, tt.want.id, id, "id fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateGetMethodData(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		state    stateGetMethod
		method   string
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "base id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID),
			},
			want: want{
				response: nil,
				state:    getMethodFinishState,
				method:   utils.CommandID,
				err:      assert.NoError,
			},
		},
		{
			name: "base filters",
			input: input{
				command: testutils.AddNewRow(utils.CommandFilters),
			},
			want: want{
				response: nil,
				state:    getMethodFinishState,
				method:   utils.CommandFilters,
				err:      assert.NoError,
			},
		},
		{
			name: "base all",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll),
			},
			want: want{
				response: nil,
				state:    getMethodFinishState,
				method:   utils.CommandAll,
				err:      assert.NoError,
			},
		},
		{
			name: "base invalid",
			input: input{
				command: testutils.AddNewRow("invalid"),
			},
			want: want{
				response: []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:    getMethodSelectionState,
				method:   empty,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: nil,
				state:    getMethodSelectionState,
				method:   empty,
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandCancel,
			},
			want: want{
				response: nil,
				state:    getMethodSelectionState,
				method:   empty,
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

			state, method, err := s.stateGetMethodData()
			tt.want.err(t, err, "stateGetMethodData()")
			assert.Equalf(t, tt.want.state, state, "stateGetMethodData()")
			assert.Equalf(t, tt.want.method, method, "stateGetMethodData()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateGetMethodInitial(t *testing.T) {
	type want struct {
		response []byte
		state    stateGetMethod
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): "),
				state:    getMethodSelectionState,
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

			assert.Equalf(t, tt.want.state, s.stateGetMethodInitial(), "stateGetMethodInitial()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_getMethod(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		method   string
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
				command: testutils.AddNewRow(utils.CommandID),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): "),
				method:   utils.CommandID,
				err:      assert.NoError,
			},
		},
		{
			name: "invalid method",
			input: input{
				command: testutils.AddNewRow(utils.CommandDelete),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): incorrect input, try again or interrupt by 'cancel' command: "),
				method:   "",
				err:      assert.Error,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): "),
				method:   "",
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

			method, err := s.getMethod()
			tt.want.err(t, err, "getMethod()")
			assert.Equalf(t, tt.want.method, method, "method fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_getStateAccordingMethod(t *testing.T) {
	type args struct {
		method string
	}
	type want struct {
		state stateGet
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base id",
			args: args{
				method: utils.CommandID,
			},
			want: want{
				state: getSearchByID,
			},
		},
		{
			name: "base filters",
			args: args{
				method: utils.CommandFilters,
			},
			want: want{
				state: getSearchByFilters,
			},
		},
		{
			name: "base all",
			args: args{
				method: utils.CommandAll,
			},
			want: want{
				state: getSearchAllByType,
			},
		},
		{
			name: "base invalid",
			args: args{
				method: utils.CommandDelete,
			},
			want: want{
				state: getInitialState,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader(nil)
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}
			assert.Equalf(t, tt.want.state, s.getStateAccordingMethod(tt.args.method), "getStateAccordingMethod(%v)", tt.args.method)
		})
	}
}

func TestStateMachines_Get(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		id       int64
		filters  map[string]string
		err      assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "base all",
			input: input{
				command: testutils.AddNewRow(utils.CommandAll),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): "),
				id:       9,
				filters:  map[string]string{},
				err:      assert.NoError,
			},
		},
		{
			name: "base id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) + testutils.AddNewRow(positive),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: "),
				id:       9,
				filters:  nil,
				err:      assert.NoError,
			},
		},
		{
			name: "cancel id",
			input: input{
				command: testutils.AddNewRow(utils.CommandID) + testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter record id: "),
				err:      assert.Error,
			},
		},
		{
			name: "base filters",
			input: input{
				command: testutils.AddNewRow(utils.CommandFilters) + testutils.AddNewRow(metaKeyVal) + testutils.AddNewRow(utils.CommandContinue),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': entered filters: map[key:val]\n"),
				id:       9,
				filters:  map[string]string{metaKey: metaVal},
				err:      assert.NoError,
			},
		},
		{
			name: "cancel filters",
			input: input{
				command: testutils.AddNewRow(utils.CommandFilters) + testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): enter filters through meta models(format: 'key : value') or 'cancel' or 'continue': "),
				id:       9,
				filters:  nil,
				err:      assert.Error,
			},
		},
		{
			name: "invalid",
			input: input{
				command: testutils.AddNewRow(utils.CommandDelete),
			},
			want: want{
				response: []byte("enter search method('id' or 'filters' or 'all'): incorrect input, try again or interrupt by 'cancel' command: "),
				id:       9,
				filters:  nil,
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

			id, filters, err := s.Get()
			tt.want.err(t, err, "Get()")

			if id != nil {
				assert.Equalf(t, tt.want.id, *id, "id fail")
			}

			assert.True(t, reflect.DeepEqual(tt.want.filters, filters), "filters fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
