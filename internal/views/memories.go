package views

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/megalypse/go-cli-components/clicomponents"
	"github.com/megalypse/memory_cli/internal/components"
	"github.com/megalypse/memory_cli/internal/memory"
	"github.com/megalypse/memory_cli/internal/msgs"
)

type Memories struct {
	memoryGroupId int
	width, height int
	memories      []*memory.Memory
	repository    memory.Repository
	err           error
	cursor        *clicomponents.CursorList
	searchInput   textinput.Model
	searching     bool
	footer        *memoriesFooter
}

func NewMemories(groupId int) *Memories {
	repository := memory.GetRepositorySqlLite(nil)
	searchInput := textinput.New()
	searchInput.Prompt = "SEARCH: "
	searchInput.Placeholder = "TITLE OR CONTENT"

	return &Memories{
		memoryGroupId: groupId,
		repository:    repository,
		cursor: &clicomponents.CursorList{
			RenderSize: 10,
		},
		searchInput: searchInput,
		footer:      newMemoriesFooter(groupId),
	}
}

func (m *Memories) Init() tea.Cmd {
	return func() tea.Msg {
		return m.loadMemories()
	}
}

func (m *Memories) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if result, ok := msg.(memoriesQueryResult); ok {
		if result.query != m.searchInput.Value() {
			return m, nil
		}

		m.err = result.err
		m.memories = result.memories
		m.resetCursor()
		return m, nil
	}

	if m.searching {
		return m.updateSearch(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return GetRootInstance().PopRoute()

		case "q":
			return m, tea.Quit
		case "/":
			m.searching = true
			return m, m.searchInput.Focus()
		case "enter":
			if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
				selectedMemory := m.memories[m.cursor.Cursor]
				detailsView := NewMemoryDetails(selectedMemory)
				detailsView.SetSize(m.width, m.height)

				return GetRootInstance().PushRoute(detailsView)
			}
			return m, nil
		case "n":
			createView := NewMemoryCreate(m.memoryGroupId)

			return GetRootInstance().PushRoute(createView)
		}
	case msgs.NewMemory:
		return m, func() tea.Msg {
			return m.loadMemories()
		}
	case msgs.EditMemory:
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, nil
		}
	case msgs.DeleteMemory:
		if len(m.memories) > 0 && m.cursor.Cursor < len(m.memories) {
			return m, func() tea.Msg {
				return m.loadMemories()
			}
		}
	}

	model, cmd := m.cursor.Update(msg)
	if model != nil {
		return model, cmd
	}

	model, cmd = m.footer.Update(msg)
	if model != nil {
		return model, cmd
	}

	return m, nil
}

func (m *Memories) View() string {
	if m.err != nil {
		return fmt.Sprintf("ERROR LOADING MEMORIES: %v", m.err)
	}

	items := make([]string, len(m.memories))
	for i, memory := range m.memories {
		items[i] = fmt.Sprintf("%s", memory.Name)
	}
	m.cursor.Items = items

	return m.cursor.View()
}

func (m *Memories) RenderFooter() string {
	if m.searching {
		return lipgloss.NewStyle().
			Width(m.width).
			Background(components.ColorMain).
			Foreground(components.ColorMainContrast).
			Render(m.searchInput.View())
	}

	return m.footer.View()
}

func (m *Memories) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.footer.SetSize(width, height)
}

func (m *Memories) loadMemories() tea.Msg {
	memories, err := m.repository.GetAllByGroup(context.Background(), m.memoryGroupId)
	if err != nil {
		m.err = err
		return nil
	}

	m.memories = memories
	m.resetCursor()
	return nil
}

func (m *Memories) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.searching = false
			m.searchInput.Blur()
			m.searchInput.SetValue("")
			return m, func() tea.Msg {
				return m.loadMemories()
			}
		case "enter":
			m.searching = false
			m.searchInput.Blur()
			return m, nil
		}
	}

	var inputCmd tea.Cmd
	m.searchInput, inputCmd = m.searchInput.Update(msg)
	query := m.searchInput.Value()

	return m, tea.Batch(inputCmd, m.queryMemories(query))
}

func (m *Memories) queryMemories(query string) tea.Cmd {
	return func() tea.Msg {
		memories, err := m.repository.QueryMemories(
			context.Background(),
			m.memoryGroupId,
			query,
		)
		return memoriesQueryResult{
			query:    query,
			memories: memories,
			err:      err,
		}
	}
}

func (m *Memories) resetCursor() {
	if len(m.memories) == 0 {
		m.cursor.Cursor = 0
		return
	}

	if m.cursor.Cursor >= len(m.memories) {
		m.cursor.Cursor = len(m.memories) - 1
	}
}

type memoriesQueryResult struct {
	query    string
	memories []*memory.Memory
	err      error
}
