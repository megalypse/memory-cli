package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/memorygroup"
)

type MemoryGroup struct {
	width  int
	height int

	groupsCursor *components.CursorListVertical
	groups       []*memorygroup.MemoryGroup
}

func (m *MemoryGroup) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *MemoryGroup) Init() tea.Cmd {
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
	return m.groupsCursor.View()
}
