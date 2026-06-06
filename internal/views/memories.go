package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/memory"
)

type Memories struct {
	memoryGroupId int
	width, height int
	memories      []*memory.Memory
	repository    memory.Repository
	err           error
	cursor        *clicomponents.CursorList
}

func NewMemories(groupId int) *Memories {
	repository := memory.GetRepositorySqlLite(nil)
	return &Memories{
		memoryGroupId: groupId,
		repository:    repository,
		cursor: &clicomponents.CursorList{
			RenderSize: 10,
		},
	}
}

func (m *Memories) Init() tea.Cmd {
	return m.loadMemories
}

func (m *Memories) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "enter":
			// When enter is pressed, show the full memory content in a new view or window
			// For now we'll just return nil to indicate no change in model
			return m, nil
		}
	}

	// Handle cursor list updates
	model, cmd := m.cursor.Update(msg)
	if model != nil {
		return model, cmd
	}

	return m, nil
}

func (m *Memories) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error loading memories: %v", m.err)
	}

	if len(m.memories) == 0 {
		return "No memories found for this group."
	}

	// Create items for the cursor list
	items := make([]string, len(m.memories))
	for i, memory := range m.memories {
		items[i] = fmt.Sprintf("[%d] %s", i+1, memory.Name)
	}
	m.cursor.Items = items

	// Render with cursor
	return m.cursor.View()
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height
	// Since we're using a CursorList that handles scrolling internally,
	// we don't need to set a specific Height here
}

func (m *Memories) GetWidth() int {
	return m.width
}

func (m *Memories) GetHeight() int {
	return m.height
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
