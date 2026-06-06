package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/msgs"
	"github.com/megalypse/memory_cli/internal/utils"
)

type memoriesFooter struct {
	base

	Options *clicomponents.CursorList
}

func (m *memoriesFooter) Init() tea.Cmd {
	return nil
}

func (m *memoriesFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle key actions specifically for this footer
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Handle adding a new memory
			return m, msgs.PublishMessage(msgs.NewMemory{})
		case "e":
			// Handle editing - this is just a placeholder in footer; main view will handle actual editing logic
			return m, msgs.PublishMessage(msgs.EditMemory{})
		case "d":
			// Handle deleting - this is just a placeholder in footer; main view will handle actual deletion logic
			return m, msgs.PublishMessage(msgs.DeleteMemory{})
		}
	}

	model, cmd := m.Options.Update(msg)
	return model, cmd
}

func (m *memoriesFooter) View() string {
	cellWidth := m.width / len(m.Options.Items)

	optionsBuilder := strings.Builder{}

	for i := range m.Options.Items {
		option := lipgloss.PlaceHorizontal(cellWidth, lipgloss.Center, m.Options.Items[i])
		optionsBuilder.WriteString(option)
	}

	style := lipgloss.NewStyle().Background(components.ColorMain).Foreground(components.ColorMainContrast)
	return style.Render(optionsBuilder.String())
}

func (m *memoriesFooter) SetSize(width, height int) {
	m.base.SetSize(utils.CalcRatio(width, 1), utils.CalcRatio(height, 0.1))
}
