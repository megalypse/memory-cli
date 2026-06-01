-- +goose Up
CREATE TABLE IF NOT EXISTS memory_group
(
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS memory_group;
