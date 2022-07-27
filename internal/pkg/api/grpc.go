package api

import (
	"fmt"
	"net"

	cfg "github.com/digiexpress/dlocator/internal/pkg/config"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DLocatorGrpcServer struct {
	listener net.Listener
	server   *grpc.Server
}

func NewDLocatorGrpcServer(courierLocator CourierLocatorServer, config *cfg.AppConfig) (*DLocatorGrpcServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Grpc.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", config.Grpc.Port, err)
	}

	logEntry := logrus.WithFields(map[string]interface{}{"app": config.Grpc.AppName})
	interceptors := []grpc.UnaryServerInterceptor{
		grpclogrus.UnaryServerInterceptor(logEntry),
		grpcrecovery.UnaryServerInterceptor(),
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	RegisterCourierLocatorServer(grpcServer, courierLocator)

	return &DLocatorGrpcServer{listener: listener, server: grpcServer}, nil
}

func (s *DLocatorGrpcServer) Serve() error {
	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}

func (s *DLocatorGrpcServer) Stop() {
	s.server.GracefulStop()
}
