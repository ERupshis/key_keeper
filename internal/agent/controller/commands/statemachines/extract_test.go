package statemachines

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/agent/utils/testutils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

const (
	pathNotAbs = "text.txt"
)

var (
	wd, _         = os.Getwd()
	pathAbs       = fmt.Sprintf("%s%ctext.txt", wd, filepath.Separator)
	pathFormatted = fmt.Sprintf("%s%c", wd, filepath.Separator)
)

func TestStateMachines_stateFilePath(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		state    extractPathState
		path     string
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
				command: testutils.AddNewRow(pathAbs),
			},
			want: want{
				response: nil,
				state:    extractPathFinishState,
				path:     pathFormatted,
				err:      assert.NoError,
			},
		},
		{
			name: "not abs path",
			input: input{
				command: testutils.AddNewRow(pathNotAbs),
			},
			want: want{
				response: []byte("entered local path. Try to set absolute path: "),
				state:    extractPathFilePathState,
				path:     "",
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
				state:    extractPathFilePathState,
				path:     "",
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
				state:    extractPathFilePathState,
				path:     "",
				err:      assert.Error,
			},
		},
		{
			name: "empty",
			input: input{
				command: "",
			},
			want: want{
				response: nil,
				state:    extractPathFilePathState,
				path:     "",
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

			state, path, err := s.stateFilePath()
			if !tt.want.err(t, err, "stateFilePath()") {
				return
			}

			assert.Equalf(t, tt.want.state, state, "state fail")
			assert.Equalf(t, tt.want.path, path, "path fail")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_stateInitial(t *testing.T) {
	type want struct {
		response []byte
		state    extractPathState
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "base",
			want: want{
				response: []byte("enter absolute path to file: "),
				state:    extractPathFilePathState,
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

			assert.Equalf(t, tt.want.state, s.stateInitial(), "stateInitial()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_extractFilePath(t *testing.T) {
	type input struct {
		command string
	}
	type want struct {
		response []byte
		path     string
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
				command: testutils.AddNewRow(pathAbs),
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				path:     pathFormatted,
				err:      assert.NoError,
			},
		},
		{
			name: "err from input parser",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				path:     "",
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

			path, err := s.extractFilePath()
			if !tt.want.err(t, err, "extractFilePath()") {
				return
			}
			assert.Equalf(t, tt.want.path, path, "extractFilePath()")
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}

func TestStateMachines_Extract(t *testing.T) {
	type input struct {
		command string
	}
	type args struct {
		cfg ExtractConfig
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
				command: testutils.AddNewRow(pathAbs),
			},
			args: args{
				cfg: ExtractConfig{
					Record: &models.Record{},
					FileSave: func(record *models.Record, savePath string) error {
						return nil
					},
				},
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				err:      assert.NoError,
			},
		},
		{
			name: "cancel",
			input: input{
				command: testutils.AddNewRow(utils.CommandCancel),
			},
			args: args{
				cfg: ExtractConfig{
					Record: &models.Record{},
					FileSave: func(record *models.Record, savePath string) error {
						return nil
					},
				},
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				err:      assert.Error,
			},
		},
		{
			name: "eof",
			input: input{
				command: "",
			},
			args: args{
				cfg: ExtractConfig{
					Record: &models.Record{},
					FileSave: func(record *models.Record, savePath string) error {
						return nil
					},
				},
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				err:      assert.Error,
			},
		},
		{
			name: "cancel on fileSave",
			input: input{
				command: testutils.AddNewRow(pathAbs),
			},
			args: args{
				cfg: ExtractConfig{
					Record: &models.Record{},
					FileSave: func(record *models.Record, savePath string) error {
						return errs.ErrInterruptedByUser
					},
				},
			},
			want: want{
				response: []byte("enter absolute path to file: "),
				err:      assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input.command))
			writer := bytes.NewBuffer(nil)
			iactr := testutils.CreateUserInteractor(reader, writer, logger.CreateMock())

			s := &StateMachines{
				iactr: iactr,
			}

			tt.want.err(t, s.Extract(tt.args.cfg), fmt.Sprintf("Extract(%v)", tt.args.cfg))
			assert.Equal(t, tt.want.response, writer.Bytes(), "response fail")
		})
	}
}
