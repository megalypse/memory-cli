package views

import (
	"context"
	"testing"

	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCreateLinksFirstReferenceForEachTerm(t *testing.T) {
	t.Parallel()

	rio := &memory.Memory{ID: 10, Name: "Rio de Janeiro"}
	brasil := &memory.Memory{ID: 20, Name: "História do Brasil"}
	repository := &memoryCreateRepositoryStub{
		references: map[string][]*memory.Memory{
			"Rio de Janeiro": {
				rio,
				{ID: 11, Name: "Outra referência"},
			},
			"História do Brasil": {
				brasil,
				{ID: 21, Name: "Referência secundária"},
			},
			"Brasil": {
				brasil,
			},
		},
	}
	view := &MemoryCreate{
		memoryGroupId: 7,
		repository:    repository,
	}

	err := view.createMemory(
		context.Background(),
		"viagem ao Rio de Janeiro",
		"estudo sobre História do Brasil e Brasil.",
	)

	assert.NoError(t, err)
	assert.Len(t, repository.findKeys, 3)
	assert.NotNil(t, repository.created)
	assert.Equal(t, 99, repository.created.ID)
	assert.Equal(t, 7, repository.created.GroupID)
	assert.Same(t, repository.created, repository.linkedMemory)
	assert.Equal(t, []*memory.Memory{rio, brasil}, repository.linkedReferences)
}

type memoryCreateRepositoryStub struct {
	references       map[string][]*memory.Memory
	findKeys         []string
	created          *memory.Memory
	linkedMemory     *memory.Memory
	linkedReferences []*memory.Memory
}

func (r *memoryCreateRepositoryStub) Create(_ context.Context, item *memory.Memory) error {
	item.ID = 99
	r.created = item
	return nil
}

func (r *memoryCreateRepositoryStub) Put(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoryCreateRepositoryStub) GetAllByGroup(context.Context, int) ([]*memory.Memory, error) {
	return nil, nil
}

func (r *memoryCreateRepositoryStub) LinkMemories(
	_ context.Context,
	item *memory.Memory,
	references []*memory.Memory,
) error {
	r.linkedMemory = item
	r.linkedReferences = references
	return nil
}

func (r *memoryCreateRepositoryStub) FindReferences(
	_ context.Context,
	keys []string,
) ([]*memory.Memory, error) {
	r.findKeys = append(r.findKeys, keys[0])
	return r.references[keys[0]], nil
}
