package views

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryDetailsLoadsAndNavigatesRelations(t *testing.T) {
	current := &memory.Memory{ID: 1, Name: "Current Memory"}
	first := &memory.Memory{ID: 2, Name: "First Relation"}
	second := &memory.Memory{ID: 3, Name: "Second Relation"}
	repository := &memoryDetailsRepositoryStub{
		relations: []*memory.Memory{first, second},
	}
	details := NewMemoryDetails(current)
	details.repository = repository
	details.SetSize(90, 30)

	msg := details.loadRelations()
	model, cmd := details.Update(msg)

	assert.Same(t, details, model)
	assert.Nil(t, cmd)
	assert.Equal(t, []*memory.Memory{first, second}, details.relations)
	assert.Contains(t, details.View(), "Details")
	assert.Contains(t, details.View(), "Relations")
	assert.Contains(t, details.View(), "> First Relation")

	model, cmd = details.Update(tea.KeyMsg{Type: tea.KeyDown})

	assert.Same(t, details, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, details.cursor.Cursor)
	assert.Contains(t, details.View(), "> Second Relation")
}

func TestMemoryDetailsHeadersDoNotWrap(t *testing.T) {
	details := NewMemoryDetails(&memory.Memory{Name: "Current Memory"})

	assert.Equal(t, 1, lipgloss.Height(details.renderHeader("Details", 20)))
	assert.Equal(t, 1, lipgloss.Height(details.renderHeader("Relations", 10)))
}

func TestMemoryDetailsEnterReplacesRouteWithoutGrowingStack(t *testing.T) {
	current := &memory.Memory{ID: 1, Name: "Current Memory"}
	first := &memory.Memory{ID: 2, Name: "First Relation"}
	second := &memory.Memory{ID: 3, Name: "Second Relation"}
	details := NewMemoryDetails(current)
	details.relations = []*memory.Memory{first, second}
	details.cursor.Items = []string{first.Name, second.Name}
	details.cursor.Cursor = 1
	details.SetSize(90, 30)

	root := GetRootInstance()
	previousRoute := root.route
	previousStack := root.stack
	defer func() {
		root.route = previousRoute
		root.stack = previousStack
	}()

	parent := &MemoryDetails{memory: &memory.Memory{ID: 10}}
	root.route = details
	root.stack = []View{parent}

	model, _ := details.Update(tea.KeyMsg{Type: tea.KeyEnter})

	selectedDetails, ok := model.(*MemoryDetails)
	assert.True(t, ok)
	assert.Same(t, second, selectedDetails.memory)
	assert.Same(t, selectedDetails, root.route)
	assert.Equal(t, []View{parent}, root.stack)
}

type memoryDetailsRepositoryStub struct {
	relations []*memory.Memory
	err       error
}

func (r *memoryDetailsRepositoryStub) Create(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoryDetailsRepositoryStub) Put(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoryDetailsRepositoryStub) GetAllByGroup(context.Context, int) ([]*memory.Memory, error) {
	return nil, nil
}

func (r *memoryDetailsRepositoryStub) GetRelations(
	context.Context,
	*memory.Memory,
) ([]*memory.Memory, error) {
	return r.relations, r.err
}

func (r *memoryDetailsRepositoryStub) LinkMemories(
	context.Context,
	*memory.Memory,
	[]*memory.Memory,
) error {
	return nil
}

func (r *memoryDetailsRepositoryStub) FindReferences(
	context.Context,
	[]string,
) ([]*memory.Memory, error) {
	return nil, nil
}
