-- +goose Up
CREATE TABLE IF NOT EXISTS memory_memory (
    id INTEGER PRIMARY KEY,
    memory_id_1 INTEGER,
    memory_id_2 INTEGER,
    created_at TIMESTAMP,

    CONSTRAINT fk_memory_1 FOREIGN KEY (memory_id_1) REFERENCES memories(id) ON DELETE CASCADE,
    CONSTRAINT fk_memory_2 FOREIGN KEY (memory_id_2) REFERENCES memories(id) ON DELETE CASCADE,
    UNIQUE (memory_id_1, memory_id_2)
);

-- +goose Down
DROP TABLE IF EXISTS memory_memory;
