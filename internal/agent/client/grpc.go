package client

import (
	"context"
	"fmt"
	"io"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	_ BaseClient = (*GRPC)(nil)
)

type GRPC struct {
	client pb.SyncClient
	conn   *grpc.ClientConn
}

func NewGRPC(address string, options ...grpc.DialOption) (BaseClient, error) {
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		return nil, fmt.Errorf("create connection to server: %w", err)
	}

	client := pb.NewSyncClient(conn)

	return &GRPC{
		client: client,
		conn:   conn,
	}, nil
}

func (g *GRPC) Close() error {
	return g.conn.Close()
}

func (g *GRPC) Push(ctx context.Context, storageRecords []localModels.StorageRecord) error {
	stream, err := g.client.Push(ctx)
	if err != nil {
		return fmt.Errorf("push records: %w", err)
	}

	defer func() {
		_ = stream.CloseSend()
	}()

	for _, record := range storageRecords {
		record := record
		err = stream.Send(&pb.PushRequest{Record: localModels.ConvertStorageRecordToGRPC(&record)})
		if err != nil {
			// TODO: collect all errors?
			return fmt.Errorf("send record: %w", err)
		}
	}

	return nil
}

func (g *GRPC) Pull(ctx context.Context) (map[int64]localModels.StorageRecord, error) {
	stream, err := g.client.Pull(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("pull records: %w", err)
	}

	res := make(map[int64]localModels.StorageRecord)
	for {
		received, err := stream.Recv()
		if err != io.EOF {
			break
		}
		// TODO: collect all errors?
		record := localModels.ConvertStorageRecordFromGRPC(received.GetRecord())
		res[record.ID] = *record
	}

	return res, nil
}
