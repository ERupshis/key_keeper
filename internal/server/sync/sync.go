package sync

import (
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	pb.UnimplementedSyncServer
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Push(stream pb.Sync_PushServer) error {
	_ = stream.SendAndClose(&emptypb.Empty{})
	return status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (c *Controller) Pull(_ *emptypb.Empty, _ pb.Sync_PullServer) error {
	return status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (c *Controller) PushBinary(_ pb.Sync_PushBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method PushBinary not implemented")
}
func (c *Controller) PullBinary(_ *emptypb.Empty, _ pb.Sync_PullBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method PullBinary not implemented")
}
