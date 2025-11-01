package grpc

import (
	"context"
	"mockium/internal/adapters/common"
	"mockium/pkg/ports"
	"net"
	"time"

	"google.golang.org/grpc"
)

func NewServer(logger ports.Logger, address string, handlers ...GRPCHandler[GRPCRequest]) ports.Server {
	return &server{
		logger:   logger,
		address:  address,
		server:   grpc.NewServer(),
		handlers: handlers,
	}
}

type server struct {
	started  bool
	address  string
	logger   ports.Logger
	server   *grpc.Server
	handlers []GRPCHandler[GRPCRequest]
}

func (inst *server) Start() error {
	lis, err := net.Listen("tcp", inst.address)
	if err != nil {
		return err
	}

	go func() {
		inst.started = true
		inst.logger.Info("Starting gRPC server", inst.address)
		if err := inst.server.Serve(lis); err != nil {
			inst.logger.Error("gRPC server error", err)
			inst.started = false
		}
	}()

	time.Sleep(100 * time.Millisecond) // Give the server a moment to start

	if !inst.started {
		return common.ErrorServerNotRunning
	}

	return nil
}

func (inst *server) Stop() {
	if inst.server != nil || !inst.started {
		inst.logger.Info("Stopping gRPC server")
		inst.server.GracefulStop()
	}
}

func (inst *server) Restart() error {
	if inst.started {
		inst.Stop()
	}

	return inst.Start()
}

func (inst *server) IsRunning() bool { return inst.started }

func (inst *server) Configure() error {
	if inst.started {
		inst.logger.Warn("Reconfiguring a running server. Changes will take effect after restart.")
	}

	for _, handler := range inst.handlers {
		h := handler
		methods := make([]grpc.MethodDesc, 0)
		streamMethods := make([]grpc.StreamDesc, 0)

		for _, method := range handler.UnaryMethods() {
			m := method
			methods = append(methods,
				grpc.MethodDesc{
					MethodName: method,
					Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
						req := NewRequest(ctx, m, "unary", nil, nil)
						if err := dec(&req); err != nil {
							return nil, err
						}

						resp, err := h.Handle(req)
						if err != nil {
							return nil, err
						}

						return resp, nil
					},
				},
			)
		}

		for _, ms := range handler.StreamMethods() {
			streamMethods = append(streamMethods,
				grpc.StreamDesc{
					StreamName: ms,
					Handler: func(srv any, stream grpc.ServerStream) error {
						return nil
					},
					ServerStreams: true,
					ClientStreams: true,
				},
			)
		}

		inst.server.RegisterService(
			&grpc.ServiceDesc{
				ServiceName: handler.Service(),
				HandlerType: (*GRPCHandler[GRPCRequest])(nil),
				Methods:     methods,
				Streams:     streamMethods,
			},
			handler,
		)

	}

	return nil
}
