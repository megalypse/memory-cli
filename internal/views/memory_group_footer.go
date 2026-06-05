package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/components"
	components2 "github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type memoryGroupFooter struct {
	base

	optionsCursor *components.CursorListVertical
}

func (m *memoryGroupFooter) Init() tea.Cmd {
	m.widthRatio = 1
	m.heightRatio = 0.1

	cursor := &components.CursorListVertical{
		Items:      []string{"(a) Add", "(e) Edit", "(d -> y) Delete"},
		RenderSize: 3,
	}

	cursor.KeyActions = []components.OnKeyFn{
		components.OnKey("a", func() (tea.Model, tea.Cmd) {
			return nil, msgs.PublishMessage(msgs.NewGroup{})
		}),
		components.OnKeys(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "e"),
		components.OnKeys(func() (tea.Model, tea.Cmd) {
			if GetLastPressedKey() == "d" {
				SetLastPressedKey("")

				return nil, msgs.PublishMessage(msgs.DeleteGroup{})
			}

			return nil, nil
		}, "y"),
	}

	m.optionsCursor = cursor
	return nil
}

func (m *memoryGroupFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.optionsCursor.Update(msg)
}

func (m *memoryGroupFooter) View() string {
	cellWidth := m.width / len(m.optionsCursor.Items)

	optionsBuilder := strings.Builder{}

	for i := range m.optionsCursor.Items {
		option := lipgloss.PlaceHorizontal(cellWidth, lipgloss.Center, m.optionsCursor.Items[i])
		optionsBuilder.WriteString(option)
	}

	style := lipgloss.NewStyle().Background(components2.ColorMain).Foreground(components2.ColorMainContrast)
	return style.Render(optionsBuilder.String())
}

func (m *memoryGroupFooter) SetSize(width, height int) {
	m.base.SetSize(width, height)
}
