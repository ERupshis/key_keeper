package interactor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
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
			writer := bytes.NewBuffer([]byte{})
			got := NewWriter(writer)
			require.NotNil(t, got)
		})
	}
}
