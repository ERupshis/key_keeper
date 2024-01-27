package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/grpc/interceptors/logging"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/internal/server"
	"github.com/erupshis/key_keeper/internal/server/config"
	"github.com/erupshis/key_keeper/internal/server/storage/postgres"
	"github.com/erupshis/key_keeper/internal/server/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip"
)

const (
	migrationsFolder = "file://db/records/migrations/"
)

func main() {
	logs, err := logger.NewZap("info")
	if err != nil {
		log.Fatalf("create zap logs: %v", err)
	}
	defer deferutils.ExecSilent(logs.Sync)

	cfg, err := config.Parse()
	if err != nil {
		logs.Fatalf("parse config: %v", err)
	}

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	// storage.
	dbConfig := db.Config{
		DSN:              cfg.DatabaseDSN,
		MigrationsFolder: migrationsFolder,
	}
	databaseConn, err := db.NewConnection(ctxWithCancel, dbConfig)
	if err != nil {
		logs.Fatalf("failed to connect to users database: %v", err)
	}

	recordsStorage := postgres.NewPostgres(databaseConn, logs)

	// handlers controller.
	syncController := sync.NewController(recordsStorage)

	// gRPC server options.
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	opts = append(opts, grpc.ChainUnaryInterceptor(logging.UnaryServer(logs)))
	opts = append(opts, grpc.ChainStreamInterceptor(logging.StreamServer(logs)))
	// gRPC server
	srv := server.NewGRPCServer(syncController, "grpc", opts...)
	srv.Host(cfg.Host)

	go func() {
		listener, err := net.Listen("tcp", cfg.Host)
		if err != nil {
			logs.Fatalf("failed to listen for %s dataprovider: %v", srv.GetInfo(), err)
		}

		if err = srv.Serve(listener); err != nil {
			logs.Infof("http://%s dataprovider refused to start or stop with error: %v", srv.GetInfo(), err)
			return
		}
	}()

	// shutdown.
	idleConnsClosed := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigCh

		if err = srv.GracefulStop(ctxWithCancel); err != nil {
			logs.Infof("%s dataprovider graceful stop error: %v", srv.GetInfo(), err)
		}

		cancel()
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
	logs.Infof("key_keeper service shutdown gracefully")
}
