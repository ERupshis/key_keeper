package credential

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/common/data"
)

func (с *Credential) ProcessUpdateCommand(record *data.Record) error {
	cfg := statemachines.AddConfig{
		Record:   record,
		MainData: с.addMainData,
	}

	return с.sm.Add(cfg)
}
