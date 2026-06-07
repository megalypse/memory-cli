package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/msgs"
)

func newMemoryCreateFooterCmdMode() *footer {
	return &footer{
		Options: &clicomponents.CursorList{
			Items: []string{
				"(S) SAVE",
				"(ESC) CANCEL",
			},
			KeyActions: []clicomponents.OnKeyFn{
				clicomponents.OnKey("s", func() (tea.Model, tea.Cmd) {
					return nil, msgs.PublishMessage(msgs.SaveMemory{})
				}),
				clicomponents.OnKey("esc", func() (tea.Model, tea.Cmd) {
					return GetRootInstance().PopRoute()
				}),
			},
		},
	}
}
