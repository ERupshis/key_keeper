package logger

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StreamServer(logger BaseLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Infof("Stream method %s called", info.FullMethod)

		err := handler(srv, ss)
		if err != nil {
			logger.Infof("authgrpc stream: %v", err)
		}

		return err
	}
}

func UnaryServer(logger BaseLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("Unary method %s called", info.FullMethod)

		resp, err := handler(ctx, req)

		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				logger.Infof("Unary method '%s' completed with error '%v', status: %s", info.FullMethod, err, st.Code().String())
			} else {
				logger.Infof("Unary method '%s' completed with error '%v', status: unknown", info.FullMethod, err)
			}
		} else {
			logger.Infof("Unary method '%s' completed, status: %s", info.FullMethod, codes.OK.String())
		}

		return resp, err
	}
}

func StreamClient(logger BaseLogger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		logger.Infof("Stream method %s called", method)

		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			logger.Infof("Stream method %s result with err: %v", method, err)
		} else {
			logger.Infof("Stream method %s successfully initiated", method)
		}

		return s, err
	}
}

func UnaryClient(logger BaseLogger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logger.Infof("Unary method %s called", method)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			logger.Infof("Unary method %s result with err: %v", method, err)
		} else {
			logger.Infof("Unary method %s successfully completed", method)
		}

		return err
	}
}
