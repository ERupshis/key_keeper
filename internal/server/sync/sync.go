package sync

import (
	"context"
	"fmt"
	"io"
	"strconv"

	clientModels "github.com/erupshis/key_keeper/internal/agent/client/models"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/internal/server/storage/binaries/models"
	minioS3 "github.com/erupshis/key_keeper/internal/server/storage/binaries/s3/minio"
	"github.com/erupshis/key_keeper/internal/server/storage/records"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	metaUserID = "user_id"
	location   = "eu-central-1"
)

type Controller struct {
	pb.UnimplementedSyncServer

	storage records.BaseStorage

	bucketManager *minioS3.BucketManager
	objectManager *minioS3.ObjectManager
}

func NewController(storage records.BaseStorage, bucketManager *minioS3.BucketManager, objectManager *minioS3.ObjectManager) *Controller {
	return &Controller{
		storage:       storage,
		bucketManager: bucketManager,
		objectManager: objectManager,
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

		record := clientModels.ConvertStorageRecordFromGRPC(tmpReceive.GetRecord())
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

	userRecords, err := c.storage.GetRecords(stream.Context(), userID)
	for idx := range userRecords {
		err = stream.Send(&pb.PullResponse{Record: clientModels.ConvertStorageRecordToGRPC(&userRecords[idx])})
		if err != nil {
			return status.Errorf(codes.Internal, "send record: %v", err)
		}
	}

	return nil
}
func (c *Controller) PushBinary(stream pb.Sync_PushBinaryServer) error {
	userID, err := getUserID(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "extract userID from jwt: %v", err)
	}

	userBucket := strconv.FormatInt(userID, 10)
	if err = c.createBucketIfMissing(stream.Context(), userBucket); err != nil {
		return err
	}

	if err = c.objectManager.RemoveObjectsInBucket(stream.Context(), userBucket); err != nil {
		return status.Errorf(codes.Internal, "remove current user records: %v", err)
	}

	return c.saveBinaryObjects(userBucket, stream)
}

func (c *Controller) PullBinary(_ *emptypb.Empty, stream pb.Sync_PullBinaryServer) error {
	userID, err := getUserID(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "extract userID from jwt: %v", err)
	}

	userBucket := strconv.FormatInt(userID, 10)
	binaryStats, err := c.getUserBinaryObjectsList(userBucket, stream)
	if err != nil {
		return err
	}

	return c.sendUserBinaryObjects(userBucket, binaryStats, stream)
}

func (c *Controller) createBucketIfMissing(ctx context.Context, userBucket string) error {
	if ok, err := c.bucketManager.BucketExists(ctx, userBucket); err != nil {
		return status.Errorf(codes.Internal, "check user binary storage: %v", err)
	} else if !ok {
		if err = c.bucketManager.MakeBucket(ctx, userBucket, location); err != nil {
			return status.Errorf(codes.Internal, "make user bucket: %v", err)
		}
	}

	return nil
}

func (c *Controller) saveBinaryObjects(userBucket string, stream pb.Sync_PushBinaryServer) error {
	for {
		objectRaw, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}

		binary := objectRaw.GetBinary()
		object := &models.Object{
			Name:        binary.GetName(),
			Data:        binary.GetData(),
			Size:        int64(len(binary.GetData())),
			ContentType: minioS3.TypeBinary,
			Bucket:      userBucket,
		}

		if err = c.objectManager.PutObject(stream.Context(), object); err != nil {
			return status.Errorf(codes.Internal, "save object in storage: %v", err)
		}
	}
}

func (c *Controller) getUserBinaryObjectsList(userBucket string, stream pb.Sync_PullBinaryServer) (<-chan models.ObjectStat, error) {
	if ok, err := c.bucketManager.BucketExists(stream.Context(), userBucket); err != nil {
		return nil, status.Errorf(codes.Internal, "check user binary storage: %v", err)
	} else if !ok {
		return nil, nil
	}

	return c.objectManager.ListObjects(stream.Context(), userBucket), nil
}

func (c *Controller) sendUserBinaryObjects(userBucket string, objectsStat <-chan models.ObjectStat, stream pb.Sync_PullBinaryServer) error {
	for stat := range objectsStat {
		object, err := c.objectManager.GetObject(stream.Context(), &models.ObjectName{Name: stat.Key, Bucket: userBucket})
		if err != nil {
			return status.Errorf(codes.Internal, "extract user binary from storage: %v", err)
		}

		msgBinary := &pb.Binary{Name: object.Name, Data: object.Data}
		err = stream.Send(&pb.PullBinaryResponse{Binary: msgBinary})
		if err != nil {
			return status.Errorf(codes.Internal, "send message to stream: %v", err)
		}
	}

	return nil
}

func getUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1, fmt.Errorf("invalid context") // TODO: extract error.
	}

	rawUserID := md.Get(metaUserID)
	if len(rawUserID) == 0 || rawUserID[0] == "" {
		return -1, status.Errorf(codes.Unauthenticated, "missing user id") // TODO: extract error.
	}

	userID, err := strconv.ParseInt(fmt.Sprintf("%s", rawUserID[0]), 10, 64)
	if err != nil {
		return -1, status.Errorf(codes.InvalidArgument, "bad user id") // TODO: extract error.
	}

	return userID, nil
}
