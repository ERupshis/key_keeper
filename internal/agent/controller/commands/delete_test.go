package commands

import (
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCommands_confirmAndDeleteByID(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		recordInBase   *models.Record
		recordToDelete *models.Record
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
			name: "base confirmed",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase:   &models.Record{ID: -1},
				recordToDelete: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record successfully deleted\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "base not confirmed",
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			args: args{
				recordInBase:   &models.Record{ID: -1},
				recordToDelete: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record deleting was interrupted by user\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "confirm error",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				recordInBase:   &models.Record{ID: -1},
				recordToDelete: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): "),
				err:      assert.Error,
			},
		},
		{
			name: "storage error",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase:   &models.Record{ID: -1},
				recordToDelete: &models.Record{ID: -2},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -2, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): "),
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase))
			tt.want.err(t, c.confirmAndDeleteByID(tt.args.recordToDelete, inMemoryStorage), fmt.Sprintf("confirmAndDeleteByID(%v, %v)", tt.args.recordToDelete, inMemoryStorage))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_findAndDeleteRecordByID(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		id           int64
		recordInBase *models.Record
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
			name: "base confirm",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id:           -1,
				recordInBase: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record successfully deleted\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "base not confirm",
			input: input{
				command: testutils.AddNewRow(utils.CommandNo),
			},
			args: args{
				id:           -1,
				recordInBase: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record deleting was interrupted by user\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "search record in storage error",
			input: input{
				command: testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				id:           -2,
				recordInBase: &models.Record{ID: -1},
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

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase))
			tt.want.err(t, c.findAndDeleteRecordByID(tt.args.id, inMemoryStorage), fmt.Sprintf("findAndDeleteRecordByID(%v, %v)", tt.args.id, inMemoryStorage))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_handleDelete(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		recordInBase *models.Record
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
				command: testutils.AddNewRow("-1") + testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("enter record id: Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record successfully deleted\n"),
				err:      assert.NoError,
			},
		},
		{
			name: "missing records id in storage",
			input: input{
				command: testutils.AddNewRow("-2") + testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{ID: -1},
			},
			want: want{
				response: []byte("enter record id: "),
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase))
			tt.want.err(t, c.handleDelete(inMemoryStorage), fmt.Sprintf("handleDelete(%v)", inMemoryStorage))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestCommands_Delete(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		parts        []string
		recordInBase *models.Record
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
			name: "base",
			input: input{
				command: testutils.AddNewRow("-1") + testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{ID: -1},
				parts:        []string{utils.CommandDelete},
			},
			want: want{
				response: []byte("enter record id: Do you really want to permanently delete the record '{ID: -1, MetaData: INVALID}%!(EXTRA models.MetaData=map[])'(yes/no): Record successfully deleted\n"),
			},
		},
		{
			name: "incorrect parts count",
			input: input{
				command: testutils.AddNewRow("-1") + testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{ID: -1},
				parts:        []string{utils.CommandDelete, utils.CommandID},
			},
			want: want{
				response: []byte("incorrect request. should contain command 'delete' only\n"),
			},
		},
		{
			name: "incorrect id entered",
			input: input{
				command: testutils.AddNewRow("-2") + testutils.AddNewRow(utils.CommandYes),
			},
			args: args{
				recordInBase: &models.Record{ID: -1},
				parts:        []string{utils.CommandDelete},
			},
			want: want{
				response: []byte("enter record id: request processing error: process 'delete' command: process 'get' command: record not found\n"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, inMemoryStorage, writer := getCommands(tt.input.command)

			assert.NoError(t, inMemoryStorage.AddRecord(tt.args.recordInBase))
			c.Delete(tt.args.parts, inMemoryStorage)
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
