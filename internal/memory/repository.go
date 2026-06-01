package memory

import "context"

type Repository interface {
	Create(ctx context.Context, memory *Memory) error
	Put(ctx context.Context, memory *Memory) error
	GetAllByGroup(ctx context.Context, groupId int) ([]*Memory, error)
	LinkMemories(ctx context.Context, memory *Memory, memories []*Memory) error
	FindReferences(ctx context.Context, keys []string) ([]*Memory, error)
}
