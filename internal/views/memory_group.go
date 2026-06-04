package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/memorygroup"
)

type MemoryGroup struct {
	base

	groupsCursor *components.CursorListVertical
	groups       []*memorygroup.MemoryGroup

	footer *memoryGroupFooter
}

func (m *MemoryGroup) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.widthRatio = 1
	m.heightRatio = 0.9

	m.footer.SetSize(width, height)
}

func (m *MemoryGroup) Init() tea.Cmd {
	m.footer = &memoryGroupFooter{}
	m.footer.Init()

	repository := memorygroup.NewRepositorySqlLite(nil)
	m.groups, _ = repository.GetAll(context.Background())

	items := make([]string, len(m.groups))
	for i, group := range m.groups {
		items[i] = group.Name
	}

	m.groupsCursor = &components.CursorListVertical{
		Items:      items,
		RenderSize: 3,
	}

	return nil
}

func (m *MemoryGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList []tea.Cmd

	model, cmd := m.groupsCursor.Update(msg)
	cmdList = append(cmdList, cmd)
	if model != nil {
		return model, tea.Batch(cmdList...)
	}

	return m, tea.Batch(cmdList...)
}

func (m *MemoryGroup) View() string {
	groups := lipgloss.Place(m.GetWidth(), m.GetHeight(), lipgloss.Center, lipgloss.Center, m.groupsCursor.View())
	footer := lipgloss.Place(m.footer.GetWidth(), m.footer.GetHeight(), lipgloss.Center, lipgloss.Center, m.footer.View())

	return lipgloss.JoinVertical(lipgloss.Center, groups, footer)
}
