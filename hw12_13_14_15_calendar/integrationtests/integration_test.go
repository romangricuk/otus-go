package integrationtests

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
)

var client api.EventServiceClient
var conn *grpc.ClientConn

func TestMain(m *testing.M) {
	var err error
	grpcAddress := getGRPCAddress()

	fmt.Printf("grpcAddress: %s\n", grpcAddress)

	conn, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = api.NewEventServiceClient(conn)
	code := m.Run()
	// Очистка после выполнения тестов
	conn.Close()
	os.Exit(code)
}

func getGRPCAddress() string {
	if addr := os.Getenv("GRPC_ADDRESS"); addr != "" {
		return addr
	}
	return "localhost:9090"
}

// TestCreateEvent Тест на добавление события
func TestCreateEvent(t *testing.T) {
	ctx := context.Background()
	startTime := time.Now().Add(1 * time.Hour)
	endTime := time.Now().Add(2 * time.Hour)

	// Успешное создание события
	resp, err := client.CreateEvent(ctx, &api.CreateEventRequest{
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   timestamppb.New(startTime),
		EndTime:     timestamppb.New(endTime),
		UserId:      uuid.New().String(),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	eventID := resp.Id

	// Очистка после теста
	t.Cleanup(func() {
		_, err := client.DeleteEvent(ctx, &api.DeleteEventRequest{Id: eventID})
		assert.NoError(t, err)
	})
}

// TestListEvents Тест на получение списка событий за день/неделю/месяц
func TestListEvents(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := uuid.New()

	// Создаем несколько событий
	for i := 0; i < 3; i++ {
		startTime := now.AddDate(0, 0, i)
		endTime := startTime.Add(1 * time.Hour)
		resp, err := client.CreateEvent(ctx, &api.CreateEventRequest{
			Title:       "Event " + strconv.Itoa(i),
			Description: "Event description",
			StartTime:   timestamppb.New(startTime),
			EndTime:     timestamppb.New(endTime),
			UserId:      userID.String(),
		})
		assert.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.DeleteEvent(ctx, &api.DeleteEventRequest{Id: resp.Id})
			assert.NoError(t, err)
		})
	}

	// Получение событий за день
	date := timestamppb.New(now)
	dayResp, err := client.ListEventsForDate(ctx, &api.ListEventsForDateRequest{
		Date: date,
	})
	assert.NoError(t, err)
	assert.NotNil(t, dayResp)
	assert.GreaterOrEqual(t, len(dayResp.Events), 1)

	// Получение событий за неделю
	weekResp, err := client.ListEventsForWeek(ctx, &api.ListEventsForWeekRequest{
		Date: date,
	})
	assert.NoError(t, err)
	assert.NotNil(t, weekResp)
	assert.GreaterOrEqual(t, len(weekResp.Events), 3)

	// Получение событий за месяц
	monthResp, err := client.ListEventsForMonth(ctx, &api.ListEventsForMonthRequest{
		Date: date,
	})
	assert.NoError(t, err)
	assert.NotNil(t, monthResp)
	assert.GreaterOrEqual(t, len(monthResp.Events), 3)
}

// TestCreateNotification Тест на отправку уведомлений
func TestCreateNotification(t *testing.T) {
	ctx := context.Background()

	startTime := time.Now().Add(1 * time.Hour)
	endTime := time.Now().Add(2 * time.Hour)

	// Создание события
	resp, err := client.CreateEvent(ctx, &api.CreateEventRequest{
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   timestamppb.New(startTime),
		EndTime:     timestamppb.New(endTime),
		UserId:      uuid.New().String(),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	eventID := resp.Id

	// Очистка после теста
	t.Cleanup(func() {
		_, err := client.DeleteEvent(ctx, &api.DeleteEventRequest{Id: eventID})
		assert.NoError(t, err)
	})

	//
	notifClient := api.NewNotificationServiceClient(conn)
	notifTime := time.Now().Add(30 * time.Minute)

	// Создаем уведомление
	notifyResp, err := notifClient.CreateNotification(ctx, &api.CreateNotificationRequest{
		EventId: eventID,
		UserId:  uuid.New().String(),
		Time:    timestamppb.New(notifTime),
		Message: "Reminder for your event",
		Sent:    dto.NotificationOnWait,
	})
	require.NoError(t, err)
	require.NotNil(t, notifyResp)
	notifID := notifyResp.Id

	// Очистка после теста
	t.Cleanup(func() {
		_, err := notifClient.DeleteNotification(ctx, &api.DeleteNotificationRequest{Id: notifID})
		assert.NoError(t, err)
	})
}
