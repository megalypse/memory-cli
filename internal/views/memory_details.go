package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/memory"
)

type MemoryDetails struct {
	memory *memory.Memory
	width  int
	height int
	footer *footer
}

func NewMemoryDetails(mem *memory.Memory) *MemoryDetails {
	return &MemoryDetails{
		memory: mem,
		footer: &footer{
			Options: &clicomponents.CursorList{
				Items: []string{"ESC/Q: Return"},
			},
		},
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

	model, cmd := m.footer.Update(msg)
	if model != nil {
		return model, cmd
	}

	return m, nil
}

func (m *MemoryDetails) View() string {
	body := fmt.Sprintf("Name: %s\nContent: %s", m.memory.Name, m.memory.Content)

	return body
}

func (m *MemoryDetails) RenderFooter() string {
	return m.footer.View()
}

func (m *MemoryDetails) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(m.width, m.height)
}
