package views

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	tea.Model

	SetSize(width, height int)
	RenderFooter() string
}
