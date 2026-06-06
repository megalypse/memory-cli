package views

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	tea.Model

	SetSize(width, height int)
	GetWidth() int
	GetHeight() int
	RenderFooter() string
}
