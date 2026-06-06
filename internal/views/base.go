package views

import tea "github.com/charmbracelet/bubbletea"

type base struct {
	width  int
	height int
}

func (b *base) Init() tea.Cmd {
	return nil
}

func (b *base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (b *base) View() string {
	return ""
}

func (b *base) SetSize(width, height int) {
	b.width = width
	b.height = height
}
