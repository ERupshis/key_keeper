package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	reader := bufio.NewReader(os.Stdin)
loop:
	for {
		fmt.Print("Введите команду (exit для выхода): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		command := strings.TrimSpace(input)
		parts := strings.Split(command, " ")
		if len(parts) == 0 {
			fmt.Println("Empty command. Try again.")
			continue
		}

		switch strings.ToLower(parts[0]) {
		case utils.CommandAdd:
			ok, needToContinueExecution := c.handleAddCommand(parts)
			if !needToContinueExecution {
				if ok {
					continue
				} else {
					break
				}
			}
		case utils.CommandGet:
			ok, needToContinueExecution := c.handleGetCommand(parts)
			if !needToContinueExecution {
				if ok {
					continue
				} else {
					break
				}
			}
		case utils.CommandExit:
			fmt.Println("Выход из приложения.")
			break loop
		default:
			fmt.Println("Unknown command.")
		}

		c.processCommand(command)
	}

	return nil
}

func (c *Controller) processCommand(command string) {
	fmt.Printf("Вы ввели команду: %s\n", command)
	// Здесь можно добавить логику для обработки конкретных команд
}
