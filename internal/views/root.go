package views

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/utils"
)

var GetRootInstance = sync.OnceValue(func() *Root {
	root := &Root{}
	return root
})

type Root struct {
	route  View
	stack  []View
	width  int
	height int
}

func (r *Root) Init() tea.Cmd {
	if r.route == nil {
		var cmd tea.Cmd
		r.route, cmd = r.PushRoute(&MemoryGroup{})
		return cmd
	}

	return nil
}

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
	}

	_, cmd := r.route.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		SetLastPressedKey(key)
		switch key {
		case "ctrl+c":
			return r, tea.Quit
		}
	}

	return r, cmd
}

func (r *Root) View() string {
	r.resizeCurrentRoute()
	screenHeight := r.height

	title := r.getTitle()
	_, titleHeight := lipgloss.Size(title)
	footerHeight := utils.CalcRatio(screenHeight, 0.1)
	bodyHeight := screenHeight - footerHeight - titleHeight

	title = lipgloss.PlaceVertical(
		titleHeight,
		lipgloss.Center,
		title,
	)

	body := lipgloss.PlaceVertical(
		bodyHeight,
		lipgloss.Center,
		r.route.View(),
	)

	footer := lipgloss.PlaceVertical(
		footerHeight,
		lipgloss.Bottom,
		r.route.RenderFooter(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		body,
		footer,
	)
}

func (r *Root) PushRoute(route View) (View, tea.Cmd) {
	cmd := route.Init()

	r.stack = append(r.stack, r.route)
	r.route = route

	return route, cmd
}

func (r *Root) PopRoute() (View, tea.Cmd) {
	if len(r.stack) == 0 {
		return nil, nil
	}

	newRoute := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	r.route = newRoute

	return nil, nil
}

func (r *Root) resizeCurrentRoute() {
	title := r.getTitle()

	r.route.SetSize(r.width, r.height-lipgloss.Height(title))
}

func (r *Root) getTitle() string {
	title := figure.NewFigure("Memory", "", true).String()
	title = lipgloss.NewStyle().Foreground(components.ColorMain).Render(title)
	title = lipgloss.PlaceHorizontal(r.width, lipgloss.Center, title)
	return title
}
