package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type Config struct {
	Inmemory *inmemory.Storage
	Local    *local.FileManager

	Interactor *interactor.Interactor
	Cmds       *commands.Commands
}

type Controller struct {
	inmemory *inmemory.Storage
	local    *local.FileManager

	iactr *interactor.Interactor
	cmds  *commands.Commands
}

func NewController(cfg *Config) *Controller {
	return &Controller{
		inmemory: cfg.Inmemory,
		local:    cfg.Local,
		iactr:    cfg.Interactor,
		cmds:     cfg.Cmds,
	}
}

func (c *Controller) Serve(ctx context.Context) error {
	if err := c.cmds.RestoreLocalStorage(ctx, c.inmemory, c.local); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

loop:
	for {
		commandParts, ok := c.iactr.ReadCommand()
		if !ok {
			continue
		}

		switch strings.ToLower(commandParts[0]) {
		case utils.CommandAdd:
			c.cmds.Add(commandParts, c.inmemory)
		case utils.CommandDelete:
			c.cmds.Delete(commandParts, c.inmemory)
		case utils.CommandGet:
			c.cmds.Get(commandParts, c.inmemory)
		case utils.CommandUpdate:
			c.cmds.Update(commandParts, c.inmemory)
		case utils.CommandExit:
			c.iactr.Printf("Exit from app\n")
			break loop
		default:
			if len(commandParts) != 0 && commandParts[0] != "" {
				c.iactr.Printf("Unknown command: '%s'\n", strings.Join(commandParts, " "))
			}
		}
	}

	return nil
}
