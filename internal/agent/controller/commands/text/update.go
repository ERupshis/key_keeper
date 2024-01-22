package text

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (t *Text) ProcessUpdateCommand(record *data.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: t.addMainData,
	}

	return t.sm.Add(cfg)
}
