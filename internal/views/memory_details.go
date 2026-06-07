package views

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memory"
)

type MemoryDetails struct {
	memory     *memory.Memory
	relations  []*memory.Memory
	repository memory.Repository
	cursor     *clicomponents.CursorList
	err        error
	width      int
	height     int
	footer     *footer
}

func NewMemoryDetails(mem *memory.Memory) *MemoryDetails {
	return &MemoryDetails{
		memory:     mem,
		repository: memory.GetRepositorySqlLite(nil),
		cursor: &clicomponents.CursorList{
			RenderSize: 10,
		},
		footer: &footer{
			Options: &clicomponents.CursorList{
				Items: []string{"UP/DOWN: Navigate", "ENTER: Open", "ESC/Q: Return"},
			},
		},
	}
}

func (m *MemoryDetails) Init() tea.Cmd {
	return m.loadRelations
}

func (m *MemoryDetails) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return GetRootInstance().PopRoute()
		case "enter":
			if len(m.relations) == 0 || m.cursor.Cursor >= len(m.relations) {
				return m, nil
			}

			detailsView := NewMemoryDetails(m.relations[m.cursor.Cursor])
			detailsView.SetSize(m.width, m.height)
			return GetRootInstance().ReplaceRoute(detailsView)
		}
	case memoryRelationsLoaded:
		m.relations = msg.relations
		m.err = msg.err
		m.cursor.Cursor = 0
	}

	model, cmd := m.cursor.Update(msg)
	if model != nil {
		return model, cmd
	}

	model, footerCmd := m.footer.Update(msg)
	if model != nil {
		return model, footerCmd
	}

	return m, tea.Batch(cmd, footerCmd)
}

func (m *MemoryDetails) View() string {
	contentWidth := m.width * 2 / 3
	relationsWidth := m.width - contentWidth

	content := fmt.Sprintf("Name: %s\n\nContent:\n%s", m.memory.Name, m.memory.Content)
	contentBodyWidth := max(contentWidth-2, 0)
	contentBody := lipgloss.NewStyle().
		Width(contentBodyWidth).
		Render(content)
	contentColumn := lipgloss.NewStyle().PaddingRight(2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderHeader("Details", contentBodyWidth),
			"",
			contentBody,
		),
	)

	relationsBodyWidth := max(relationsWidth-3, 0)
	relationsBody := lipgloss.NewStyle().
		Width(relationsBodyWidth).
		Render(m.renderRelations())
	relationsColumn := lipgloss.NewStyle().
		BorderLeft(true).
		PaddingLeft(2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.renderHeader("Relations", relationsBodyWidth),
				"",
				relationsBody,
			),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top, contentColumn, relationsColumn)
}

func (m *MemoryDetails) RenderFooter() string {
	return m.footer.View()
}

func (m *MemoryDetails) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(m.width, m.height)
}

func (m *MemoryDetails) loadRelations() tea.Msg {
	relations, err := m.repository.GetRelations(context.Background(), m.memory)
	return memoryRelationsLoaded{
		relations: relations,
		err:       err,
	}
}

func (m *MemoryDetails) renderRelations() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if len(m.relations) == 0 {
		return "No relations"
	}

	items := make([]string, len(m.relations))
	for index, relation := range m.relations {
		items[index] = relation.Name
	}
	m.cursor.Items = items

	return m.cursor.View()
}

func (m *MemoryDetails) renderHeader(title string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Background(components.ColorMain).
		Foreground(components.ColorMainContrast).
		Render(title)
}

type memoryRelationsLoaded struct {
	relations []*memory.Memory
	err       error
}
