-- +goose Up
CREATE TABLE IF NOT EXISTS memories
(
    id         INTEGER PRIMARY KEY,
    group_id   INTEGER,
    name      TEXT,
    content    TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE VIRTUAL TABLE IF NOT EXISTS memory_fts USING fts5
(
    id UNINDEXED,
    name,
    content
);

-- +goose Down
DROP TABLE IF EXISTS memories;
DROP TABLE IF EXISTS memory_fts;
