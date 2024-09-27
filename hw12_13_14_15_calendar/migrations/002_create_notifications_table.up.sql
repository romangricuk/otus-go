CREATE TABLE IF NOT EXISTS notifications
(
    id       UUID PRIMARY KEY,
    event_id UUID REFERENCES events (id) ON DELETE CASCADE,
    time     TIMESTAMP NOT NULL,
    message  TEXT      NOT NULL,
    sent     TEXT      NOT NULL DEFAULT 'wait'
);

CREATE INDEX IF NOT EXISTS idx_notifications_time ON notifications (time);
CREATE INDEX IF NOT EXISTS idx_notifications_event_id_time ON notifications (event_id, time);
