package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/memory"
)

type MemoryDetails struct {
	memory *memory.Memory
	width  int
	height int
}

func NewMemoryDetails(mem *memory.Memory) *MemoryDetails {
	return &MemoryDetails{
		memory: mem,
	}
}

func (m *MemoryDetails) Init() tea.Cmd {
	return nil
}

func (m *MemoryDetails) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			return GetRootInstance().PopRoute()
		}
	}

	return m, nil
}

func (m *MemoryDetails) View() string {
	return fmt.Sprintf("Name: %s\nContent: %s\n\n[Press ESC or Q to return]", m.memory.Name, m.memory.Content)
}

func (m *MemoryDetails) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *MemoryDetails) GetWidth() int {
	return m.width
}

func (m *MemoryDetails) GetHeight() int {
	return m.height
}
