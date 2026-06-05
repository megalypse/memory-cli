package views

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/components"
	"github.com/megalypse/memory_cli/internal/utils"
)

func NewMemoryGroupCreate() View {
	return &MemoryGroupCreate{}
}

type MemoryGroupCreate struct {
	nameInput        *textinput.Model
	descriptionInput *textarea.Model
	inputList        components.InputList
	width            int
	height           int
}

func (m *MemoryGroupCreate) Init() tea.Cmd {
	// Initialize inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.Focus()

	descriptionInput := textarea.New()
	descriptionInput.Placeholder = "Description"

	m.nameInput = &nameInput
	m.descriptionInput = &descriptionInput

	m.inputList.Inputs = append(m.inputList.Inputs, &nameInput)
	m.inputList.Inputs = append(m.inputList.Inputs, &descriptionInput)

	return nil
}

func (m *MemoryGroupCreate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.inputList.Update(msg)
	if model != nil {
		return model, cmd
	}

	m.nameInput.Focus()

	return m, cmd
}

func (m *MemoryGroupCreate) View() string {
	if m.nameInput == nil || m.descriptionInput == nil {
		return ""
	}

	nameView := m.nameInput.View()
	descriptionView := m.descriptionInput.View()

	return nameView + "\n" + descriptionView
}

func (m *MemoryGroupCreate) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *MemoryGroupCreate) GetWidth() int {
	return utils.CalcRatio(m.width, 1)
}

func (m *MemoryGroupCreate) GetHeight() int {
	return utils.CalcRatio(m.height, 1)
}
