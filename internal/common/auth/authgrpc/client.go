package authgrpc

import (
	"context"
	"strings"

	"github.com/erupshis/key_keeper/internal/common/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor struct {
	token string
}

func NewClientInterceptor() *ClientInterceptor {
	return &ClientInterceptor{}
}

func (c *ClientInterceptor) StreamClient() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		procName := method[strings.LastIndex(method, ".")+1:]
		_, methodOk := procExclusions[procName]
		if !methodOk {
			ctx = c.addTokenInHeader(ctx)
		}

		var header metadata.MD
		opts = append(opts, grpc.Header(&header))

		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return s, err
		}

		if methodOk {
			c.extractTokenFromHeader(&header)
		}

		return s, nil
	}
}

func (c *ClientInterceptor) UnaryClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		procName := method[strings.LastIndex(method, ".")+1:]
		_, methodOk := procExclusions[procName]
		if !methodOk {
			ctx = c.addTokenInHeader(ctx)
		}

		var header metadata.MD
		opts = append(opts, grpc.Header(&header))

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}

		if methodOk {
			c.extractTokenFromHeader(&header)
		}

		return nil
	}
}

func (c *ClientInterceptor) addTokenInHeader(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{auth.TokenHeader: c.token})
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *ClientInterceptor) extractTokenFromHeader(header *metadata.MD) {
	authHeader := (*header)["authorization"]
	c.token = authHeader[0]
}
