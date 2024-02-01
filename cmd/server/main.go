package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	authCommon "github.com/erupshis/key_keeper/internal/common/auth"
	"github.com/erupshis/key_keeper/internal/common/auth/authgrpc"
	authPostgres "github.com/erupshis/key_keeper/internal/common/auth/storage/postgres"
	"github.com/erupshis/key_keeper/internal/common/db"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/jwtgenerator"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/internal/server"
	"github.com/erupshis/key_keeper/internal/server/auth"
	"github.com/erupshis/key_keeper/internal/server/config"
	minioS3 "github.com/erupshis/key_keeper/internal/server/storage/binaries/s3/minio"
	"github.com/erupshis/key_keeper/internal/server/storage/records/postgres"
	"github.com/erupshis/key_keeper/internal/server/sync"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip"
)

const (
	migrationsFolder = "file://db/migrations/"
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

	// s3.
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3Login, cfg.S3Password, ""),
		Secure: false,
	})
	if err != nil {
		logs.Fatalf("connect to s3 storage: %v", err)
	}

	bucketManager := minioS3.NewBucketManager(minioClient)
	objectManager := minioS3.NewObjectManager(minioClient)

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
	syncController := sync.NewController(recordsStorage, bucketManager, objectManager)

	// jwt tokens.
	jwtGenerator, err := jwtgenerator.NewJWTGenerator(cfg.JWT, 2)

	// auth.
	hash := hasher.CreateHasher(cfg.HashKey, hasher.TypeSHA256, logs)
	authStorage := authPostgres.NewPostgres(databaseConn, logs)
	authManagerConfig := &authCommon.Config{
		Storage: authStorage,
		JWT:     jwtGenerator,
		Hasher:  hash,
	}
	authManager := authCommon.NewManager(authManagerConfig)
	authController := auth.NewController(authManager)

	// gRPC server options.
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	opts = append(opts, grpc.ChainUnaryInterceptor(
		logger.UnaryServer(logs),
		authgrpc.UnaryServer(jwtGenerator),
	))
	opts = append(opts, grpc.ChainStreamInterceptor(
		logger.StreamServer(logs),
		authgrpc.StreamServer(jwtGenerator),
	))
	// gRPC server
	srv := server.NewGRPCServer(syncController, authController, "grpc", opts...)
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
