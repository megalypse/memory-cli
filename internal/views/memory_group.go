package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type MemoryGroup struct {
	base

	groupsCursor *components.CursorListVertical
	groups       []*memorygroup.MemoryGroup

	footer *memoryGroupFooter

	repository memorygroup.RepositoryMemoryGroup
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

	m.groupsCursor = &components.CursorListVertical{
		RenderSize: 3,
	}

	m.repository = memorygroup.NewRepositorySqlLite(nil)
	m.updateGroups()

	return nil
}

func (m *MemoryGroup) updateGroups() tea.Cmd {
	var err error
	m.groups, err = m.repository.GetAll(context.Background())
	if err != nil {
		return msgs.PublishMessage(msgs.Err{Err: err})
	}

	items := make([]string, len(m.groups))
	for i, group := range m.groups {
		items[i] = group.Name
	}

	m.groupsCursor.Items = items
	if len(m.groups) == 0 {
		m.groupsCursor.Cursor = 0
		return nil
	}

	if m.groupsCursor.Cursor > len(m.groups)-1 {
		m.groupsCursor.Cursor = len(m.groups) - 1
	}

	return nil
}

func (m *MemoryGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList []tea.Cmd

	switch msg.(type) {
	case msgs.NewGroup:
		return getRootInstance().PushRoute(&MemoryGroupCreate{})
	case msgs.UpdateGroups:
		return m, m.updateGroups()
	case msgs.DeleteGroup:
		if len(m.groups) == 0 {
			return m, nil
		}

		group := m.groups[m.groupsCursor.Cursor]
		err := m.repository.Delete(context.Background(), group.ID)
		if err != nil {
			return m, msgs.PublishMessage(msgs.Err{Err: err})
		}

		return m, msgs.PublishMessage(msgs.UpdateGroups{})
	}

	model, cmd := m.groupsCursor.Update(msg)
	cmdList = append(cmdList, cmd)
	if model != nil {
		return model, tea.Batch(cmdList...)
	}

	model, cmd = m.footer.Update(msg)
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
