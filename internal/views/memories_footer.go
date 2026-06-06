package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type memoriesFooter struct {
	*footer
	repository memory.Repository
	memories   []*memory.Memory
	cursor     *clicomponents.CursorList
}

func newMemoriesFooter(groupId int) *memoriesFooter {
	return &memoriesFooter{
		footer: &footer{
			Options: &clicomponents.CursorList{
				Items: []string{"(n) Nova memória", "(e) Editar", "(d -> y) Deletar"},
				KeyActions: []clicomponents.OnKeyFn{
					clicomponents.OnKeys(func() (tea.Model, tea.Cmd) {
						return nil, msgs.PublishMessage(msgs.NewMemory{})
					}, "n"),
					clicomponents.OnKeys(func() (tea.Model, tea.Cmd) {
						return nil, msgs.PublishMessage(msgs.EditMemory{})
					}, "e"),
					clicomponents.OnKeys(func() (tea.Model, tea.Cmd) {
						return nil, msgs.PublishMessage(msgs.DeleteMemory{})
					}, "d"),
				},
			},
		},
		repository: memory.GetRepositorySqlLite(nil),
	}
}

func (m *memoriesFooter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.footer.Update(msg)
	if model != nil {
		return model, cmd
	}

	return m, nil
}
