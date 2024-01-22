package credential

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (c *Credential) ProcessUpdateCommand(record *data.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: c.addMainData,
	}

	return c.sm.Add(cfg)
}
