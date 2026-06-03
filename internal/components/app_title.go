package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func GetAppTitle() string {
	title := figure.NewFigure("Memory", "", true).String()
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMain)).Render(title)
}
