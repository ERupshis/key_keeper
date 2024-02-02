package authgrpc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/erupshis/key_keeper/internal/common/auth"
	"github.com/erupshis/key_keeper/internal/common/jwtgenerator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	Login    = "Auth/Login"
	Register = "Auth/Register"
)

var (
	procExclusions = map[string]struct{}{
		Login:    struct{}{},
		Register: struct{}{},
	}
)

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *wrappedStream) Context() context.Context {
	return s.ctx
}

func Authorize(ctx context.Context, jwt *jwtgenerator.JwtGenerator) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader, ok := md[auth.TokenHeader]
	if !ok || len(authHeader) == 0 {
		return -1, status.Error(codes.Unauthenticated, "missing token in metadata")
	}

	token := strings.Split(authHeader[0], " ")
	if len(token) != 2 || token[0] != auth.TokenType {
		return -1, status.Error(codes.InvalidArgument, "incorrect authorization data")
	}

	userID, err := jwt.GetUserID(token[1])
	if err != nil {
		return -1, status.Errorf(codes.Unauthenticated, "jwt validation: %v", err)
	}

	return userID, nil
}

func StreamServer(jwt *jwtgenerator.JwtGenerator) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		procName := info.FullMethod[strings.LastIndex(info.FullMethod, ".")+1:]
		if _, ok := procExclusions[procName]; !ok {
			userID, err := Authorize(ss.Context(), jwt)
			if err != nil {
				return err
			}

			md, okMD := metadata.FromIncomingContext(ctx)
			if !okMD {
				return fmt.Errorf("read headers failed")
			}

			md.Append(auth.UserID, strconv.FormatInt(userID, 10))
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		return handler(srv, &wrappedStream{ss, ctx})
	}
}

func UnaryServer(jwt *jwtgenerator.JwtGenerator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		procName := info.FullMethod[strings.LastIndex(info.FullMethod, ".")+1:]
		if _, ok := procExclusions[procName]; !ok {
			userID, err := Authorize(ctx, jwt)
			if err != nil {
				return nil, err
			}

			md, okMD := metadata.FromIncomingContext(ctx)
			if !okMD {
				return nil, fmt.Errorf("read headers failed")
			}

			md.Append(auth.UserID, strconv.FormatInt(userID, 10))
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}
