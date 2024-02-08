package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBinary(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "base",
			args: args{
				cfg: &Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewBinary(tt.args.cfg))
		})
	}
}
