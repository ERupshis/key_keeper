package deferutils

import (
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/common/logger"
)

func TestExecuteWithLogError(t *testing.T) {
	type args struct {
		callback func() error
		log      logger.BaseLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid def case",
			args: args{
				callback: func() error {
					return nil
				},
				log: logger.CreateMock(),
			},
		},
		{
			name: "error from callback",
			args: args{
				callback: func() error {
					return fmt.Errorf("test err")
				},
				log: logger.CreateMock(),
			},
		},
	}
	for _, ttTmp := range tests {
		tt := ttTmp
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ExecWithLogError(tt.args.callback, tt.args.log)
		})
	}
}

func TestExecSilent(t *testing.T) {
	type args struct {
		callback func() error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid def case",
			args: args{
				callback: func() error {
					return nil
				},
			},
		},
		{
			name: "error from callback",
			args: args{
				callback: func() error {
					return fmt.Errorf("test err")
				},
			},
		},
	}
	for _, ttTmp := range tests {
		tt := ttTmp
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ExecSilent(tt.args.callback)
		})
	}
}
