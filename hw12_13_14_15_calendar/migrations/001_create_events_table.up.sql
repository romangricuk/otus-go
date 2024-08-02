CREATE TABLE IF NOT EXISTS events
(
    id          UUID PRIMARY KEY,
    title       TEXT      NOT NULL,
    description TEXT,
    start_time  TIMESTAMP NOT NULL,
    end_time    TIMESTAMP NOT NULL,
    user_id     UUID      NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_start_time ON events (start_time);
CREATE INDEX IF NOT EXISTS idx_events_end_time ON events (end_time);
