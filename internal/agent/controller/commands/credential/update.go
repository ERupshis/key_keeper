package credential

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/models"
)

func (c *Credential) ProcessUpdateCommand(record *models.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: c.addMainData,
	}

	return c.sm.Add(cfg)
}
