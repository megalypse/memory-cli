package views

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/megalypse/memory_cli/internal/msgs"
	"github.com/megalypse/memory_cli/internal/utils"
)

func NewMemoryGroupCreate() View {
	repository := memorygroup.GetRepositorySqlLite(nil)

	return &MemoryGroupCreate{
		inputList: &clicomponents.InputList{
			EditMode: true,
		},
		footer:     newMemoryGroupCreateFooter(),
		repository: repository,
	}
}

type MemoryGroupCreate struct {
	nameInput        *components.TextInput
	descriptionInput *components.TextArea
	inputList        *clicomponents.InputList
	width            int
	height           int

	footer *memoryGroupCreateFooter

	repository memorygroup.RepositoryMemoryGroup
}

func (m *MemoryGroupCreate) Init() tea.Cmd {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.Focus()

	descriptionInput := textarea.New()
	descriptionInput.Placeholder = "Description"

	m.nameInput = &components.TextInput{Model: nameInput}
	m.descriptionInput = &components.TextArea{Model: descriptionInput}

	m.inputList.Inputs = append(m.inputList.Inputs, m.nameInput)
	m.inputList.Inputs = append(m.inputList.Inputs, m.descriptionInput)

	return nil
}

func (m *MemoryGroupCreate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, inputCmd := m.inputList.Update(msg)
	_, footerCmd := m.footer.Update(msg)

	m.footer.editMode = m.inputList.EditMode

	switch msg.(type) {
	case msgs.SaveGroup:
		groupName := m.nameInput.Value()
		groupDescription := m.nameInput.Value()

		if err := m.repository.Create(context.Background(), &memorygroup.MemoryGroup{
			Name:        groupName,
			Description: groupDescription,
		}); err != nil {
			return nil, msgs.PublishMessage(msgs.Err{Err: err})
		}

		_, cmd := GetRootInstance().PopRoute()
		return nil, tea.Batch(cmd, msgs.PublishMessage(msgs.UpdateGroups{}))
	}

	return m, tea.Batch(inputCmd, footerCmd)
}

func (m *MemoryGroupCreate) View() string {
	if m.nameInput == nil || m.descriptionInput == nil {
		return ""
	}

	return m.nameInput.View() + components.LineBreak + m.descriptionInput.View()
}

func (m *MemoryGroupCreate) RenderFooter() string {
	return m.footer.View()
}

func (m *MemoryGroupCreate) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(m.width, m.height)
}

func (m *MemoryGroupCreate) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *MemoryGroupCreate) GetHeight() int {
	return utils.CalcRatio(m.height, 0.9)
}
