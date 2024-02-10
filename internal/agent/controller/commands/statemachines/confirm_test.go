package statemachines

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	invalid = "incorrect_command"
)

func TestStateMachines_stateConfirmApprove(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response          []byte
		state             stateConfirm
		confirmSuccessful bool
		err               assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "base",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			want: want{
				response:          nil,
				state:             confirmFinishState,
				confirmSuccessful: true,
				err:               assert.NoError,
			},
		},
		{
			name: "base 2",
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			want: want{
				response:          nil,
				state:             confirmFinishState,
				confirmSuccessful: false,
				err:               assert.NoError,
			},
		},
		{
			name: "invalid command",
			input: input{
				command: testutils.AddNewRow(invalid),
			},
			want: want{
				response:          []byte("incorrect input, try again or interrupt by 'cancel' command: "),
				state:             confirmApproveState,
				confirmSuccessful: false,
				err:               assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response:          nil,
				state:             confirmApproveState,
				confirmSuccessful: false,
				err:               assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: utils.CommandYes,
			},
			want: want{
				response:          nil,
				state:             confirmApproveState,
				confirmSuccessful: false,
				err:               assert.Error,
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
			state, valid, err := s.stateConfirmApprove()
			tt.want.err(t, err)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
			assert.Equal(t, tt.want.state, state, "operation state fail")
			assert.Equal(t, tt.want.confirmSuccessful, valid, "incorrect input")
		})
	}
}

func TestStateMachines_stateConfirmInitial(t *testing.T) {
	type args struct {
		record  *models.Record
		command string
	}
	type want struct {
		response []byte
		state    stateConfirm
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "update",
			args: args{
				record:  &models.Record{},
				command: utils.CommandUpdate,
			},
			want: want{
				response: []byte(fmt.Sprintf("Do you really want to update the record '%s'(yes/no): ", &models.Record{})),
				state:    confirmApproveState,
			},
		},
		{
			name: "delete",
			args: args{
				record:  &models.Record{},
				command: utils.CommandDelete,
			},
			want: want{
				response: []byte(fmt.Sprintf("Do you really want to permanently delete the record '%s'(yes/no): ", &models.Record{})),
				state:    confirmApproveState,
			},
		},
		{
			name: "another command",
			args: args{
				record:  &models.Record{},
				command: utils.CommandServer,
			},
			want: want{
				response: []byte(fmt.Sprintf("Do you really want to commit action with record '%s'(yes/no): ", &models.Record{})),
				state:    confirmApproveState,
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
			assert.Equalf(t, tt.want.state, s.stateConfirmInitial(tt.args.record, tt.args.command), "stateConfirmInitial(%v, %v)", tt.args.record, tt.args.command)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_Confirm(t *testing.T) {
	type args struct {
		record  *models.Record
		command string
	}
	type input struct {
		command string
	}
	type want struct {
		response  []byte
		confirmed bool
		err       assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		args  args
		input input
		want  want
	}{
		{
			name: "base update confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandUpdate,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to update the record '%s'(yes/no): ", &models.Record{})),
				confirmed: true,
				err:       assert.NoError,
			},
		},
		{
			name: "base update not confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandUpdate,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to update the record '%s'(yes/no): ", &models.Record{})),
				confirmed: false,
				err:       assert.NoError,
			},
		},
		{
			name: "base delete confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandDelete,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to permanently delete the record '%s'(yes/no): ", &models.Record{})),
				confirmed: true,
				err:       assert.NoError,
			},
		},
		{
			name: "base delete not confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandDelete,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to permanently delete the record '%s'(yes/no): ", &models.Record{})),
				confirmed: false,
				err:       assert.NoError,
			},
		},
		{
			name: "base unknown confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandServer,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to commit action with record '%s'(yes/no): ", &models.Record{})),
				confirmed: true,
				err:       assert.NoError,
			},
		},
		{
			name: "base unknown not confirmed",
			args: args{
				record:  &models.Record{},
				command: utils.CommandServer,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to commit action with record '%s'(yes/no): ", &models.Record{})),
				confirmed: false,
				err:       assert.NoError,
			},
		},
		{
			name: "cancel",
			args: args{
				record:  &models.Record{},
				command: utils.CommandUpdate,
			},
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to update the record '%s'(yes/no): ", &models.Record{})),
				confirmed: false,
				err:       assert.Error,
			},
		},
		{
			name: "eof",
			args: args{
				record:  &models.Record{},
				command: utils.CommandUpdate,
			},
			input: input{
				command: utils.CommandCancel,
			},
			want: want{
				response:  []byte(fmt.Sprintf("Do you really want to update the record '%s'(yes/no): ", &models.Record{})),
				confirmed: false,
				err:       assert.Error,
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

			confirmed, err := s.Confirm(tt.args.record, tt.args.command)
			if !tt.want.err(t, err, fmt.Sprintf("Confirm(%v, %v)", tt.args.record, tt.args.command)) {
				return
			}
			assert.Equalf(t, tt.want.confirmed, confirmed, "Confirm(%v, %v)", tt.args.record, tt.args.command)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
