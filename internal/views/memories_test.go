package views

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
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
	assert.Contains(t, view.RenderFooter(), "SEARCH:")

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

func TestMemoriesDeleteSelectedMemory(t *testing.T) {
	first := &memory.Memory{ID: 1, GroupID: 7, Name: "FIRST MEMORY"}
	second := &memory.Memory{ID: 2, GroupID: 7, Name: "SECOND MEMORY"}
	repository := &memoriesRepositoryStub{
		all: []*memory.Memory{first, second},
	}
	view := NewMemories(7)
	view.repository = repository
	view.memories = repository.all
	view.cursor.Cursor = 1

	_, cmd := view.Update(msgs.DeleteMemory{})
	if cmd != nil {
		cmd()
	}

	assert.Same(t, second, repository.deleted)
	assert.Equal(t, []*memory.Memory{first}, view.memories)
	assert.Zero(t, view.cursor.Cursor)
}

func TestMemoriesDeleteRequiresConfirmation(t *testing.T) {
	footer := newMemoriesFooter(7)

	_, cmd := footer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	assert.Nil(t, cmd)

	SetLastPressedKey("d")
	t.Cleanup(func() {
		SetLastPressedKey("")
	})

	_, cmd = footer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if assert.NotNil(t, cmd) {
		assert.IsType(t, msgs.DeleteMemory{}, cmd())
	}
}

type memoriesRepositoryStub struct {
	all          []*memory.Memory
	queryResult  []*memory.Memory
	queryGroupID int
	query        string
	deleted      *memory.Memory
}

func (r *memoriesRepositoryStub) Create(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoriesRepositoryStub) Put(context.Context, *memory.Memory) error {
	return nil
}

func (r *memoriesRepositoryStub) Delete(_ context.Context, item *memory.Memory) error {
	r.deleted = item
	for index, existing := range r.all {
		if existing.ID == item.ID {
			r.all = append(r.all[:index], r.all[index+1:]...)
			break
		}
	}
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
	int,
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
