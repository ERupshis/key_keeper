package controller

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/data"
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
			ok, needToContinueExecution := c.commandAdd(parts)
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

func (c *Controller) commandAdd(parts []string) (bool, bool) {
	if len(parts) != 2 {
		fmt.Println("incorrect request. Should contain command 'add' and object type('bank', 'cards', 'text', 'bin').")
		return true, false
	}

	record, err := c.processAddCommand(data.ConvertStringToRecordType(parts[1]))
	if err != nil {
		if errors.Is(err, errs.ErrInterruptedByUser) {
			fmt.Printf("add operation was canceled by user\n")
			return false, false
		}

		fmt.Printf("request parsing error: %v\n", err)
		if errors.Is(err, errs.ErrIncorrectRecordType) {
			fmt.Printf("only ('bank', 'cards', 'text', 'bin') are supported\n")
		}
		return false, false
	}

	fmt.Printf("record added: %+v\n", record)
	return true, true
}

func (c *Controller) processAddCommand(recordType data.RecordType) (*data.Record, error) {
	errMsg := "process add command: %w"

	newRecord := &data.Record{
		Id: -1,
	}

	var err error
	switch recordType {
	case data.TypeBankCard:
		err = bankcard.ProcessAddCommand(newRecord)
	default:
		return nil, fmt.Errorf(errMsg, errs.ErrIncorrectRecordType)
	}

	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	if err = c.inmemory.AddRecord(newRecord); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return newRecord, nil
}

func (c *Controller) processCommand(command string) {
	fmt.Printf("Вы ввели команду: %s\n", command)
	// Здесь можно добавить логику для обработки конкретных команд
}
