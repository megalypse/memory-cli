package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/msgs"
	"github.com/megalypse/memory_cli/internal/utils"
)

type memoryCreateFooter struct {
	width, height int
	editMode      bool
}

func newMemoryCreateFooter() *memoryCreateFooter {
	return &memoryCreateFooter{
		editMode: true,
	}
}

func (m *memoryCreateFooter) Init() tea.Cmd {
	return nil
}

func (m *memoryCreateFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.editMode {
				return m, nil
			}
		case "esc":
			return m, msgs.PublishMessage(msgs.CancelCreate{})
		}
	}

	return m, nil
}

func (m *memoryCreateFooter) View() string {
	if m.editMode {
		return "Press ENTER to save, ESC to cancel"
	}

	return "Press ENTER to create, ESC to cancel"
}

func (m *memoryCreateFooter) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *memoryCreateFooter) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *memoryCreateFooter) GetHeight() int {
	return utils.CalcRatio(m.height, 0.1)
}
