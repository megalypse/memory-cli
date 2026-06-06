package views

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
	"github.com/megalypse/memory_cli/internal/utils"
)

func NewMemoryCreate(groupId int) View {
	repository := memory.GetRepositorySqlLite(nil)

	return &MemoryCreate{
		memoryGroupId: groupId,
		inputList: &clicomponents.InputList{
			EditMode: true,
		},
		footer:     newMemoryCreateFooter(),
		repository: repository,
	}
}

type MemoryCreate struct {
	nameInput    *components.TextInput
	contentInput *components.TextArea
	inputList    *clicomponents.InputList
	width        int
	height       int

	footer *memoryCreateFooter

	memoryGroupId int

	repository memory.Repository
}

func (m *MemoryCreate) Init() tea.Cmd {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.Focus()

	contentInput := textarea.New()
	contentInput.Placeholder = "Content"

	m.nameInput = &components.TextInput{Model: nameInput}
	m.contentInput = &components.TextArea{Model: contentInput}

	m.inputList.Inputs = append(m.inputList.Inputs, m.nameInput)
	m.inputList.Inputs = append(m.inputList.Inputs, m.contentInput)

	return nil
}

func (m *MemoryCreate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, inputCmd := m.inputList.Update(msg)
	_, footerCmd := m.footer.Update(msg)

	m.footer.editMode = m.inputList.EditMode

	switch msg.(type) {
	case msgs.SaveMemory:
		memoryName := m.nameInput.Value()
		memoryContent := m.contentInput.Value()

		if err := m.repository.Create(context.Background(), &memory.Memory{
			GroupID: m.memoryGroupId,
			Name:    memoryName,
			Content: memoryContent,
		}); err != nil {
			return nil, msgs.PublishMessage(msgs.Err{Err: err})
		}

		_, cmd := GetRootInstance().PopRoute()
		return nil, tea.Batch(cmd, msgs.PublishMessage(msgs.NewMemory{}))
	}

	return m, tea.Batch(inputCmd, footerCmd)
}

func (m *MemoryCreate) View() string {
	if m.nameInput == nil || m.contentInput == nil {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.PlaceVertical(
			m.GetHeight(),
			lipgloss.Center,
			m.nameInput.View()+components.LineBreak+m.contentInput.View()),
		lipgloss.PlaceVertical(
			m.footer.GetHeight(),
			lipgloss.Center,
			m.footer.View()),
	)
}

func (m *MemoryCreate) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(m.width, m.height)
}

func (m *MemoryCreate) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *MemoryCreate) GetHeight() int {
	return utils.CalcRatio(m.height, 0.9)
}
