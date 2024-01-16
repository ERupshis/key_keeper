package main

import (
	"log"
	"os"

	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
)

func main() {

	reader := interactor.NewReader(os.Stdin)
	writer := interactor.NewWriter(os.Stdout)
	userInteractor := interactor.NewInteractor(reader, writer)

	sm := statemachines.NewStateMachines(userInteractor)
	bankCard := bankcard.NewBankCard(userInteractor, sm)

	cmds := commands.NewCommands(userInteractor, sm, bankCard)

	inMemoryStorage := inmemory.NewStorage()
	mainController := controller.NewController(inMemoryStorage, userInteractor, cmds)

	if err := mainController.Serve(); err != nil {
		log.Fatalf("problem with controller: %v", err)
	}

}
