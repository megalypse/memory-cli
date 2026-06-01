package memory_group

import "context"

type RepositoryMemoryGroup interface {
	Create(ctx context.Context, group *MemoryGroup) error
	GetAll(ctx context.Context) ([]*MemoryGroup, error)
	Put(ctx context.Context, group *MemoryGroup) error
}
