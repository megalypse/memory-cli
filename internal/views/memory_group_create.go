package views

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewMemoryGroupCreate() View {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.Focus()

	descriptionInput := textarea.New()
	descriptionInput.Placeholder = "Description"

	return &MemoryGroupCreate{
		nameInput:        &nameInput,
		descriptionInput: &descriptionInput,
	}
}

type MemoryGroupCreate struct {
	nameInput        *textinput.Model
	descriptionInput *textarea.Model
}

func (m *MemoryGroupCreate) Init() tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryGroupCreate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryGroupCreate) View() string {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryGroupCreate) SetSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryGroupCreate) GetWidth() int {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryGroupCreate) GetHeight() int {
	//TODO implement me
	panic("implement me")
}
