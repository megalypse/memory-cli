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
	third := &memory.Memory{ID: 4, Name: "Third Relation"}
	fourth := &memory.Memory{ID: 5, Name: "Fourth Relation"}
	fifth := &memory.Memory{ID: 6, Name: "Fifth Relation"}
	repository := &memoryDetailsRepositoryStub{
		relations: map[int][]*memory.Memory{
			current.ID: {first, second},
			first.ID:   {current, second, third},
			second.ID:  {current, third, fourth},
			third.ID:   {first, fifth},
			fourth.ID:  {second, fifth},
		},
	}
	details := NewMemoryDetails(current)
	details.repository = repository
	details.SetSize(90, 30)

	msg := details.loadRelations()
	model, cmd := details.Update(msg)

	assert.Same(t, details, model)
	assert.Nil(t, cmd)
	assert.Equal(t, []*memory.Memory{first, second}, details.relations[0])
	assert.Equal(t, []*memory.Memory{third, fourth}, details.relations[1])
	assert.Equal(t, []*memory.Memory{fifth}, details.relations[2])
	assert.Contains(t, details.View(), "Details")
	assert.Contains(t, details.View(), "Relations")
	assert.Contains(t, details.View(), "1º")
	assert.Contains(t, details.View(), "2º")
	assert.Contains(t, details.View(), "3º")
	assert.Contains(t, details.View(), "> First Relation")
	assert.NotContains(t, details.View(), "> Third Relation")
	assert.NotContains(t, details.View(), "> Fifth Relation")

	model, cmd = details.Update(tea.KeyMsg{Type: tea.KeyDown})

	assert.Same(t, details, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, details.cursors[0].Cursor)
	assert.Contains(t, details.View(), "> Second Relation")

	model, cmd = details.Update(tea.KeyMsg{Type: tea.KeyRight})

	assert.Same(t, details, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, details.activeLevel)
	assert.Contains(t, details.View(), "> Third Relation")
	assert.NotContains(t, details.View(), "> Second Relation")
}

func TestMemoryDetailsHeadersDoNotWrap(t *testing.T) {
	details := NewMemoryDetails(&memory.Memory{Name: "Current Memory"})

	assert.Equal(t, 1, lipgloss.Height(details.renderHeader("Details", 20)))
	assert.Equal(t, 1, lipgloss.Height(details.renderHeader("Relations", 10)))
}

func TestMemoryDetailsUsesFullBodyHeight(t *testing.T) {
	details := NewMemoryDetails(&memory.Memory{Name: "Current Memory"})
	details.relations[0] = []*memory.Memory{{ID: 2, Name: "First Relation"}}
	details.SetSize(90, 24)

	view := details.View()

	assert.Equal(t, 24, lipgloss.Height(view))
	assert.Equal(t, 9, details.cursors[0].RenderSize)
	assert.Equal(t, 9, details.cursors[1].RenderSize)
	assert.Equal(t, 9, details.cursors[2].RenderSize)
}

func TestMemoryDetailsEnterReplacesRouteWithoutGrowingStack(t *testing.T) {
	current := &memory.Memory{ID: 1, Name: "Current Memory"}
	first := &memory.Memory{ID: 2, Name: "First Relation"}
	second := &memory.Memory{ID: 3, Name: "Second Relation"}
	details := NewMemoryDetails(current)
	details.relations[1] = []*memory.Memory{first, second}
	details.cursors[1].Items = []string{first.Name, second.Name}
	details.cursors[1].Cursor = 1
	details.activeLevel = 1
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
	relations map[int][]*memory.Memory
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
	_ context.Context,
	item *memory.Memory,
) ([]*memory.Memory, error) {
	return r.relations[item.ID], r.err
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
