package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type Memories struct {
	memoryGroupId int
	width, height int
	memories      []*memory.Memory
	repository    memory.Repository
	err           error
	cursor        *clicomponents.CursorList
	footer        *memoriesFooter
}

func NewMemories(groupId int) *Memories {
	repository := memory.GetRepositorySqlLite(nil)
	return &Memories{
		memoryGroupId: groupId,
		repository:    repository,
		cursor: &clicomponents.CursorList{
			RenderSize: 10,
		},
		footer: newMemoriesFooter(groupId),
	}
}

func (m *Memories) Init() tea.Cmd {
	return func() tea.Msg {
		return m.loadMemories()
	}
}

func (m *Memories) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return GetRootInstance().PopRoute()

		case "q":
			return m, tea.Quit
		case "enter":
			if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
				selectedMemory := m.memories[m.cursor.Cursor]
				detailsView := NewMemoryDetails(selectedMemory)
				detailsView.SetSize(m.width, m.height)

				return GetRootInstance().PushRoute(detailsView)
			}
			return m, nil
		case "n":
			createView := NewMemoryCreate(m.memoryGroupId)

			return GetRootInstance().PushRoute(createView)
		}
	case msgs.NewMemory:
		return m, func() tea.Msg {
			return m.loadMemories()
		}
	case msgs.EditMemory:
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, nil
		}
	case msgs.DeleteMemory:
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, func() tea.Msg {
				return m.loadMemories()
			}
		}
	}

	model, cmd := m.cursor.Update(msg)
	if model != nil {
		return model, cmd
	}

	model, cmd = m.footer.Update(msg)
	if model != nil {
		return model, cmd
	}

	return m, nil
}

func (m *Memories) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error loading memories: %v", m.err)
	}

	items := make([]string, len(m.memories))
	for i, memory := range m.memories {
		items[i] = fmt.Sprintf("[%d] %s", i+1, memory.Name)
	}
	m.cursor.Items = items

	return m.cursor.View()
}

func (m *Memories) RenderFooter() string {
	return m.footer.View()
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(width, height)
}

func (m *Memories) loadMemories() tea.Msg {
	memories, err := m.repository.GetAllByGroup(context.Background(), m.memoryGroupId)
	if err != nil {
		m.err = err
		return nil
	}

	m.memories = memories
	return nil
}
