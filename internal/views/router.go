package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/views/mainmenu"
)

func NewRouter(startRoute tea.Model) *Router {
	return &Router{
		router: startRoute,
	}
}

type Router struct {
	router  tea.Model
	history []tea.Model
	height  int
	width   int
}

func (r *Router) GetWidth() int {
	return r.width
}

func (r *Router) GetHeight() int {
	return r.height
}

func (r *Router) Push(model tea.Model) (tea.Model, tea.Cmd) {
	r.history = append(r.history, r.router)
	return model, model.Init()
}

func (r *Router) Pop() (tea.Model, tea.Cmd) {
	if len(r.history) == 0 {
		return r.router, nil
	}

	next := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]

	return next, nil
}

func (r *Router) Init() tea.Cmd {
	return nil
}

func (r *Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
	}

	var cmd tea.Cmd
	r.router, cmd = r.router.Update(msg)

	return r, cmd
}

func (r *Router) View() string {
	if r.router == nil {
		r.router = mainmenu.NewMainMenu(r)
		r.router.Init()
	}

	flexHeight := r.height / 10

	viewBuilder := strings.Builder{}

	title := lipgloss.PlaceVertical(flexHeight*2, lipgloss.Bottom, components.GetAppTitle())
	viewBuilder.WriteString(title)

	body := lipgloss.PlaceVertical(flexHeight*8, lipgloss.Center, r.router.View())
	body = lipgloss.JoinVertical(lipgloss.Center, components.GetAppTitle(), r.router.View())
	viewBuilder.WriteString(body)

	return viewBuilder.String()
}
