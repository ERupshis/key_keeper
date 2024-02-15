package interactor

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReader(t *testing.T) {
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
			reader := bytes.NewReader([]byte("some stream"))
			got := NewReader(reader)
			require.NotNil(t, got)
		})
	}
}

func TestReader_getUserInput(t *testing.T) {
	type args struct {
		Reader *bufio.Reader
	}
	type want struct {
		data     string
		errOccur bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base",
			args: args{
				Reader: bufio.NewReader(bytes.NewReader([]byte("some stream\n"))),
			},
			want: want{
				data:     "some stream\n",
				errOccur: false,
			},
		},
		{
			name: "empty EOF",
			args: args{
				Reader: bufio.NewReader(bytes.NewReader([]byte(""))),
			},
			want: want{
				data:     "",
				errOccur: true,
			},
		},
		{
			name: "without delim \n EOF",
			args: args{
				Reader: bufio.NewReader(bytes.NewReader([]byte("some stream"))),
			},
			want: want{
				data:     "",
				errOccur: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Reader{
				Reader: tt.args.Reader,
			}

			got, err := r.getUserInput()
			if (err != nil) != tt.want.errOccur {
				t.Errorf("getUserInput() error = %v, wantErr %v", err, tt.want.errOccur)
				return
			}
			if got != tt.want.data {
				t.Errorf("getUserInput() got = %v, want %v", got, tt.want)
			}
		})
	}
}
