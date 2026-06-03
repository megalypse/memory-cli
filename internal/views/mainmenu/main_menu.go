package mainmenu

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	components2 "github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/megalypse/memory_cli/internal/msgs"
)

func NewMainMenu(parent components.Container) tea.Model {
	mainMenu := &MainMenu{
		footer: newFooter(parent),
		groupCursor: components2.CursorListVertical{
			RenderSize: 3,
		},
	}

	mainMenu.updateGroups()

	return mainMenu
}

type MainMenu struct {
	footer      tea.Model
	groupCursor components2.CursorListVertical
	groups      []*memorygroup.MemoryGroup
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

	model, cmd = m.groupCursor.Update(msg)
	cmdBatch = append(cmdBatch, cmd)
	if model != nil {
		return m, tea.Batch(cmdBatch...)
	}

	switch msg.(type) {
	case msgs.NewGroup:
		m.updateGroups()
	}

	return m, tea.Batch(cmdBatch...)
}

func (m *MainMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, m.groupCursor.View(), m.footer.View())
}

func (m *MainMenu) updateGroups() {
	repository := memorygroup.NewRepositorySqlLite(nil)
	groups, _ := repository.GetAll(context.Background())

	var groupNames []string
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}

	m.groupCursor.Items = groupNames
	m.groups = groups
}
