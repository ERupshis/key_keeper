package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/key_keeper/internal/agent/client"
	"github.com/erupshis/key_keeper/internal/agent/config"
	"github.com/erupshis/key_keeper/internal/agent/controller"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/bankcard"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/binary"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/credential"
	localCmd "github.com/erupshis/key_keeper/internal/agent/controller/commands/local"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/server"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/text"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/agent/storage/binaries"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
	"github.com/erupshis/key_keeper/internal/common/auth/authgrpc"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
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

	reader := interactor.NewReader(bufio.NewReader(os.Stdin))
	writer := interactor.NewWriter(bufio.NewWriter(os.Stdout))
	userInteractor := interactor.NewInteractor(reader, writer, logs)

	sm := statemachines.NewStateMachines(userInteractor)
	bankCard := bankcard.NewBankCard(userInteractor, sm)
	cred := credential.NewCredentials(userInteractor, sm)
	txt := text.NewText(userInteractor, sm)

	dataCryptor := ska.NewSKA("some user key", ska.Key16) // TODO: need to move param in config.
	hash := hasher.CreateHasher(cfg.HashKey, hasher.TypeSHA256, logs)

	binaryConfig := binary.Config{
		Iactr:     userInteractor,
		Sm:        sm,
		Hash:      hash,
		Cryptor:   dataCryptor,
		StorePath: cfg.LocalStoragePath,
	}
	bin := binary.NewBinary(&binaryConfig)

	inMemoryStorage := inmemory.NewStorage(dataCryptor)
	binaryManager := binaries.NewBinaryManager(cfg.LocalStoragePath)
	localAutoSaveConfig := local.AutoSaveConfig{
		SaveInterval:    cfg.LocalStoreInterval,
		InMemoryStorage: inMemoryStorage,
		BinaryManager:   binaryManager,
		Logs:            logs,
	}
	localStorage := local.NewFileManager(cfg.LocalStoragePath, logs, userInteractor, &localAutoSaveConfig, dataCryptor)

	cmdLocal := localCmd.NewLocal(userInteractor)

	authInterceptor := authgrpc.NewClientInterceptor()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithChainUnaryInterceptor(
		logger.UnaryClient(logs),
		authInterceptor.UnaryClient(),
	))
	opts = append(opts, grpc.WithChainStreamInterceptor(
		logger.StreamClient(logs),
		authInterceptor.StreamClient(),
	))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	grpcClient, err := client.NewGRPC(cfg.ServerHost, opts...)
	if err != nil {
		logs.Fatalf("client: %v", err)
	}
	defer deferutils.ExecWithLogError(grpcClient.Close, logs)

	serverCommandConfig := server.Config{
		Inmemory: inMemoryStorage,
		Local:    localStorage,
		Client:   grpcClient,
		Iactr:    userInteractor,
		Binary:   binaryManager,
	}
	serverCommand := server.NewServer(&serverCommandConfig)

	cmdConfig := commands.Config{
		StateMachines:   sm,
		BankCard:        bankCard,
		Credential:      cred,
		Text:            txt,
		Binary:          bin,
		LocalStorageCmd: cmdLocal,
		Server:          serverCommand,
	}
	cmds := commands.NewCommands(userInteractor, &cmdConfig)

	controllerConfig := controller.Config{
		Inmemory:   inMemoryStorage,
		Local:      localStorage,
		Interactor: userInteractor,
		Cmds:       cmds,
	}
	mainController := controller.NewController(&controllerConfig)

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// shutdown.
	idleConnsClosed := make(chan struct{})
	sigCh := make(chan os.Signal, 5)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigCh
		cancel()
		if err = mainController.SaveRecordsLocally(); err != nil {
			logs.Infof("failed to save records locally: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err = mainController.Serve(ctxWithCancel); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
		logs.Infof("problem with controller: %v", err)
	}
	sigCh <- syscall.SIGQUIT

	<-idleConnsClosed
	logs.Infof("agent shutdown gracefully")
}
