package memorygroup

import "context"

type RepositoryMemoryGroup interface {
	Create(ctx context.Context, group *MemoryGroup) error
	GetAll(ctx context.Context) ([]*MemoryGroup, error)
	Put(ctx context.Context, group *MemoryGroup) error
	Delete(ctx context.Context, id int) error
}
