package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Memories struct {
	memoryGroupId int
	width, height int
}

func (m *Memories) Init() tea.Cmd {
	return nil
}

func (m *Memories) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *Memories) View() string {
	return ""
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Memories) GetWidth() int {
	return m.width
}

func (m *Memories) GetHeight() int {
	return m.height
}
