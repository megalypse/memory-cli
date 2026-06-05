package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/msgs"
)

func newMemoryGroupCreateFooterCmdMode() *footer {
	return &footer{
		Options: &clicomponents.CursorListVertical{
			Items: []string{
				"(s) Save",
				"(esc) Cancel",
			},
			KeyActions: []clicomponents.OnKeyFn{
				clicomponents.OnKey("s", func() (tea.Model, tea.Cmd) {
					return nil, msgs.PublishMessage(msgs.SaveGroup{})
				}),
				clicomponents.OnKey("esc", func() (tea.Model, tea.Cmd) {
					return nil, msgs.PublishMessage(msgs.CancelGroup{})
				}),
			},
		},
	}
}
