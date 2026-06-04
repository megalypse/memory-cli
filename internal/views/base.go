package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/memory_cli/internal/utils"
)

type base struct {
	width  int
	height int

	widthRatio  float64
	heightRatio float64
}

func (b *base) GetWidth() int {
	return utils.CalcRatio(b.width, b.widthRatio)
}

func (b *base) GetHeight() int {
	return utils.CalcRatio(b.height, b.heightRatio)
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
