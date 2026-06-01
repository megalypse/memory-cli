package memory_group

import (
	"context"
	"database/sql"
)

func NewRepositorySqlLite(db *sql.DB) RepositoryMemoryGroup {
	return &RepositorySqlLite{db: db}
}

type RepositorySqlLite struct {
	db *sql.DB
}

func (r *RepositorySqlLite) Create(ctx context.Context, group *MemoryGroup) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO memory_group (name, description) VALUES (?, ?)
`, group.Name, group.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositorySqlLite) GetAll(ctx context.Context) ([]*MemoryGroup, error) {
	var result []*MemoryGroup
	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, description, created_at, updated_at FROM memory_group
ORDER BY created_at DESC;
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group MemoryGroup
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &group)
	}

	return result, nil
}

func (r *RepositorySqlLite) Put(ctx context.Context, group *MemoryGroup) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE memory_group SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`, group.Name, group.Description, group.ID)
	if err != nil {
		return err
	}

	return nil
}
