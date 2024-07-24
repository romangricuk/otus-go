package services

import (
	"context"
	"fmt"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type HealthServiceImpl struct {
	storage storage.Storage
}

func NewHealthService(store storage.Storage) *HealthServiceImpl {
	return &HealthServiceImpl{storage: store}
}

func (s HealthServiceImpl) HealthCheck(ctx context.Context) error {
	err := s.storage.HealthCheck(ctx)
	if err != nil {
		err = fmt.Errorf("on storage health check: %w", err)
	}

	return err
}
