package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	// example of run: go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(cmd.exe /c "echo %DATE%")' -X 'main.buildCommit=$(git rev-parse HEAD)'" main.go
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	// cfg, err := config.Parse()
	// if err != nil {
	// 	log.Fatalf("error parse config: %v", err)
	// 	return
	// }

	logger, err := logger.NewZap("info")
	if err != nil {
		log.Fatalf("create zap logger: %v", err)
	}
	defer deferutils.ExecSilent(logger.Sync)

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
