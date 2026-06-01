package memory

import (
	"context"
	"database/sql"
	"strings"
)

func NewRepositorySqlLite(db *sql.DB) *RepositorySqlLite {
	return &RepositorySqlLite{db: db}
}

type RepositorySqlLite struct {
	db *sql.DB
}

func (r *RepositorySqlLite) FindReferences(ctx context.Context, keys []string) ([]*Memory, error) {
	orMatch := strings.Join(keys, " OR ")
	var result []*Memory

	query := `
SELECT *, bm25(memories_fts, 10.0, 2.0) AS score FROM memories
JOIN memory_fts ON memories.id = memory_fts.id
WHERE memory_fts MATCH ?
ORDER BY score;
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
	query := strings.Builder{}
	query.WriteString(`INSERT INTO memory_memory (memory_id_1, memory_id_2)`)
	query.WriteString(" VALUES ")

	placeholders := make([]string, len(memories))
	args := make([]any, 0, len(memories)*2)

	for i, _ := range memories {
		if i > 0 {
			query.WriteString(", ")
		}

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

func (r *RepositorySqlLite) Create(ctx context.Context, memory *Memory) error {
	var id int
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
INSERT INTO memories (group_id, name, description) VALUES (?, ?, ?) RETURNING id
`, memory.GroupID, memory.Name, memory.Content).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO memory_fts (rowid, name, description) VALUES (?, ?, ?)
`, id, memory.Name, memory.Content)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *RepositorySqlLite) Put(ctx context.Context, memory *Memory) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE memories SET group_id = ?, name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`, memory.GroupID, memory.Name, memory.Content, memory.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositorySqlLite) GetAllByGroup(ctx context.Context, groupId int) ([]*Memory, error) {
	var result []*Memory
	rows, err := r.db.QueryContext(ctx, `
SELECT id, group_id, name, description, created_at, updated_at FROM memories
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
