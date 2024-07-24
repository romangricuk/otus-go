package memorystorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Connect(t *testing.T) {
	memStore := New()
	err := memStore.Connect(context.Background())
	assert.NoError(t, err, "expected no error on connect")
}

func TestMemoryStorage_Close(t *testing.T) {
	memStore := New()
	err := memStore.Close()
	assert.NoError(t, err, "expected no error on close")
}

func TestMemoryStorage_EventRepository(t *testing.T) {
	memStore := New()
	repo := memStore.EventRepository()
	assert.NotNil(t, repo, "expected non-nil EventRepository")
}

func TestMemoryStorage_NotificationRepository(t *testing.T) {
	memStore := New()
	repo := memStore.NotificationRepository()
	assert.NotNil(t, repo, "expected non-nil NotificationRepository")
}
