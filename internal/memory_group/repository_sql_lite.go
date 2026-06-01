package memory_group

import (
	"context"
	"database/sql"
)

func NewRepositorySqlLiteMemoryGroup(db *sql.DB) RepositoryMemoryGroup {
	return &RepositorySqlLiteMemoryGroup{db: db}
}

type RepositorySqlLiteMemoryGroup struct {
	db *sql.DB
}

func (r *RepositorySqlLiteMemoryGroup) Create(ctx context.Context, group *MemoryGroup) error {
	r.db.ExecContext(ctx, `
INSERT INTO memory_group (name, description) VALUES (?, ?)
`, group.Name, group.Description)
	return nil
}

func (r *RepositorySqlLiteMemoryGroup) GetAll(ctx context.Context) ([]*MemoryGroup, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RepositorySqlLiteMemoryGroup) Put(ctx context.Context, group *MemoryGroup) error {
	//TODO implement me
	panic("implement me")
}
