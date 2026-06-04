package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type memoryGroupFooter struct {
	base

	optionsCursor *components.CursorListVertical
}

func (m *memoryGroupFooter) Init() tea.Cmd {
	cursor := &components.CursorListVertical{
		Items:      []string{"(a) Add", "(e) Edit", "(d) Delete"},
		RenderSize: 3,
	}

	cursor.KeyActions = []components.OnKeyFn{
		components.OnKey(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "a"),
		components.OnKey(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "e"),
		components.OnKey(func() (tea.Model, tea.Cmd) {
			if GetLastPressedKey() == "d" {
				SetLastPressedKey("")

				return nil, func() tea.Msg {
					return msgs.DeleteGroup{}
				}
			}

			return nil, nil
		}, "d"),
	}
	return nil
}

func (m *memoryGroupFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *memoryGroupFooter) View() string {
	return ""
}

func (m *memoryGroupFooter) SetSize(width, height int) {
	m.base.SetSize(width, height)
}
