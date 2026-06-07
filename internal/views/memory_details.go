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
	memory      *memory.Memory
	relations   [3][]*memory.Memory
	repository  memory.Repository
	cursors     [3]*clicomponents.CursorList
	activeLevel int
	err         error
	width       int
	height      int
	footer      *footer
}

func NewMemoryDetails(mem *memory.Memory) *MemoryDetails {
	return &MemoryDetails{
		memory:     mem,
		repository: memory.GetRepositorySqlLite(nil),
		cursors: [3]*clicomponents.CursorList{
			{RenderSize: 10},
			{RenderSize: 10},
			{RenderSize: 10},
		},
		footer: &footer{
			Options: &clicomponents.CursorList{
				Items: []string{
					"LEFT/RIGHT: Level",
					"UP/DOWN: Navigate",
					"ENTER: Open",
					"ESC/Q: Return",
				},
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
		case "left":
			m.moveActiveLevel(-1)
			return m, nil
		case "right":
			m.moveActiveLevel(1)
			return m, nil
		case "enter":
			selected := m.selectedRelation()
			if selected == nil {
				return m, nil
			}

			detailsView := NewMemoryDetails(selected)
			detailsView.SetSize(m.width, m.height)
			return GetRootInstance().ReplaceRoute(detailsView)
		}
	case *memoryRelationsLoaded:
		m.relations = msg.relations
		m.err = msg.err
		m.activeLevel = firstPopulatedLevel(m.relations)

		for level := range m.cursors {
			m.cursors[level].Cursor = 0
		}
	}

	model, cmd := m.cursors[m.activeLevel].Update(msg)
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
	contentWidth := m.width * 2 / 5
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
		Render(m.renderRelations(relationsBodyWidth))
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
	relations, err := m.getRelationLevels(context.Background())
	return &memoryRelationsLoaded{
		relations: relations,
		err:       err,
	}
}

func (m *MemoryDetails) renderRelations(width int) string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if relationCount(m.relations) == 0 {
		return "No relations"
	}

	columnWidth := max(width/3, 1)
	columns := make([]string, len(m.relations))

	for level, relations := range m.relations {
		items := make([]string, len(relations))
		for index, relation := range relations {
			items[index] = relation.Name
		}
		m.cursors[level].Items = items

		title := fmt.Sprintf("Level %d", level+1)
		if level == m.activeLevel {
			title = "> " + title
		}

		columns[level] = lipgloss.NewStyle().
			Width(columnWidth).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					"",
					m.cursors[level].View(),
				),
			)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (m *MemoryDetails) renderHeader(title string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Background(components.ColorMain).
		Foreground(components.ColorMainContrast).
		Render(title)
}

type memoryRelationsLoaded struct {
	relations [3][]*memory.Memory
	err       error
}

func (m *MemoryDetails) getRelationLevels(ctx context.Context) ([3][]*memory.Memory, error) {
	var levels [3][]*memory.Memory
	visited := map[int]bool{m.memory.ID: true}
	frontier := []*memory.Memory{m.memory}

	for level := range levels {
		next := make([]*memory.Memory, 0)

		for _, current := range frontier {
			relations, err := m.repository.GetRelations(ctx, current)
			if err != nil {
				return levels, err
			}

			for _, relation := range relations {
				if visited[relation.ID] {
					continue
				}

				visited[relation.ID] = true
				levels[level] = append(levels[level], relation)
				next = append(next, relation)
			}
		}

		frontier = next
	}

	return levels, nil
}

func (m *MemoryDetails) moveActiveLevel(delta int) {
	next := m.activeLevel + delta
	if next < 0 || next >= len(m.relations) {
		return
	}

	m.activeLevel = next
}

func (m *MemoryDetails) selectedRelation() *memory.Memory {
	relations := m.relations[m.activeLevel]
	cursor := m.cursors[m.activeLevel].Cursor
	if len(relations) == 0 || cursor >= len(relations) {
		return nil
	}

	return relations[cursor]
}

func firstPopulatedLevel(levels [3][]*memory.Memory) int {
	for level, relations := range levels {
		if len(relations) > 0 {
			return level
		}
	}

	return 0
}

func relationCount(levels [3][]*memory.Memory) int {
	count := 0
	for _, relations := range levels {
		count += len(relations)
	}

	return count
}
