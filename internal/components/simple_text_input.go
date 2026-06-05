package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SimpleTextInput struct {
	textinput.Model
}

func (s SimpleTextInput) Focus() tea.Cmd {
	return s.Model.Focus()
}

func (s SimpleTextInput) Blur() {
	s.Model.Blur()
}

func (s SimpleTextInput) Up() {
	// Não há navegação vertical em um campo de texto simples, mas o método precisa existir conforme a interface
	// Vamos mover para o início do campo de texto
	pos := s.Model.Position()
	if pos > 0 {
		s.Model.SetCursor(0)
	}
}

func (s SimpleTextInput) Down() {
	// Não há navegação vertical em um campo de texto simples, mas o método precisa existir conforme a interface
	// Vamos mover para o final do campo de texto
	pos := s.Model.Position()
	text := s.Model.Value()
	if pos < len(text) {
		s.Model.SetCursor(len(text))
	}
}

func (s SimpleTextInput) Left() {
	pos := s.Model.Position()
	if pos > 0 {
		s.Model.SetCursor(pos - 1)
	}
}

func (s SimpleTextInput) Right() {
	pos := s.Model.Position()
	text := s.Model.Value()
	if pos < len(text) {
		s.Model.SetCursor(pos + 1)
	}
}
