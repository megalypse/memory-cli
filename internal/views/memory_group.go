package views

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type MemoryGroup struct {
	base

	groupsCursor *clicomponents.CursorList
	groups       []*memorygroup.MemoryGroup

	footer *footer

	repository memorygroup.RepositoryMemoryGroup
}

func (m *MemoryGroup) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(width, height)
}

func (m *MemoryGroup) Init() tea.Cmd {
	m.footer = &footer{}

	cursor := &clicomponents.CursorList{
		Items:      []string{"(a) Add", "(e) Edit", "(d -> y) Delete"},
		RenderSize: 3,
	}

	cursor.KeyActions = []clicomponents.OnKeyFn{
		clicomponents.OnKey("a", func() (tea.Model, tea.Cmd) {
			return nil, msgs.PublishMessage(msgs.NewGroup{})
		}),
		clicomponents.OnKeys(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "e"),
		clicomponents.OnKeys(func() (tea.Model, tea.Cmd) {
			if GetLastPressedKey() == "d" {
				SetLastPressedKey("")

				return nil, msgs.PublishMessage(msgs.DeleteGroup{})
			}

			return nil, nil
		}, "y"),
		clicomponents.OnKey("enter", func() (tea.Model, tea.Cmd) {
			if len(m.groups) == 0 {
				return m, nil
			}

			group := m.groups[m.groupsCursor.Cursor]
			newRoute := NewMemories(group.ID)
			newRoute.SetSize(m.width, m.height)
			return GetRootInstance().PushRoute(newRoute)
		}),
	}

	m.footer.Options = cursor
	m.footer.Init()

	m.groupsCursor = &clicomponents.CursorList{
		RenderSize: 3,
	}

	m.repository = memorygroup.GetRepositorySqlLite(nil)
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
		newRoute := NewMemoryGroupCreate()
		return GetRootInstance().PushRoute(newRoute)
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
	return m.groupsCursor.View()
}

func (m *MemoryGroup) RenderFooter() string {
	return m.footer.View()
}
