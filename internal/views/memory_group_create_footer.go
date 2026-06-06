package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/utils"
)

func newMemoryGroupCreateFooter() *memoryGroupCreateFooter {
	return &memoryGroupCreateFooter{
		footerEditMode: newMemoryGroupCreateFooterEditMode(),
		footerCmdMode:  newMemoryGroupCreateFooterCmdMode(),
		editMode:       true,
	}
}

type memoryGroupCreateFooter struct {
	base

	footerEditMode *footer
	footerCmdMode  *footer

	editMode bool
}

func (m *memoryGroupCreateFooter) Init() tea.Cmd {
	return nil
}

func (m *memoryGroupCreateFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editMode {
		return m.footerEditMode.Update(msg)
	}

	return m.footerCmdMode.Update(msg)
}

func (m *memoryGroupCreateFooter) View() string {
	if m.editMode {
		return m.footerEditMode.View()
	}

	return m.footerCmdMode.View()
}

func (m *memoryGroupCreateFooter) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footerEditMode.SetSize(width, height)
	m.footerCmdMode.SetSize(width, height)
}

func (m *memoryGroupCreateFooter) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *memoryGroupCreateFooter) GetHeight() int {
	return utils.CalcRatio(m.height, 0.1)
}
