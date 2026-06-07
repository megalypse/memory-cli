package views

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memory"
)

const (
	memoryRelationLevelCount = 3

	memoryDetailsWidthRatioNumerator   = 2
	memoryDetailsWidthRatioDenominator = 5

	memoryDetailsColumnGap     = 2
	memoryRelationsBorderWidth = 1
	memoryRelationsLeftPadding = 2
	memoryRelationsColumnInset = memoryRelationsBorderWidth + memoryRelationsLeftPadding
	memoryRelationsColumnCount = memoryRelationLevelCount

	memoryDetailsHeaderHeight   = 1
	memoryDetailsHeaderSpacing  = 1
	memoryRelationsTitleHeight  = 1
	memoryRelationsTitleSpacing = 1

	memoryRelationsCursorCenterRows   = 1
	memoryRelationsCursorSides        = 2
	memoryRelationsMinimumColumnWidth = 1
	memoryRelationsMinimumRenderSize  = 1
	memoryRelationDisplayOffset       = 1

	previousRelationLevel = -1
	nextRelationLevel     = 1
)

type MemoryDetails struct {
	memory      *memory.Memory
	relations   [memoryRelationLevelCount][]*memory.Memory
	repository  memory.Repository
	cursors     [memoryRelationLevelCount]*clicomponents.CursorList
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
		cursors: [memoryRelationLevelCount]*clicomponents.CursorList{
			{},
			{},
			{},
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
			m.moveActiveLevel(previousRelationLevel)
			return m, nil
		case "right":
			m.moveActiveLevel(nextRelationLevel)
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
	contentWidth := m.width *
		memoryDetailsWidthRatioNumerator /
		memoryDetailsWidthRatioDenominator
	relationsWidth := m.width - contentWidth

	content := fmt.Sprintf("Name: %s\n\nContent:\n%s", m.memory.Name, m.memory.Content)
	contentBodyWidth := max(contentWidth-memoryDetailsColumnGap, 0)
	contentBodyHeight := max(
		m.height-memoryDetailsHeaderHeight-memoryDetailsHeaderSpacing,
		0,
	)
	contentBody := lipgloss.NewStyle().
		Width(contentBodyWidth).
		Height(contentBodyHeight).
		Render(content)
	contentColumn := lipgloss.NewStyle().
		Height(m.height).
		PaddingRight(memoryDetailsColumnGap).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.renderHeader("Details", contentBodyWidth),
				"",
				contentBody,
			),
		)

	relationsBodyWidth := max(relationsWidth-memoryRelationsColumnInset, 0)
	relationsBodyHeight := max(
		m.height-memoryDetailsHeaderHeight-memoryDetailsHeaderSpacing,
		0,
	)
	relationsBody := lipgloss.NewStyle().
		Width(relationsBodyWidth).
		Height(relationsBodyHeight).
		Render(m.renderRelations(relationsBodyWidth, relationsBodyHeight))
	relationsColumn := lipgloss.NewStyle().
		Height(m.height).
		BorderLeft(true).
		PaddingLeft(memoryRelationsLeftPadding).
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

func (m *MemoryDetails) renderRelations(width, height int) string {
	if m.err != nil {
		return lipgloss.NewStyle().Height(height).Render(fmt.Sprintf("Error: %v", m.err))
	}

	if relationCount(m.relations) == 0 {
		return lipgloss.NewStyle().Height(height).Render("No relations")
	}

	columnWidth := max(
		width/memoryRelationsColumnCount,
		memoryRelationsMinimumColumnWidth,
	)
	listHeight := max(
		height-memoryRelationsTitleHeight-memoryRelationsTitleSpacing,
		0,
	)
	renderSize := max(
		(listHeight-memoryRelationsCursorCenterRows)/memoryRelationsCursorSides,
		memoryRelationsMinimumRenderSize,
	)
	columns := make([]string, len(m.relations))

	for level, relations := range m.relations {
		items := make([]string, len(relations))
		for index, relation := range relations {
			items[index] = relation.Name
		}
		m.cursors[level].Items = items
		m.cursors[level].RenderSize = renderSize

		title := fmt.Sprintf("%dº", level+memoryRelationDisplayOffset)
		list := lipgloss.NewStyle().
			Height(listHeight).
			Render(strings.Join(items, "\n"))
		if level == m.activeLevel {
			list = lipgloss.NewStyle().
				Height(listHeight).
				Render(m.cursors[level].View())
		}

		columns[level] = lipgloss.NewStyle().
			Width(columnWidth).
			Height(height).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					"",
					list,
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
	relations [memoryRelationLevelCount][]*memory.Memory
	err       error
}

func (m *MemoryDetails) getRelationLevels(
	ctx context.Context,
) ([memoryRelationLevelCount][]*memory.Memory, error) {
	var levels [memoryRelationLevelCount][]*memory.Memory
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

func firstPopulatedLevel(
	levels [memoryRelationLevelCount][]*memory.Memory,
) int {
	for level, relations := range levels {
		if len(relations) > 0 {
			return level
		}
	}

	return 0
}

func relationCount(levels [memoryRelationLevelCount][]*memory.Memory) int {
	count := 0
	for _, relations := range levels {
		count += len(relations)
	}

	return count
}
