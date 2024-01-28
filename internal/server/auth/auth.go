package auth

import (
	"context"
	"errors"

	"github.com/erupshis/key_keeper/internal/common/auth"
	"github.com/erupshis/key_keeper/internal/common/auth/models"
	"github.com/erupshis/key_keeper/internal/common/auth/storage"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	pb.UnimplementedAuthServer

	authManager *auth.Manager
}

func NewController(authManager *auth.Manager) *Controller {
	return &Controller{
		authManager: authManager,
	}
}

func (c *Controller) Login(ctx context.Context, in *pb.LoginRequest) (*emptypb.Empty, error) {
	token, err := c.authManager.Login(ctx, models.ConvertUserFromGRPC(in.GetCreds()))
	if err != nil {
		if errors.Is(err, auth.ErrMismatchPassword) || errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		}

		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	md := metadata.Pairs(auth.TokenHeader, token)
	if err = grpc.SendHeader(ctx, md); err != nil {
		return nil, status.Errorf(codes.Internal, "add token in header: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *Controller) Register(ctx context.Context, in *pb.RegisterRequest) (*emptypb.Empty, error) {
	if err := c.authManager.Register(ctx, models.ConvertUserFromGRPC(in.GetCreds())); err != nil {
		if errors.Is(err, auth.ErrLoginOccupied) {
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		}

		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &emptypb.Empty{}, nil
}
