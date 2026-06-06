package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/memory"
)

type Memories struct {
	memoryGroupId int
	width, height int
	memories      []*memory.Memory
	repository    memory.Repository
	err           error
}

func NewMemories(groupId int) *Memories {
	repository := memory.GetRepositorySqlLite(nil)
	return &Memories{
		memoryGroupId: groupId,
		repository:    repository,
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
		}
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

	var result string
	for i, memory := range m.memories {
		result += fmt.Sprintf("[%d] %s\n", i+1, memory.Name)
		if memory.Content != "" {
			result += fmt.Sprintf("   %s\n", memory.Content)
		}
		result += "\n"
	}

	return result
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height
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
