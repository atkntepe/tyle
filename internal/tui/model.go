package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/primefaces/tyle/internal/layout"
)

const columns = 3

type Model struct {
	layouts   []layout.Layout
	cursor    int
	selected  *layout.Layout
	cancelled bool
	width     int
	height    int
}

func NewModel(layouts []layout.Layout) Model {
	return Model{
		layouts: layouts,
		cursor:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			m.selected = &m.layouts[m.cursor]
			return m, tea.Quit

		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right", "l":
			if m.cursor < len(m.layouts)-1 {
				m.cursor++
			}

		case "up", "k":
			if m.cursor-columns >= 0 {
				m.cursor -= columns
			}

		case "down", "j":
			if m.cursor+columns < len(m.layouts) {
				m.cursor += columns
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	header := headerStyle.Render("⊞ tyle")
	grid := renderGrid(m.layouts, m.cursor)
	help := helpStyle.Render(
		helpKeyStyle.Render("←→↑↓") + " navigate  " +
			helpKeyStyle.Render("enter") + " select  " +
			helpKeyStyle.Render("esc") + " cancel  " +
			helpKeyStyle.Render("q") + " quit",
	)
	return lipgloss.JoinVertical(lipgloss.Left, header, grid, help)
}

func (m Model) Selected() *layout.Layout {
	return m.selected
}

func (m Model) Cancelled() bool {
	return m.cancelled
}
