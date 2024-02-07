package bankcard

import (
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/stretchr/testify/assert"
)

func TestNewBankCard(t *testing.T) {
	type args struct {
		iactr    *interactor.Interactor
		machines *statemachines.StateMachines
	}
	tests := []struct {
		name string
		args args
		want *BankCard
	}{
		{
			name: "base",
			args: args{
				iactr:    nil,
				machines: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewBankCard(tt.args.iactr, tt.args.machines))
		})
	}
}
