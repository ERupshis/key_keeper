package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type Config struct {
	Inmemory *inmemory.Storage
	Local    *local.FileManager

	Client client.BaseClient

	Interactor *interactor.Interactor
	Cmds       *commands.Commands
}

type Controller struct {
	inmemory *inmemory.Storage
	local    *local.FileManager

	client client.BaseClient

	iactr *interactor.Interactor
	cmds  *commands.Commands
}

func NewController(cfg *Config) *Controller {
	return &Controller{
		inmemory: cfg.Inmemory,
		local:    cfg.Local,
		client:   cfg.Client,
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
			return nil
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
			case utils.CommandExtract:
				c.cmds.Extract(commandParts, c.inmemory)
			case utils.CommandGet:
				c.cmds.Get(commandParts, c.inmemory)
			case utils.CommandPull:
				c.cmds.Pull(ctx, c.client, c.inmemory)
			case utils.CommandPush:
				c.cmds.Push(ctx, c.client, c.inmemory)
			case utils.CommandUpdate:
				c.cmds.Update(commandParts, c.inmemory)
			case utils.CommandExit:
				c.iactr.Printf("Exit from app\n")
				return nil
			default:
				if len(commandParts) != 0 && commandParts[0] != "" {
					c.iactr.Printf("Unknown command: '%s'\n", strings.Join(commandParts, " "))
				}
			}
		}
	}
}

func (c *Controller) SaveRecordsLocally() error {
	errMsg := "save records locally: %w"
	records, err := c.inmemory.GetAllRecords()
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if records == nil {
		return nil
	}

	if err = c.local.SaveUserData(records); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
