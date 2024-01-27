package authgrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/erupshis/key_keeper/internal/common/auth"
	"github.com/erupshis/key_keeper/internal/common/jwtgenerator"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Authorize(ctx context.Context, jwt *jwtgenerator.JwtGenerator) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders, ok := md[auth.TokenHeader]
	if !ok || len(authHeaders) == 0 {
		return status.Error(codes.Unauthenticated, "missing token in metadata")
	}

	token := strings.TrimSpace(authHeaders[0])

	userID, err := jwt.GetUserID(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "jwt validation: %v", err)
	}

	mdPairs := metadata.Pairs(
		auth.UserID, strconv.FormatInt(userID, 10),
	)

	ctx = metadata.NewIncomingContext(ctx, mdPairs)
	return nil
}

func StreamServer(jwt *jwtgenerator.JwtGenerator, logger logger.BaseLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := Authorize(ss.Context(), jwt); err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func UnaryServer(jwt *jwtgenerator.JwtGenerator, logger logger.BaseLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := Authorize(ctx, jwt); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
