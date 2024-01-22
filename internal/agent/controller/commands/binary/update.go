package binary

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (b *Binary) ProcessUpdateCommand(record *data.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: b.addMainData,
	}

	return b.sm.Add(cfg)
}
