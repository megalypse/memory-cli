package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	components2 "github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/utils"
)

type footer struct {
	base

	Options *clicomponents.CursorList
}

func (m *footer) Init() tea.Cmd {
	return nil
}

func (m *footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Options.Update(msg)
}

func (m *footer) View() string {
	cellWidth := m.width / len(m.Options.Items)

	optionsBuilder := strings.Builder{}

	for i := range m.Options.Items {
		option := lipgloss.PlaceHorizontal(cellWidth, lipgloss.Center, m.Options.Items[i])
		optionsBuilder.WriteString(option)
	}

	style := lipgloss.NewStyle().Background(components2.ColorMain).Foreground(components2.ColorMainContrast)
	return style.Render(optionsBuilder.String())
}

func (m *footer) SetSize(width, height int) {
	m.base.SetSize(utils.CalcRatio(width, 1), utils.CalcRatio(height, 0.1))
}
