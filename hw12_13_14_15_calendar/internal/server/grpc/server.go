package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	api.UnimplementedEventServiceServer
	api.UnimplementedNotificationServiceServer
	grpcServer          *grpc.Server
	config              config.GRPCServerConfig
	eventService        services.EventService
	notificationService services.NotificationService
	healthService       services.HealthService
	logger              logger.Logger
}

func New(
	eventService services.EventService,
	notificationService services.NotificationService,
	healthService services.HealthService,
	logger logger.Logger,
	config config.GRPCServerConfig,
) (*Server, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor(logger)),
	)

	server := &Server{
		eventService:        eventService,
		notificationService: notificationService,
		healthService:       healthService,
		logger:              logger,
		config:              config,
		grpcServer:          grpcServer,
	}

	return server, nil
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("on net.Listen: %w", err)
	}

	// Register gRPC services
	api.RegisterEventServiceServer(s.grpcServer, s)
	api.RegisterNotificationServiceServer(s.grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	s.logger.Infof("gRPC server started. Listening on %s", s.config.Address)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(_ context.Context) {
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC server stopped gracefully")
}

func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	event := dto.EventData{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		StartTime:   req.GetStartTime().AsTime(),
		EndTime:     req.GetEndTime().AsTime(),
		UserID:      uuid.MustParse(req.GetUserId()),
	}
	id, err := s.eventService.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &api.CreateEventResponse{Id: id.String()}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
	event := dto.EventData{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		StartTime:   req.GetStartTime().AsTime(),
		EndTime:     req.GetEndTime().AsTime(),
		UserID:      uuid.MustParse(req.GetUserId()),
	}
	err := s.eventService.UpdateEvent(ctx, uuid.MustParse(req.GetId()), event)
	if err != nil {
		return nil, err
	}
	return &api.UpdateEventResponse{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	err := s.eventService.DeleteEvent(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &api.DeleteEventResponse{}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
	event, err := s.eventService.GetEvent(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &api.GetEventResponse{Event: dto.ToAPIEvent(event)}, nil
}

func (s *Server) ListEvents(ctx context.Context, req *api.ListEventsRequest) (*api.ListEventsResponse, error) {
	start := req.GetStartTime().AsTime()
	end := req.GetEndTime().AsTime()
	events, err := s.eventService.ListEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	apiEvents := make([]*api.Event, len(events))
	for i, event := range events {
		apiEvents[i] = dto.ToAPIEvent(event)
	}
	return &api.ListEventsResponse{Events: apiEvents}, nil
}

func (s *Server) ListEventsForDate(
	ctx context.Context,
	req *api.ListEventsForDateRequest,
) (*api.ListEventsResponse, error) {
	start := req.GetDate().AsTime()
	end := start.AddDate(0, 0, 1) // Добавляем 1 день

	events, err := s.eventService.ListEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	apiEvents := make([]*api.Event, len(events))
	for i, event := range events {
		apiEvents[i] = dto.ToAPIEvent(event)
	}
	return &api.ListEventsResponse{Events: apiEvents}, nil
}

func (s *Server) ListEventsForWeek(
	ctx context.Context,
	req *api.ListEventsForWeekRequest,
) (*api.ListEventsResponse, error) {
	start := req.GetDate().AsTime()
	end := start.AddDate(0, 0, 7) // Добавляем 7 дней

	events, err := s.eventService.ListEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	apiEvents := make([]*api.Event, len(events))
	for i, event := range events {
		apiEvents[i] = dto.ToAPIEvent(event)
	}
	return &api.ListEventsResponse{Events: apiEvents}, nil
}

func (s *Server) ListEventsForMonth(
	ctx context.Context,
	req *api.ListEventsForMonthRequest,
) (*api.ListEventsResponse, error) {
	start := req.GetDate().AsTime()
	end := start.AddDate(0, 1, 0) // Добавляем 1 месяц

	events, err := s.eventService.ListEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	apiEvents := make([]*api.Event, len(events))
	for i, event := range events {
		apiEvents[i] = dto.ToAPIEvent(event)
	}
	return &api.ListEventsResponse{Events: apiEvents}, nil
}

func (s *Server) CreateNotification(
	ctx context.Context,
	req *api.CreateNotificationRequest,
) (*api.CreateNotificationResponse, error) {
	notification := dto.NotificationData{
		EventID: uuid.MustParse(req.GetEventId()),
		UserID:  uuid.MustParse(req.GetUserId()),
		Time:    req.GetTime().AsTime(),
		Message: req.GetMessage(),
		Sent:    req.GetSent(),
	}
	id, err := s.notificationService.CreateNotification(ctx, notification)
	if err != nil {
		return nil, err
	}
	return &api.CreateNotificationResponse{Id: id.String()}, nil
}

func (s *Server) UpdateNotification(
	ctx context.Context,
	req *api.UpdateNotificationRequest,
) (*api.UpdateNotificationResponse, error) {
	notification := dto.NotificationData{
		EventID: uuid.MustParse(req.GetEventId()),
		UserID:  uuid.MustParse(req.GetUserId()),
		Time:    req.GetTime().AsTime(),
		Message: req.GetMessage(),
		Sent:    req.GetSent(),
	}
	err := s.notificationService.UpdateNotification(ctx, uuid.MustParse(req.GetId()), notification)
	if err != nil {
		return nil, err
	}
	return &api.UpdateNotificationResponse{}, nil
}

func (s *Server) DeleteNotification(
	ctx context.Context,
	req *api.DeleteNotificationRequest,
) (*api.DeleteNotificationResponse, error) {
	err := s.notificationService.DeleteNotification(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &api.DeleteNotificationResponse{}, nil
}

func (s *Server) GetNotification(
	ctx context.Context,
	req *api.GetNotificationRequest,
) (*api.GetNotificationResponse, error) {
	notification, err := s.notificationService.GetNotification(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &api.GetNotificationResponse{Notification: dto.ToAPINotification(notification)}, nil
}

func (s *Server) ListNotifications(
	ctx context.Context,
	req *api.ListNotificationsRequest,
) (*api.ListNotificationsResponse, error) {
	start := req.GetStartTime().AsTime()
	end := req.GetEndTime().AsTime()
	notifications, err := s.notificationService.ListNotifications(ctx, start, end)
	if err != nil {
		return nil, err
	}
	apiNotifications := make([]*api.Notification, len(notifications))
	for i, notification := range notifications {
		apiNotifications[i] = dto.ToAPINotification(notification)
	}
	return &api.ListNotificationsResponse{Notifications: apiNotifications}, nil
}
