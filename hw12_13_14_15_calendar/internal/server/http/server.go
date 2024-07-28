package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// _ "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/docs" нужен для инициализации документации Swagger.
	_ "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/docs"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
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
	cfg config.ServerConfig,
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

	// Добавляем маршруты и обработчики
	router.HandleFunc("/events", server.createEventHandler).Methods("POST")
	router.HandleFunc("/events/{id}", server.updateEventHandler).Methods("PUT")
	router.HandleFunc("/events/{id}", server.deleteEventHandler).Methods("DELETE")
	router.HandleFunc("/events/{id}", server.getEventHandler).Methods("GET")
	router.HandleFunc("/events", server.listEventsHandler).Methods("GET")
	router.HandleFunc("/health", server.healthCheckHandler).Methods("GET")

	// Маршрут для Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Добавляем middleware
	router.Use(RequestIDMiddleware)
	router.Use(LoggingMiddleware(logger))

	return server
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// @Summary Проверка состояния здоровья
// @Description Проверяет состояние сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health [get].
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := s.healthService.HealthCheck(ctx)
	var response Response
	if err != nil {
		response = NewResponse(nil, []string{"Сервис недоступен"}, http.StatusServiceUnavailable, "")
	} else {
		response = NewResponse(map[string]string{"status": "ok"}, nil, http.StatusOK, "")
	}

	s.writeJSONResponse(w, r, response)
}

// eventRequest представляет запрос на создание/обновление события.
type eventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uuid.UUID `json:"userId"`
}

// @Summary Создать событие
// @Description Создает новое событие
// @Tags events
// @Accept json
// @Produce json
// @Param event body eventRequest true "Запрос на создание события"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /events [post].
func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventRequest eventRequest

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	id, err := s.eventService.CreateEvent(
		r.Context(),
		eventRequest.Title,
		eventRequest.Description,
		eventRequest.StartTime,
		eventRequest.EndTime,
		eventRequest.UserID,
	)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(map[string]interface{}{"id": id}, nil, http.StatusOK, "")
	s.writeJSONResponse(w, r, response)
}

// @Summary Обновить событие
// @Description Обновляет существующее событие
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Param event body eventRequest true "Запрос на обновление события"
// @Success 204 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /events/{id} [put].
func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventRequest eventRequest

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	err = s.eventService.UpdateEvent(
		r.Context(),
		id,
		eventRequest.Title,
		eventRequest.Description,
		eventRequest.StartTime,
		eventRequest.EndTime,
		eventRequest.UserID,
	)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent, "")
	s.writeJSONResponse(w, r, response)
}

// @Summary Удалить событие
// @Description Удаляет существующее событие
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Success 204 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /events/{id} [delete].
func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	err = s.eventService.DeleteEvent(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(nil, nil, http.StatusNoContent, "")
	s.writeJSONResponse(w, r, response)
}

// @Summary Получить событие
// @Description Получает событие по ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID события"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /events/{id} [get].
func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	event, err := s.eventService.GetEvent(r.Context(), id)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(event, nil, http.StatusOK, "")
	s.writeJSONResponse(w, r, response)
}

// @Summary Список событий
// @Description Получает список событий между указанными датами
// @Tags events
// @Accept json
// @Produce json
// @Param startTime query string true "Время начала"
// @Param endTime query string true "Время окончания"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /events [get].
func (s *Server) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		response := NewResponse(nil, []string{"Некорректное время начала"}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		response := NewResponse(nil, []string{"Некорректное время окончания"}, http.StatusBadRequest, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	events, err := s.eventService.ListEvents(r.Context(), start, end)
	if err != nil {
		response := NewResponse(nil, []string{err.Error()}, http.StatusInternalServerError, "")
		s.writeJSONResponse(w, r, response)
		return
	}

	response := NewResponse(events, nil, http.StatusOK, "")
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
