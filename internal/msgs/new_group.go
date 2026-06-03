package msgs

import tea "github.com/charmbracelet/bubbletea"

type NewGroup struct{}

func MsgNewGroup() tea.Msg {
	return NewGroup{}
}
