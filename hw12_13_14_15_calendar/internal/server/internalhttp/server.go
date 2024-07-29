package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// _ "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api" нужен для инициализации документации Swagger.
	_ "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/api"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API Календаря
// @version 1.0
// @description Это простой API для управления событиями календаря.
// @host localhost:8080
// @BasePath /

type Server struct {
	httpServer          *http.Server
	eventService        services.EventService
	notificationService services.NotificationService
	healthService       services.HealthService
	logger              logger.Logger
}

func New(
	cfg config.HTTPServerConfig,
	logger logger.Logger,
	eventService services.EventService,
	notificationService services.NotificationService,
	healthService services.HealthService,
) *Server {
	router := mux.NewRouter()
	server := &Server{
		httpServer: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		eventService:        eventService,
		notificationService: notificationService,
		logger:              logger,
		healthService:       healthService,
	}

	// Роутинг для событий (events)
	router.HandleFunc("/events/day", server.listEventsForDateHandler).Methods("GET")
	router.HandleFunc("/events/week", server.listEventsForWeekHandler).Methods("GET")
	router.HandleFunc("/events/month", server.listEventsForMonthHandler).Methods("GET")
	router.HandleFunc("/events/{id}", server.updateEventHandler).Methods("PUT")
	router.HandleFunc("/events/{id}", server.deleteEventHandler).Methods("DELETE")
	router.HandleFunc("/events/{id}", server.getEventHandler).Methods("GET")
	router.HandleFunc("/events", server.createEventHandler).Methods("POST")
	router.HandleFunc("/events", server.listEventsHandler).Methods("GET")

	// Роутинг для уведомлений (notifications)
	router.HandleFunc("/notifications", server.createNotificationHandler).Methods("POST")
	router.HandleFunc("/notifications/{id}", server.updateNotificationHandler).Methods("PUT")
	router.HandleFunc("/notifications/{id}", server.deleteNotificationHandler).Methods("DELETE")
	router.HandleFunc("/notifications/{id}", server.getNotificationHandler).Methods("GET")
	router.HandleFunc("/notifications", server.listNotificationsHandler).Methods("GET")

	// Роутинг для healthcheck
	router.HandleFunc("/health", server.healthCheckHandler).Methods("GET")

	// Маршрут для Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Добавляем middleware
	router.Use(RequestIDMiddleware)
	router.Use(LoggingMiddleware(logger))

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("запуск http сервера")
	go func() {
		<-ctx.Done()
		err := s.Stop(ctx)
		if err != nil {
			s.logger.Errorf("ошибка остановки internalhttp сервера: %v", err)
		}
	}()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("ошибка остановки http.Server: %w", err)
	}
	return nil
}

func parseStartAndEndTime(r *http.Request) (time.Time, time.Time, error) {
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

// ErrorResponseWrapper используется для документации swagger.
type ErrorResponseWrapper struct {
	Errors    []string `json:"errors,omitempty"`
	Status    int      `json:"status"`
	RequestID string   `json:"requestId"`
}

// EventListResponseWrapper используется для документации swagger.
type EventListResponseWrapper struct {
	Data      []dto.EventData `json:"data"`
	Errors    []string        `json:"errors,omitempty"`
	Status    int             `json:"status"`
	RequestID string          `json:"requestId"`
}

// EventResponseWrapper используется для документации swagger.
type EventResponseWrapper struct {
	Data      dto.EventData `json:"data"`
	Errors    []string      `json:"errors,omitempty"`
	Status    int           `json:"status"`
	RequestID string        `json:"requestId"`
}

// NotificationListResponseWrapper используется для документации swagger.
type NotificationListResponseWrapper struct {
	Data      []dto.NotificationData `json:"data"`
	Errors    []string               `json:"errors,omitempty"`
	Status    int                    `json:"status"`
	RequestID string                 `json:"requestId"`
}

// NotificationResponseWrapper используется для документации swagger.
type NotificationResponseWrapper struct {
	Data      dto.NotificationData `json:"data"`
	Errors    []string             `json:"errors,omitempty"`
	Status    int                  `json:"status"`
	RequestID string               `json:"requestId"`
}

// @Summary Проверка состояния здоровья
// @Description Проверяет состояние сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 503 {object} ErrorResponseWrapper
// @Router /health [get].
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := s.healthService.HealthCheck(ctx)
	var response Response
	if err != nil {
		response = NewResponse(nil, []string{"Сервис недоступен"}, http.StatusServiceUnavailable)
	} else {
		response = NewResponse(map[string]string{"status": "ok"}, nil, http.StatusOK)
	}

	s.writeJSONResponse(w, r, response)
}

// @Summary Создать событие
// @Description Создает новое событие
// @Tags events
// @Accept json
// @Produce json
// @Param event body dto.EventData true "Запрос на создание события"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events [post].
func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventRequest dto.EventData

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	id, err := s.eventService.CreateEvent(r.Context(), eventRequest)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(map[string]interface{}{"id": id}, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Обновить событие
// @Description Обновляет существующее событие
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Param event body dto.EventData true "Запрос на обновление события"
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/{id} [put].
func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventRequest dto.EventData

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	err = s.eventService.UpdateEvent(r.Context(), id, eventRequest)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent)
	s.writeJSONResponse(w, r, response)
}

// @Summary Удалить событие
// @Description Удаляет существующее событие
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/{id} [delete].
func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	err = s.eventService.DeleteEvent(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent)
	s.writeJSONResponse(w, r, response)
}

// @Summary Получить событие
// @Description Получает событие по ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Success 200 {object} EventResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/{id} [get].
func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	event, err := s.eventService.GetEvent(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(event, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Список событий
// @Description Получает список событий между указанными датами
// @Tags events
// @Accept json
// @Produce json
// @Param startTime query string true "Время начала"
// @Param endTime query string true "Время окончания"
// @Success 200 {object} EventListResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events [get].
func (s *Server) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	start, end, err := parseStartAndEndTime(r)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	events, err := s.eventService.ListEvents(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(events, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Список событий на указанный день
// @Description Получает список событий на указанный день
// @Tags events
// @Accept json
// @Produce json
// @Param date query string true "Дата" format(date) example(2024-07-24)
// @Success 200 {object} EventListResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/day [get].
func (s *Server) listEventsForDateHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	start, err := time.Parse(time.DateOnly, date)
	if err != nil {
		response := NewResponse(nil, []string{"Некорректная дата"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	end := start.AddDate(0, 0, 1) // Добавляем 1 день к начальной дате для получения конца недели

	events, err := s.eventService.ListEvents(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(events, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Список событий на указанную неделю
// @Description Получает список событий на указанную неделю
// @Tags events
// @Accept json
// @Produce json
// @Param date query string true "Дата начала недели" format(date) example(2024-07-22)
// @Success 200 {object} EventListResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/week [get].
func (s *Server) listEventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	start, err := time.Parse(time.DateOnly, date)
	if err != nil {
		response := NewResponse(nil, []string{"Некорректная дата"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	end := start.AddDate(0, 0, 7) // Добавляем 7 дней к начальной дате для получения конца недели

	events, err := s.eventService.ListEvents(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(events, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Список событий на указанный месяц
// @Description Получает список событий на указанный месяц
// @Tags events
// @Accept json
// @Produce json
// @Param date query string true "Дата начала месяца" format(date) example(2024-07-01)
// @Success 200 {object} EventListResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /events/month [get].
func (s *Server) listEventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	start, err := time.Parse(time.DateOnly, date)
	if err != nil {
		response := NewResponse(nil, []string{"Некорректная дата"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	end := start.AddDate(0, 1, 0) // Добавляем 1 месяц к начальной дате

	events, err := s.eventService.ListEvents(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(events, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, r *http.Request, response Response) {
	requestID := r.Context().Value(requestIDKey).(string)
	response.RequestID = requestID

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.logger.Errorf("on marshal response %v: %v", response, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	_, err = w.Write(jsonResponse)
	if err != nil {
		s.logger.Errorf("on writeJSONResponse: %v", err)
		return
	}
}

// @Summary Создать уведомление
// @Description Создает новое уведомление
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body dto.NotificationData true "Запрос на создание уведомления"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /notifications [post].
func (s *Server) createNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notificationRequest dto.NotificationData

	if err := json.NewDecoder(r.Body).Decode(&notificationRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	notificationRequest.Sent = false

	id, err := s.notificationService.CreateNotification(r.Context(), notificationRequest)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(map[string]interface{}{"id": id}, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Обновить уведомление
// @Description Обновляет существующее уведомление
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "ID уведомления"
// @Param notification body dto.NotificationData true "Запрос на обновление уведомления"
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /notifications/{id} [put].
func (s *Server) updateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notificationRequest dto.NotificationData

	if err := json.NewDecoder(r.Body).Decode(&notificationRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response := NewResponse(nil, []string{"Invalid ID format"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	notificationRequest.ID = id
	notificationRequest.Sent = false

	err = s.notificationService.UpdateNotification(r.Context(), id, notificationRequest)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent)
	s.writeJSONResponse(w, r, response)
}

// @Summary Удалить уведомление
// @Description Удаляет существующее уведомление
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "ID уведомления"
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /notifications/{id} [delete].
func (s *Server) deleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response := NewResponse(nil, []string{"Invalid ID format"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	err = s.notificationService.DeleteNotification(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent)
	s.writeJSONResponse(w, r, response)
}

// @Summary Получить уведомление
// @Description Получает уведомление по ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "ID уведомления"
// @Success 200 {object} NotificationResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /notifications/{id} [get].
func (s *Server) getNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response := NewResponse(nil, []string{"Invalid ID format"}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	notification, err := s.notificationService.GetNotification(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(notification, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}

// @Summary Список уведомлений
// @Description Получает список уведомлений между указанными датами
// @Tags notifications
// @Accept json
// @Produce json
// @Param start_time query string true "Время начала"
// @Param end_time query string true "Время окончания"
// @Success 200 {object} NotificationListResponseWrapper
// @Failure 400 {object} ErrorResponseWrapper
// @Failure 500 {object} ErrorResponseWrapper
// @Router /notifications [get].
func (s *Server) listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	start, end, err := parseStartAndEndTime(r)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest)
		s.writeJSONResponse(w, r, response)
		return
	}

	notifications, err := s.notificationService.ListNotifications(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError)
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(notifications, nil, http.StatusOK)
	s.writeJSONResponse(w, r, response)
}
