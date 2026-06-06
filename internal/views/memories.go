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
		case "q", "esc":
			return m, tea.Quit
		case "enter":
			// When enter is pressed, show the full memory content in a new view or window
			// For now we'll just return nil to indicate no change in model
			return m, nil
		}
	case msgs.NewMemory:
		// Handle creating a new memory - this will be triggered by header action
		return m, func() tea.Msg {
			return m.loadMemories()
		}
	case msgs.EditMemory:
		// Handle editing a memory - this will be triggered by header action
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, nil // Would typically open edit view in real implementation
		}
	case msgs.DeleteMemory:
		// Handle deleting a memory - this will be triggered by header action
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, func() tea.Msg {
				return m.loadMemories()
			} // Reload to reflect deletion
		}
	}

	// Handle cursor list updates
	model, cmd := m.cursor.Update(msg)
	if model != nil {
		return model, cmd
	}

	// Handle footer updates (agora usando o novo tipo memoriesFooter)
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
	content := m.cursor.View()

	// Add footer if present
	if m.footer != nil {
		footer := m.footer.View()
		return fmt.Sprintf("%s\n%s", content, footer)
	}

	return content
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Configure footer size
	if m.footer != nil {
		m.footer.SetSize(width, height)
	}
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
