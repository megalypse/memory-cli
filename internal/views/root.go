package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/megalypse/memory_cli/internal/components"
)

type Root struct {
	route  View
	stack  []View
	width  int
	height int
}

func (r *Root) Init() tea.Cmd {
	if r.route == nil {
		newView := &MemoryGroup{}
		initCmd := newView.Init()
		r.PushRoute(newView)

		return initCmd
	}

	return nil
}

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		r.resizeCurrentRoute()
	}

	_, cmd := r.route.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		SetLastPressedKey(key)

	}

	return r, cmd
}

func (r *Root) View() string {
	title := figure.NewFigure("Memory", "", true).String()
	title = lipgloss.NewStyle().Foreground(components.ColorMain).Render(title)

	return lipgloss.JoinVertical(lipgloss.Center, title, r.route.View())
}

func (r *Root) PushRoute(route View) tea.Cmd {
	cmd := route.Init()
	r.route = route
	return cmd
}

func (r *Root) PopRoute() tea.Cmd {
	if len(r.stack) == 0 {
		return nil
	}

	newRoute := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]

	return r.PushRoute(newRoute)
}

func (r *Root) resizeCurrentRoute() {
	title := figure.NewFigure("Memory", "", true).String()

	r.route.SetSize(r.width, r.height-lipgloss.Height(title))
}
