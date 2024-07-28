package grpc

import (
	"context"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	grpcServer          *grpc.Server
	eventService        services.EventService
	notificationService services.NotificationService
	healthService       services.HealthService
	logger              logger.Logger
	listener            net.Listener
}

func New(
	eventService services.EventService,
	notificationService services.NotificationService,
	healthService services.HealthService,
	logger logger.Logger,
	address string,
) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s := &Server{
		eventService:        eventService,
		notificationService: notificationService,
		healthService:       healthService,
		logger:              logger,
		listener:            listener,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor(logger)),
	)

	s.grpcServer = grpcServer

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("gRPC сервер запущен")

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC сервер остановлен")
}
