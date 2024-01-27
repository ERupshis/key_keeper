package sync

import (
	"context"
	"fmt"
	"io"
	"strconv"

	localModels "github.com/erupshis/key_keeper/internal/agent/storage/models"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/internal/server/storage"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	metaUserID = "user_id"
)

type Controller struct {
	pb.UnimplementedSyncServer

	storage storage.BaseStorage
}

func NewController(storage storage.BaseStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (c *Controller) Push(stream pb.Sync_PushServer) error {
	defer deferutils.ExecSilent(func() error {
		return stream.SendAndClose(&emptypb.Empty{})
	})

	userID, err := getUserID(stream.Context())
	if err != nil {
		return err
	}

	for {
		tmpReceive := &pb.PushRequest{}
		tmpReceive, err = stream.Recv()
		if err != nil {
			break
		}

		record := localModels.ConvertStorageRecordFromGRPC(tmpReceive.GetRecord())
		if err = c.storage.UpsertRecord(stream.Context(), userID, record); err != nil {
			break
		}
	}

	if err == io.EOF {
		return nil
	}

	return status.Errorf(codes.Internal, "receive record: %v", err)
}

func (c *Controller) Pull(_ *emptypb.Empty, stream pb.Sync_PullServer) error {
	userID, err := getUserID(stream.Context())
	if err != nil {
		return err
	}

	records, err := c.storage.GetRecords(stream.Context(), userID)
	for idx := range records {
		err = stream.Send(&pb.PullResponse{Record: localModels.ConvertStorageRecordToGRPC(&records[idx])})
		if err != nil {
			return status.Errorf(codes.Internal, "send record: %v", err)
		}
	}

	return nil
}
func (c *Controller) PushBinary(_ pb.Sync_PushBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method PushBinary not implemented")
}
func (c *Controller) PullBinary(_ *emptypb.Empty, _ pb.Sync_PullBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method PullBinary not implemented")
}

func getUserID(ctx context.Context) (int64, error) {
	rawUserID := ctx.Value(metaUserID)
	if rawUserID == "" {
		return -1, status.Errorf(codes.Unauthenticated, "missing user id")
	}

	userID, err := strconv.ParseInt(fmt.Sprintf("%s", rawUserID), 10, 64)
	if err != nil {
		return -1, status.Errorf(codes.InvalidArgument, "bad user id")
	}

	return userID, nil
}
