package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func GetAppTitle() string {
	title := figure.NewFigure("MEMORY", "", true).String()
	return lipgloss.NewStyle().Foreground(ColorMain).Render(title)
}
