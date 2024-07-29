package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"time"

	"github.com/google/uuid"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type EventRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewEventRepo(db *sql.DB, logger logger.Logger) *EventRepo {
	return &EventRepo{
		db:     db,
		logger: logger,
	}
}

func (r *EventRepo) CreateEvent(ctx context.Context, event storage.Event) (uuid.UUID, error) {
	query := `INSERT INTO events (id, title, description, start_time, end_time, user_id) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	r.logger.Debugf("CreateEvent SQL: %s", query)

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
	)
	return event.ID, err
}

func (r *EventRepo) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	query := `UPDATE events SET title=$1, description=$2, start_time=$3, end_time=$4, user_id=$5 WHERE id=$6`
	r.logger.Debugf("UpdateEvent SQL: %s", query)

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UserID,
		id,
	)
	return err
}

func (r *EventRepo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id=$1`
	r.logger.Debugf("DeleteEvent SQL: %s", query)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *EventRepo) GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	query := `SELECT id, title, description, start_time, end_time, user_id FROM events WHERE id=$1`
	r.logger.Debugf("GetEvent SQL: %s", query)

	row := r.db.QueryRowContext(ctx, query, id)
	var event storage.Event
	err := row.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return event, err
}

func (r *EventRepo) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, description, start_time, end_time, user_id 
				FROM events WHERE start_time >= $1 AND end_time <= $2`
	r.logger.Debugf("ListEvents SQL: %s", query)

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Errorf("on closing rows in ListEvents: %v", err)
		}
	}(rows)

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
