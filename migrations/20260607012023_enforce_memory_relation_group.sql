-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER enforce_memory_relation_group
BEFORE INSERT ON memory_memory
FOR EACH ROW
WHEN (
    SELECT group_id FROM memories WHERE id = NEW.memory_id_1
) != (
    SELECT group_id FROM memories WHERE id = NEW.memory_id_2
)
BEGIN
    SELECT RAISE(ABORT, 'MEMORIES MUST BELONG TO THE SAME GROUP');
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS enforce_memory_relation_group;
