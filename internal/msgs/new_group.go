package msgs

import tea "github.com/charmbracelet/bubbletea"

func PublishMessage[T any](msg T) func() tea.Msg {
	return func() tea.Msg {
		return msg
	}
}

type NewGroup struct{}

type DeleteGroup struct {
}
