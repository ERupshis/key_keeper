package interactor

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInteractor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "base",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writeBuf := bytes.NewBuffer([]byte{})
			writer := NewWriter(writeBuf)

			byteReader := bytes.NewReader([]byte("some stream"))
			reader := NewReader(byteReader)

			interactor := NewInteractor(reader, writer, logger.CreateMock())
			require.NotNil(t, interactor)
			require.NotNil(t, interactor.rd)
			require.NotNil(t, interactor.wr)
		})
	}
}

func TestInteractor_Printf(t *testing.T) {
	type fields struct {
		rd   *Reader
		wr   *Writer
		logs logger.BaseLogger
	}
	type args struct {
		format string
		a      []any
	}
	type want struct {
		n    int
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base string",
			fields: fields{
				rd:   nil,
				wr:   nil,
				logs: logger.CreateMock(),
			},
			args: args{
				format: "%s",
				a:      []any{"some text"},
			},
			want: want{
				n:    9,
				data: []byte("some text"),
			},
		},
		{
			name: "base digit",
			fields: fields{
				rd:   nil,
				wr:   nil,
				logs: logger.CreateMock(),
			},
			args: args{
				format: "%d",
				a:      []any{9},
			},
			want: want{
				n:    1,
				data: []byte("9"),
			},
		},
		{
			name: "empty format",
			fields: fields{
				rd:   nil,
				wr:   nil,
				logs: logger.CreateMock(),
			},
			args: args{
				format: "",
				a:      []any{"some text"},
			},
			want: want{
				n:    26,
				data: []byte("%!(EXTRA string=some text)"),
			},
		},
		{
			name: "empty",
			fields: fields{
				rd:   nil,
				wr:   nil,
				logs: logger.CreateMock(),
			},
			args: args{
				format: "",
				a:      []any{""},
			},
			want: want{
				n:    17,
				data: []byte("%!(EXTRA string=)"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			byteBuff := bytes.NewBuffer(nil)
			i := &Interactor{
				rd:   tt.fields.rd,
				wr:   NewWriter(byteBuff),
				logs: tt.fields.logs,
			}

			n := i.Printf(tt.args.format, tt.args.a...)
			assert.Equalf(t, tt.want.n, n, "incorrect count of writed data: (%d != %d)", tt.want.n, n)
			assert.Equalf(t, tt.want.data, byteBuff.Bytes(), "data was not added in buf")
		})
	}
}

func TestInteractor_Writer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writeBuf := bytes.NewBuffer([]byte{})
			writer := NewWriter(writeBuf)

			byteReader := bytes.NewReader([]byte("some stream"))
			reader := NewReader(byteReader)

			interactor := NewInteractor(reader, writer, logger.CreateMock())
			require.NotNil(t, interactor.Writer())
			require.Equal(t, writer, interactor.Writer())
		})
	}
}

func TestInteractor_GetUserInputAndValidate(t *testing.T) {
	type fields struct {
		wrBuf *bytes.Buffer
		rd    *Reader
		wr    *Writer
		logs  logger.BaseLogger
	}
	type args struct {
		regex *regexp.Regexp
	}
	type want struct {
		text     string
		success  assert.BoolAssertionFunc
		errOccur assert.ErrorAssertionFunc
		errType  error
		buffer   []byte
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
				rd:    NewReader(bytes.NewReader([]byte("add\n"))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "add",
				success:  assert.True,
				errOccur: assert.NoError,
				errType:  nil,
				buffer:   nil,
			},
		},
		{
			name: "base without delim \\n",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("add"))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "",
				success:  assert.True,
				errOccur: assert.Error,
				errType:  io.EOF,
				buffer:   nil,
			},
		},
		{
			name: "base without delim \\n check interruption error",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("add"))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "",
				success:  assert.True,
				errOccur: assert.Error,
				errType:  errs.ErrInterruptedByUser,
				buffer:   nil,
			},
		},
		{
			name: "empty",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("\n"))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "",
				success:  assert.True,
				errOccur: assert.NoError,
				errType:  nil,
				buffer:   nil,
			},
		},
		{
			name: "with spaces",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("  add\n  "))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "add",
				success:  assert.True,
				errOccur: assert.NoError,
				errType:  nil,
				buffer:   nil,
			},
		},
		{
			name: "with spaces 2",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("  add\r\n  "))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "add",
				success:  assert.True,
				errOccur: assert.NoError,
				errType:  nil,
				buffer:   nil,
			},
		},
		{
			name: "cancel",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("cancel\n"))),
				wrBuf: nil,
				wr:    NewWriter(bytes.NewBuffer(nil)),
				logs:  logger.CreateMock(),
			},
			args: args{
				nil,
			},
			want: want{
				text:     "cancel",
				success:  assert.True,
				errOccur: assert.Error,
				errType:  errs.ErrInterruptedByUser,
				buffer:   nil,
			},
		},
		{
			name: "regex check",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("add\n"))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			args: args{
				regexp.MustCompile(`^(yes|no)$`),
			},
			want: want{
				text:     "",
				success:  assert.False,
				errOccur: assert.NoError,
				errType:  nil,
				buffer:   []byte(fmt.Sprintf("incorrect input, try again or interrupt by '%s' command: ", utils.CommandCancel)),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := &Interactor{
				rd:   tt.fields.rd,
				wr:   tt.fields.wr,
				logs: tt.fields.logs,
			}

			if tt.fields.wrBuf != nil {
				i.wr = NewWriter(tt.fields.wrBuf)
			}

			text, successful, err := i.GetUserInputAndValidate(tt.args.regex)
			assert.Equalf(t, tt.want.text, text, "invalid text")
			tt.want.success(t, successful)
			tt.want.errOccur(t, err)

			if tt.want.errType != nil {
				assert.ErrorIs(t, err, tt.want.errType)
			}

			if tt.want.buffer != nil {
				assert.Truef(t, reflect.DeepEqual(tt.want.buffer, tt.fields.wrBuf.Bytes()),
					"write buff is not equal ('%s' != '%s')", tt.want.buffer, tt.fields.wrBuf.Bytes())
			}
		})
	}
}

func TestInteractor_ReadCommand(t *testing.T) {
	type fields struct {
		wrBuf *bytes.Buffer
		rd    *Reader
		wr    *Writer
		logs  logger.BaseLogger
	}
	type want struct {
		parsedCommand []string
		success       assert.BoolAssertionFunc
		buffer        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   want
		want1  bool
	}{
		{
			name: "base",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("add\n"))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			want: want{
				parsedCommand: []string{"add"},
				success:       assert.True,
				buffer:        []byte(fmt.Sprintf("enter command (or '%s'): ", utils.CommandExit)),
			},
		},
		{
			name: "with trim spaces",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("  add\n  "))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			want: want{
				parsedCommand: []string{"add"},
				success:       assert.True,
				buffer:        []byte(fmt.Sprintf("enter command (or '%s'): ", utils.CommandExit)),
			},
		},
		{
			name: "two-word command",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("add creds\n"))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			want: want{
				parsedCommand: []string{"add", "creds"},
				success:       assert.True,
				buffer:        []byte(fmt.Sprintf("enter command (or '%s'): ", utils.CommandExit)),
			},
		},
		{
			name: "new row symbol",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte("\n"))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			want: want{
				parsedCommand: nil,
				success:       assert.False,
				buffer:        []byte(fmt.Sprintf("enter command (or '%s'): Empty command. Try again\n", utils.CommandExit)),
			},
		},
		{
			name: "empty",
			fields: fields{
				rd:    NewReader(bytes.NewReader([]byte(""))),
				wrBuf: bytes.NewBuffer(nil),
				wr:    nil,
				logs:  logger.CreateMock(),
			},
			want: want{
				parsedCommand: []string{"exit"},
				success:       assert.True,
				buffer:        []byte(fmt.Sprintf("enter command (or '%s'): ", utils.CommandExit)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := &Interactor{
				rd:   tt.fields.rd,
				wr:   tt.fields.wr,
				logs: tt.fields.logs,
			}

			if tt.fields.wrBuf != nil {
				i.wr = NewWriter(tt.fields.wrBuf)
			}

			command, ok := i.ReadCommand()
			assert.Equalf(t, tt.want.parsedCommand, command, "ReadCommand()")
			tt.want.success(t, ok)

			if tt.want.buffer != nil {
				assert.Truef(t, reflect.DeepEqual(tt.want.buffer, tt.fields.wrBuf.Bytes()),
					"write buff is not equal ('%s' != '%s')", tt.want.buffer, tt.fields.wrBuf.Bytes())
			}
		})
	}
}
