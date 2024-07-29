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

func TestGRPCServer(t *testing.T) {
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

	client := api.NewEventServiceClient(conn)

	t.Run("CreateEvent", func(t *testing.T) {
		req := &api.CreateEventRequest{
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   timestamppb.Now(),
			EndTime:     timestamppb.Now(),
			UserId:      uuid.New().String(),
		}

		resp, err := client.CreateEvent(context.Background(), req)
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

		createResp, err := client.CreateEvent(context.Background(), createReq)
		require.NoError(t, err)
		require.NotNil(t, createResp)
		require.NotEmpty(t, createResp.Id)

		req := &api.GetEventRequest{Id: createResp.Id}

		resp, err := client.GetEvent(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "Test Event", resp.GetEvent().Title)
	})

	grpcServer.Stop(context.Background())
}
