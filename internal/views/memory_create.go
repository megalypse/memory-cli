package views

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/keyterm"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
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
	nameInput.Placeholder = "NAME"
	nameInput.Focus()

	contentInput := textarea.New()
	contentInput.Placeholder = "CONTENT"

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

		if err := m.createMemory(context.Background(), memoryName, memoryContent); err != nil {
			return nil, msgs.PublishMessage(msgs.Err{Err: err})
		}

		_, cmd := GetRootInstance().PopRoute()
		return nil, tea.Batch(cmd, msgs.PublishMessage(msgs.NewMemory{}))
	}

	return m, tea.Batch(inputCmd, footerCmd)
}

func (m *MemoryCreate) createMemory(ctx context.Context, name, content string) error {
	terms := append(keyterm.Extract(name), keyterm.Extract(content)...)
	references := make([]*memory.Memory, 0, len(terms))
	referenceIDs := make(map[int]bool, len(terms))

	for _, term := range terms {
		found, err := m.repository.FindReferences(ctx, m.memoryGroupId, []string{term})
		if err != nil {
			return err
		}

		if len(found) == 0 || referenceIDs[found[0].ID] {
			continue
		}

		references = append(references, found[0])
		referenceIDs[found[0].ID] = true
	}

	newMemory := &memory.Memory{
		GroupID: m.memoryGroupId,
		Name:    name,
		Content: content,
	}

	if err := m.repository.Create(ctx, newMemory); err != nil {
		return err
	}

	if len(references) == 0 {
		return nil
	}

	return m.repository.LinkMemories(ctx, newMemory, references)
}

func (m *MemoryCreate) View() string {
	if m.nameInput == nil || m.contentInput == nil {
		return ""
	}

	return m.nameInput.View() + components.LineBreak + m.contentInput.View()
}

func (m *MemoryCreate) RenderFooter() string {
	return m.footer.View()
}

func (m *MemoryCreate) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(m.width, m.height)
}
