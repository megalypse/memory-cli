package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TextInput é a interface que define os métodos necessários para um componente de entrada de texto
type TextInput interface {
	Focus() tea.Cmd
	Blur()
	Up()
	Down()
	Left()
	Right()
}
