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

	for {
		select {
		case <-ctx.Done():
			return c.saveRecordsLocally(ctx)
		default:
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
				return c.saveRecordsLocally(ctx) // TODO: need to move in stop controller func.
			default:
				if len(commandParts) != 0 && commandParts[0] != "" {
					c.iactr.Printf("Unknown command: '%s'\n", strings.Join(commandParts, " "))
				}
			}
		}
	}
}

func (c *Controller) saveRecordsLocally(ctx context.Context) error {
	errMsg := "save records locally: %w"
	records, err := c.inmemory.GetAllRecords()
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if err = c.local.SaveUserData(ctx, records); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
