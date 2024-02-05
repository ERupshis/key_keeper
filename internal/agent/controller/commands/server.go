package commands

import (
	"context"
	"fmt"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

func (c *Commands) Server(ctx context.Context, parts []string) {
	supportedTypes := []string{utils.CommandPush, utils.CommandPull, utils.CommandLogin, utils.CommandRegister}
	if len(parts) != 2 {
		c.iactr.Printf("incorrect request. should contain command '%s' and action type(%s)\n", utils.CommandServer, supportedTypes)
		return
	}

	if err := c.handleServer(ctx, parts[1]); err != nil {
		c.handleCommandError(err, utils.CommandServer, supportedTypes)
		return
	}

	c.iactr.Printf("command %s %s done\n", parts[0], parts[1])
}

func (c *Commands) handleServer(ctx context.Context, actionType string) error {
	var err error
	switch actionType {
	case utils.CommandPush:
		err = c.server.ProcessPushCommand(ctx)
	case utils.CommandPull:
		err = c.server.ProcessPullCommand(ctx)
	case utils.CommandLogin:
		err = c.server.ProcessLoginCommand(ctx)
	case utils.CommandRegister:
		err = c.server.ProcessRegisterCommand(ctx)
	default:
		err = fmt.Errorf(errs.ErrProcessMsgBody, utils.CommandServer, errs.ErrIncorrectServerActionType)
	}

	return err
}
