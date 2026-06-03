package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/views/mainmenu"
)

type Root struct {
	router tea.Model
	height int
	width  int
}

func (r *Root) GetWidth() int {
	return r.width
}

func (r *Root) GetHeight() int {
	return r.height
}

func (r *Root) Push(model tea.Model) (tea.Model, tea.Cmd) {
	return model, model.Init()
}

func (r *Root) Init() tea.Cmd {
	return nil
}

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
	}

	var cmd tea.Cmd
	r.router, cmd = r.router.Update(msg)

	return r, cmd
}

func (r *Root) View() string {
	if r.router == nil {
		r.router = mainmenu.NewMainMenu(r)
		r.router.Init()
	}

	title := lipgloss.Place(r.width, r.height, lipgloss.Center, lipgloss.Center, components.GetAppTitle())

	viewBuilder := strings.Builder{}
	viewBuilder.WriteString(title)
	viewBuilder.WriteString(components.LineSkip)
	viewBuilder.WriteString(r.router.View())

	return viewBuilder.String()
}
