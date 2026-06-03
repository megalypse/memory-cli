package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/views/mainmenu"
)

type Root struct {
	router tea.Model
}

func (r *Root) Init() tea.Cmd {
	return nil
}

func (r *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	r.router, cmd = r.router.Update(msg)
	return r, cmd
}

func (r *Root) View() string {
	if r.router == nil {
		r.router = &mainmenu.MainMenu{}
	}

	return r.router.View()
}
