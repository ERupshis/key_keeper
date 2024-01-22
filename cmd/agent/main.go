package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/erupshis/key_keeper/internal/agent/config"
	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/credential"
	localCmd "github.com/erupshis/key_keeper/internal/agent/controller/commands/local"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
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

	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("error parse config: %v", err)
		return
	}

	logs, err := logger.NewZap("info")
	if err != nil {
		log.Fatalf("create zap logs: %v", err)
	}
	defer deferutils.ExecSilent(logs.Sync)

	reader := interactor.NewReader(os.Stdin)
	writer := interactor.NewWriter(os.Stdout)
	userInteractor := interactor.NewInteractor(reader, writer)

	sm := statemachines.NewStateMachines(userInteractor)
	bankCard := bankcard.NewBankCard(userInteractor, sm)
	cred := credential.NewCredentials(userInteractor, sm)
	cmdLocal := localCmd.NewLocal(userInteractor)

	cmdConfig := commands.Config{
		StateMachines:   sm,
		BankCard:        bankCard,
		Credential:      cred,
		LocalStorageCmd: cmdLocal,
	}
	cmds := commands.NewCommands(userInteractor, &cmdConfig)

	inMemoryStorage := inmemory.NewStorage()

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	localAutoSaveConfig := local.AutoSaveConfig{
		SaveInterval:    cfg.LocalStoreInterval,
		InMemoryStorage: inMemoryStorage,
		Logs:            logs,
	}

	dataCryptor := ska.NewSKA("some user key", ska.Key16)
	localStorage := local.NewFileManager(cfg.LocalStoragePath, logs, &localAutoSaveConfig, dataCryptor)

	controllerConfig := controller.Config{
		Inmemory:   inMemoryStorage,
		Local:      localStorage,
		Interactor: userInteractor,
		Cmds:       cmds,
	}

	mainController := controller.NewController(&controllerConfig)

	if err := mainController.Serve(ctxWithCancel); err != nil {
		log.Fatalf("problem with controller: %v", err)
	}

}
