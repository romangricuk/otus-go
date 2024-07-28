package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) CreateNotification(
	ctx context.Context,
	notification storage.Notification,
) (uuid.UUID, error) {
	query := `INSERT INTO notifications (id, event_id, time, message, sent) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(
		ctx, query,
		notification.ID,
		notification.EventID,
		notification.Time,
		notification.Message, notification.Sent,
	)
	return notification.ID, err
}

func (r *NotificationRepo) UpdateNotification(
	ctx context.Context,
	id uuid.UUID,
	notification storage.Notification,
) error {
	query := `UPDATE notifications SET event_id=$1, time=$2, message=$3, sent=$4 WHERE id=$5`
	_, err := r.db.ExecContext(
		ctx,
		query,
		notification.EventID,
		notification.Time,
		notification.Message,
		notification.Sent,
		id,
	)
	return err
}

func (r *NotificationRepo) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *NotificationRepo) GetNotification(ctx context.Context, id uuid.UUID) (storage.Notification, error) {
	query := `SELECT id, event_id, time, message, sent FROM notifications WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)
	var notification storage.Notification

	err := row.Scan(
		&notification.ID,
		&notification.EventID,
		&notification.Time,
		&notification.Message,
		&notification.Sent,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.Notification{}, storage.ErrNotificationNotFound
	}
	return notification, err
}

func (r *NotificationRepo) ListNotifications(
	ctx context.Context,
	start time.Time,
	end time.Time,
) ([]storage.Notification, error) {
	query := `SELECT id, event_id, time, message, sent FROM notifications WHERE time >= $1 AND time <= $2`
	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("on list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []storage.Notification
	for rows.Next() {
		var notification storage.Notification
		err = rows.Scan(
			&notification.ID,
			&notification.EventID,
			&notification.Time,
			&notification.Message,
			&notification.Sent,
		)
		if err != nil {
			return nil, fmt.Errorf("on scan notifications: %w", err)
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
