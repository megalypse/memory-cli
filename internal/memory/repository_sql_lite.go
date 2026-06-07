package memory

import (
	"context"
	"database/sql"
	"strings"
	"sync"
)

var (
	repositoryInstance *RepositorySqlLite
	repositoryMu       sync.RWMutex
)

func GetRepositorySqlLite(db *sql.DB) Repository {
	if db != nil {
		repositoryMu.Lock()
		repositoryInstance = &RepositorySqlLite{db: db}
		repositoryMu.Unlock()
	}

	repositoryMu.RLock()
	defer repositoryMu.RUnlock()
	return repositoryInstance
}

type RepositorySqlLite struct {
	db *sql.DB
}

func (r *RepositorySqlLite) FindReferences(ctx context.Context, keys []string) ([]*Memory, error) {
	orMatch := strings.Join(keys, " OR ")
	var result []*Memory

	query := `
SELECT memories.id, memories.group_id, memories.name, memories.content, memories.created_at, memories.updated_at
FROM memory_fts
JOIN memories ON memories.id = memory_fts.id
WHERE memory_fts MATCH ?
ORDER BY bm25(memory_fts, 10.0, 0.0);
`

	rows, err := r.db.QueryContext(ctx, query, orMatch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memory Memory
		if err := rows.Scan(&memory.ID, &memory.GroupID, &memory.Name, &memory.Content, &memory.CreatedAt, &memory.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &memory)
	}

	return result, nil
}

func (r *RepositorySqlLite) LinkMemories(ctx context.Context, memory *Memory, memories []*Memory) error {
	if len(memories) == 0 {
		return nil
	}

	query := strings.Builder{}
	query.WriteString(`INSERT INTO memory_memory (memory_id_1, memory_id_2)`)
	query.WriteString(" VALUES ")

	placeholders := make([]string, len(memories))
	args := make([]any, 0, len(memories)*2)

	for i := range memories {
		placeholders[i] = "(?, ?)"
		args = append(args, memory.ID, memories[i].ID)
	}

	query.WriteString(strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositorySqlLite) GetRelations(ctx context.Context, memory *Memory) ([]*Memory, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT related.id, related.group_id, related.name, related.content, related.created_at, related.updated_at
FROM memory_memory
JOIN memories AS related ON related.id = memory_memory.memory_id_2
WHERE memory_memory.memory_id_1 = ?
UNION
SELECT related.id, related.group_id, related.name, related.content, related.created_at, related.updated_at
FROM memory_memory
JOIN memories AS related ON related.id = memory_memory.memory_id_1
WHERE memory_memory.memory_id_2 = ?
ORDER BY name;
`, memory.ID, memory.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []*Memory
	for rows.Next() {
		var relation Memory
		if err := rows.Scan(
			&relation.ID,
			&relation.GroupID,
			&relation.Name,
			&relation.Content,
			&relation.CreatedAt,
			&relation.UpdatedAt,
		); err != nil {
			return nil, err
		}

		relations = append(relations, &relation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}

func (r *RepositorySqlLite) Create(ctx context.Context, memory *Memory) error {
	var id int
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
INSERT INTO memories (group_id, name, content) VALUES (?, ?, ?) RETURNING id
`, memory.GroupID, memory.Name, memory.Content).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO memory_fts (id, name, content) VALUES (?, ?, ?)
`, id, memory.Name, memory.Content)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	memory.ID = id

	return nil
}

func (r *RepositorySqlLite) Put(ctx context.Context, memory *Memory) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
UPDATE memories SET group_id = ?, name = ?, content = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`, memory.GroupID, memory.Name, memory.Content, memory.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
UPDATE memory_fts SET id = ?, name = ?, content = ?
WHERE id = ?;
`, memory.ID, memory.Name, memory.Content, memory.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *RepositorySqlLite) GetAllByGroup(ctx context.Context, groupId int) ([]*Memory, error) {
	var result []*Memory
	rows, err := r.db.QueryContext(ctx, `
SELECT id, group_id, name, content, created_at, updated_at FROM memories
WHERE group_id = ?
ORDER BY created_at DESC;
`, groupId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memory Memory
		if err := rows.Scan(&memory.ID, &memory.GroupID, &memory.Name, &memory.Content, &memory.CreatedAt, &memory.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &memory)
	}

	return result, nil
}
