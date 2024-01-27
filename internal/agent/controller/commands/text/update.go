package text

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/models"
)

func (t *Text) ProcessUpdateCommand(record *models.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: t.addMainData,
	}

	return t.sm.Add(cfg)
}
