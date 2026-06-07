package views

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoriesSearch(t *testing.T) {
	repository := &memoriesRepositoryStub{
		all: []*memory.Memory{
			{ID: 1, Name: "First Memory"},
			{ID: 2, Name: "Second Memory"},
		},
		queryResult: []*memory.Memory{
			{ID: 2, Name: "Second Memory"},
		},
	}
	view := NewMemories(7)
	view.repository = repository
	view.memories = repository.all
	view.SetSize(90, 24)

	listBeforeSearch := view.View()

	view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	assert.True(t, view.searching)
	assert.Equal(t, listBeforeSearch, view.View())
	assert.Contains(t, view.RenderFooter(), "Search:")

	view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("second")})
	result := view.queryMemories("second")()
	view.Update(result)

	assert.Equal(t, 7, repository.queryGroupID)
	assert.Equal(t, "second", repository.query)
	assert.Equal(t, repository.queryResult, view.memories)
	assert.Contains(t, view.View(), "Second Memory")

	view.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.False(t, view.searching)
	assert.Empty(t, view.searchInput.Value())
}

type memoriesRepositoryStub struct {
	all          []*memory.Memory
	queryResult  []*memory.Memory
	queryGroupID int
	query        string
}

func (r *memoriesRepositoryStub) Create(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoriesRepositoryStub) Put(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoriesRepositoryStub) GetAllByGroup(context.Context, int) ([]*memory.Memory, error) {
	return r.all, nil
}

func (r *memoriesRepositoryStub) GetRelations(
	context.Context,
	*memory.Memory,
) ([]*memory.Memory, error) {
	return nil, nil
}

func (r *memoriesRepositoryStub) LinkMemories(
	context.Context,
	*memory.Memory,
	[]*memory.Memory,
) error {
	return nil
}

func (r *memoriesRepositoryStub) FindReferences(
	context.Context,
	[]string,
) ([]*memory.Memory, error) {
	return nil, nil
}

func (r *memoriesRepositoryStub) QueryMemories(
	_ context.Context,
	groupID int,
	query string,
) ([]*memory.Memory, error) {
	r.queryGroupID = groupID
	r.query = query
	return r.queryResult, nil
}
