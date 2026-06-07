package memory

import "context"

type Repository interface {
	Create(ctx context.Context, memory *Memory) error
	Put(ctx context.Context, memory *Memory) error
	Delete(ctx context.Context, memory *Memory) error
	GetAllByGroup(ctx context.Context, groupId int) ([]*Memory, error)
	GetRelations(ctx context.Context, memory *Memory) ([]*Memory, error)
	LinkMemories(ctx context.Context, memory *Memory, memories []*Memory) error
	FindReferences(ctx context.Context, groupID int, keys []string) ([]*Memory, error)
	QueryMemories(ctx context.Context, groupID int, query string) ([]*Memory, error)
}
