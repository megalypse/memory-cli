package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/utils"
)

func newMemoryCreateFooter() *memoryCreateFooter {
	return &memoryCreateFooter{
		footerEditMode: newMemoryCreateFooterEditMode(),
		footerCmdMode:  newMemoryCreateFooterCmdMode(),
		editMode:       true,
	}
}

type memoryCreateFooter struct {
	base

	footerEditMode *footer
	footerCmdMode  *footer

	editMode bool
}

func (m *memoryCreateFooter) Init() tea.Cmd {
	return nil
}

func (m *memoryCreateFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editMode {
		return m.footerEditMode.Update(msg)
	}

	return m.footerCmdMode.Update(msg)
}

func (m *memoryCreateFooter) View() string {
	if m.editMode {
		return m.footerEditMode.View()
	}

	return m.footerCmdMode.View()
}

func (m *memoryCreateFooter) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footerEditMode.SetSize(width, height)
	m.footerCmdMode.SetSize(width, height)
}

func (m *memoryCreateFooter) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *memoryCreateFooter) GetHeight() int {
	return utils.CalcRatio(m.height, 0.1)
}
