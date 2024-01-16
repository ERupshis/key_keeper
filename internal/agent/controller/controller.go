package controller

import (
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type Controller struct {
	inmemory *inmemory.Storage
	iactr    *interactor.Interactor
	cmds     *commands.Commands
}

func NewController(inmemory *inmemory.Storage, interactor *interactor.Interactor, cmds *commands.Commands) *Controller {
	return &Controller{
		inmemory: inmemory,
		iactr:    interactor,
		cmds:     cmds,
	}
}

func (c *Controller) Serve() error {
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
