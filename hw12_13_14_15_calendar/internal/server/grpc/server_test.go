package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	memorystorage "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGRPCServer(t *testing.T) { //nolint:funlen
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := lis.Addr().String()
	err = lis.Close()
	require.NoError(t, err)

	store := memorystorage.New()
	eventService := services.NewEventService(store)
	notificationService := services.NewNotificationService(store)
	healthService := services.NewHealthService(store)

	logInstance, err := logger.New(config.LoggerConfig{
		Level:            "fatal",
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	require.NoError(t, err)

	grpcServer, err := New(
		eventService,
		notificationService,
		healthService,
		logInstance,
		config.GRPCServerConfig{Address: addr},
	)
	require.NoError(t, err)

	errChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Start(context.Background()); err != nil {
			errChan <- err
		}
	}()

	// Ждем, пока сервер полностью запустится
	time.Sleep(2 * time.Second)

	select {
	case err := <-errChan:
		require.NoError(t, err)
	default:
	}

	t.Logf("Listening on: %s\n", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	eventClient := api.NewEventServiceClient(conn)
	notificationClient := api.NewNotificationServiceClient(conn)

	t.Run("Events", func(t *testing.T) {
		t.Run("CreateEvent", func(t *testing.T) {
			req := &api.CreateEventRequest{
				Title:       "Test Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.Now(),
				UserId:      uuid.New().String(),
			}

			resp, err := eventClient.CreateEvent(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.Id)
		})

		t.Run("GetEvent", func(t *testing.T) {
			createReq := &api.CreateEventRequest{
				Title:       "Test Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.Now(),
				UserId:      uuid.New().String(),
			}

			createResp, err := eventClient.CreateEvent(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			req := &api.GetEventRequest{Id: createResp.Id}

			resp, err := eventClient.GetEvent(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, "Test Event", resp.GetEvent().Title)
		})

		t.Run("UpdateEvent", func(t *testing.T) {
			createReq := &api.CreateEventRequest{
				Title:       "Test Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.Now(),
				UserId:      uuid.New().String(),
			}

			createResp, err := eventClient.CreateEvent(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			updateReq := &api.UpdateEventRequest{
				Id:          createResp.Id,
				Title:       "Updated Event",
				Description: "Updated Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.Now(),
				UserId:      createReq.UserId,
			}

			updateResp, err := eventClient.UpdateEvent(context.Background(), updateReq)
			require.NoError(t, err)
			require.NotNil(t, updateResp)
		})

		t.Run("DeleteEvent", func(t *testing.T) {
			createReq := &api.CreateEventRequest{
				Title:       "Test Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.Now(),
				UserId:      uuid.New().String(),
			}

			createResp, err := eventClient.CreateEvent(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			req := &api.DeleteEventRequest{Id: createResp.Id}

			resp, err := eventClient.DeleteEvent(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
		})

		t.Run("ListEvents", func(t *testing.T) {
			startTime := timestamppb.Now()
			endTime := timestamppb.New(time.Now().Add(24 * time.Hour))

			req := &api.ListEventsRequest{
				StartTime: startTime,
				EndTime:   endTime,
			}

			resp, err := eventClient.ListEvents(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	})

	t.Run("Notifications", func(t *testing.T) {
		t.Run("CreateNotification", func(t *testing.T) {
			req := &api.CreateNotificationRequest{
				EventId: uuid.New().String(),
				UserId:  uuid.New().String(),
				Time:    timestamppb.Now(),
				Message: "Test Notification",
				Sent:    false,
			}

			resp, err := notificationClient.CreateNotification(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.Id)
		})

		t.Run("GetNotification", func(t *testing.T) {
			createReq := &api.CreateNotificationRequest{
				EventId: uuid.New().String(),
				UserId:  uuid.New().String(),
				Time:    timestamppb.Now(),
				Message: "Test Notification",
				Sent:    false,
			}

			createResp, err := notificationClient.CreateNotification(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			req := &api.GetNotificationRequest{Id: createResp.Id}

			resp, err := notificationClient.GetNotification(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, "Test Notification", resp.GetNotification().Message)
		})

		t.Run("UpdateNotification", func(t *testing.T) {
			createReq := &api.CreateNotificationRequest{
				EventId: uuid.New().String(),
				UserId:  uuid.New().String(),
				Time:    timestamppb.Now(),
				Message: "Test Notification",
				Sent:    false,
			}

			createResp, err := notificationClient.CreateNotification(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			updateReq := &api.UpdateNotificationRequest{
				Id:      createResp.Id,
				EventId: createReq.EventId,
				UserId:  createReq.UserId,
				Time:    timestamppb.Now(),
				Message: "Updated Notification",
				Sent:    true,
			}

			updateResp, err := notificationClient.UpdateNotification(context.Background(), updateReq)
			require.NoError(t, err)
			require.NotNil(t, updateResp)
		})

		t.Run("DeleteNotification", func(t *testing.T) {
			createReq := &api.CreateNotificationRequest{
				EventId: uuid.New().String(),
				UserId:  uuid.New().String(),
				Time:    timestamppb.Now(),
				Message: "Test Notification",
				Sent:    false,
			}

			createResp, err := notificationClient.CreateNotification(context.Background(), createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotEmpty(t, createResp.Id)

			req := &api.DeleteNotificationRequest{Id: createResp.Id}

			resp, err := notificationClient.DeleteNotification(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
		})

		t.Run("ListNotifications", func(t *testing.T) {
			startTime := timestamppb.Now()
			endTime := timestamppb.New(time.Now().Add(24 * time.Hour))

			req := &api.ListNotificationsRequest{
				StartTime: startTime,
				EndTime:   endTime,
			}

			resp, err := notificationClient.ListNotifications(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	})

	grpcServer.Stop(context.Background())
}
