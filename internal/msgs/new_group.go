package msgs

import tea "github.com/charmbracelet/bubbletea"

type newGroup struct{}

var signalNewGroup = &newGroup{}

func NewGroup() tea.Msg {
	return signalNewGroup
}
