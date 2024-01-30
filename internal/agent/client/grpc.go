package client

import (
	"context"
	"fmt"
	"io"

	clientModels "github.com/erupshis/key_keeper/internal/agent/client/models"
	"github.com/erupshis/key_keeper/internal/agent/models"
	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	_ BaseClient = (*GRPC)(nil)
)

type GRPC struct {
	syncClient pb.SyncClient
	authClient pb.AuthClient
	conn       *grpc.ClientConn
}

func NewGRPC(address string, options ...grpc.DialOption) (BaseClient, error) {
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		return nil, fmt.Errorf("create connection to server: %w", err)
	}

	syncClient := pb.NewSyncClient(conn)
	authClient := pb.NewAuthClient(conn)

	return &GRPC{
		syncClient: syncClient,
		authClient: authClient,
		conn:       conn,
	}, nil
}

func (g *GRPC) Close() error {
	return g.conn.Close()
}

func (g *GRPC) Login(ctx context.Context, creds *models.Credential) error {
	_, err := g.authClient.Login(ctx, &pb.LoginRequest{Creds: clientModels.ConvertCredentialToGRPC(creds)})
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}

func (g *GRPC) Register(ctx context.Context, creds *models.Credential) error {
	_, err := g.authClient.Register(ctx, &pb.RegisterRequest{Creds: clientModels.ConvertCredentialToGRPC(creds)})
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	return nil
}

func (g *GRPC) Push(ctx context.Context, storageRecords []localModels.StorageRecord) error {
	stream, err := g.syncClient.Push(ctx)
	if err != nil {
		return fmt.Errorf("push records: %w", err)
	}

	defer deferutils.ExecSilent(stream.CloseSend)

	for _, record := range storageRecords {
		record := record
		err = stream.Send(&pb.PushRequest{Record: clientModels.ConvertStorageRecordToGRPC(&record)})
		if err != nil {
			return fmt.Errorf("send record: %w", err)
		}
	}

	return nil
}

func (g *GRPC) Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error) {
	stream, err := g.syncClient.Pull(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("pull records: %w", err)
	}

	res := make(map[int64]localModels.StorageRecord)

	tmpReceive := &pb.PullResponse{}
	for {
		tmpReceive, err = stream.Recv()
		if err != nil {
			break
		}

		record := clientModels.ConvertStorageRecordFromGRPC(tmpReceive.GetRecord())
		res[record.ID] = *record
	}

	if err == io.EOF {
		return res, nil
	}

	return nil, fmt.Errorf("receive record: %w", err)
}

func (g *GRPC) PushBinary(ctx context.Context, binaries map[string]string) error {
	return nil
}

func (g *GRPC) PullBinary(ctx context.Context) error {
	// TODO: need to perform file loading and syncing.
	return nil
}
