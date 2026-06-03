package mainmenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/components"
	components2 "github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/msgs"
)

func newFooter(parent components2.Container) tea.Model {
	cursor := &components.CursorListVertical{
		Items: []string{
			"(N) Novo grupo",
			"(E) Editar grupo",
			"(D) Deletar grupo",
		},
	}

	cursor.KeyActions = []components.OnKeyFn{
		components.OnKey(func() (tea.Model, tea.Cmd) {
			return nil, func() tea.Msg {
				return msgs.MsgNewGroup()
			}
		}, "n"),
		components.OnKey(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "e"),
		components.OnKey(func() (tea.Model, tea.Cmd) {
			return nil, nil
		}, "d"),
	}

	return &footer{
		options: cursor,
		parent:  parent,
	}
}

type footer struct {
	options *components.CursorListVertical
	parent  components2.Container
}

func (f *footer) Init() tea.Cmd {
	return nil
}

func (f *footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := f.options.Update(msg)
	if model != nil {
		return model, cmd
	}

	return nil, cmd
}

func (f *footer) View() string {
	optionInterval := 1.0 / len(f.options.Items)
	cellWidth := f.parent.GetWidth() / len(f.options.Items)
	positions := make([]lipgloss.Position, len(f.options.Items))
	for i := range f.options.Items {
		positions[i] = lipgloss.Position(optionInterval * i)
	}

	var options []string
	for _, option := range f.options.Items {
		placed := lipgloss.PlaceHorizontal(cellWidth, lipgloss.Center, option)
		options = append(options, placed)
	}

	return lipgloss.NewStyle().
		Background(components2.ColorMain).
		Foreground(lipgloss.Color("#121105")).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, options...))
}
