package controller

import (
	"fmt"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type Controller struct {
	inmemory *inmemory.Storage
}

func NewController(inmemory *inmemory.Storage) *Controller {
	return &Controller{inmemory: inmemory}
}

func (c *Controller) Serve() error {
loop:
	for {
		commandParts, ok := readCommand()
		if !ok {
			continue
		}

		switch strings.ToLower(commandParts[0]) {
		case utils.CommandAdd:
			commands.Add(commandParts, c.inmemory)
		case utils.CommandGet:
			commands.Get(commandParts, c.inmemory)
		case utils.CommandExit:
			fmt.Printf("Exit from app\n")
			break loop
		default:
			if len(commandParts) != 0 && commandParts[0] != "" {
				fmt.Printf("Unknown command: '%s'\n", strings.Join(commandParts, " "))
			}
		}
	}

	return nil
}

func readCommand() ([]string, bool) {
	fmt.Printf("Insert command (or '%s'): ", utils.CommandExit)
	command, _, _ := utils.GetUserInputAndValidate(nil)
	command = strings.TrimSpace(command)
	commandParts := strings.Split(command, " ")
	if len(commandParts) == 0 {
		fmt.Printf("Empty command. Try again\n")
		return nil, false
	}

	return commandParts, true
}
