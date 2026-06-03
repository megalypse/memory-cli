package mainmenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/components"
)

func NewMainMenu(parent components.Container) tea.Model {
	return &MainMenu{
		parent: parent,
		footer: newFooter(parent),
	}
}

type MainMenu struct {
	parent components.Container
	footer tea.Model
}

func (m *MainMenu) Init() tea.Cmd {
	return nil
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	model, cmd := m.footer.Update(msg)
	cmdBatch = append(cmdBatch, cmd)
	if model != nil {
		return m, tea.Batch(cmdBatch...)
	}

	return m, tea.Batch(cmdBatch...)
}

func (m *MainMenu) View() string {
	return m.footer.View()
}
