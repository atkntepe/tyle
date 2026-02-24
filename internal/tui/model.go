package tui

import (
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/primefaces/tyle/internal/layout"
)

const maxColumns = 4

type Model struct {
	layouts   []layout.Layout
	cursor    int
	selected  *layout.Layout
	cancelled bool
	width     int
	height    int
	scroll    int
}

func NewModel(layouts []layout.Layout) Model {
	return Model{
		layouts: layouts,
		cursor:  0,
	}
}

func (m Model) cardOuterWidth() int {
	maxW := 0
	for _, l := range m.layouts {
		for _, line := range l.Preview {
			w := utf8.RuneCountInString(line)
			if w > maxW {
				maxW = w
			}
		}
	}
	return maxW + 4 + 1
}

func (m Model) cardOuterHeight() int {
	maxH := 0
	for _, l := range m.layouts {
		if len(l.Preview) > maxH {
			maxH = len(l.Preview)
		}
	}
	return maxH + 2
}

func (m Model) cols() int {
	if m.width <= 0 {
		return 3
	}
	ow := m.cardOuterWidth()
	if ow <= 0 {
		return 3
	}
	c := m.width / ow
	if c < 1 {
		return 1
	}
	if c > maxColumns {
		return maxColumns
	}
	return c
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.scroll = m.ensureVisible(m.cursor)
		return m, nil

	case tea.KeyMsg:
		cols := m.cols()

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
			if m.cursor-cols >= 0 {
				m.cursor -= cols
			}

		case "down", "j":
			if m.cursor+cols < len(m.layouts) {
				m.cursor += cols
			}
		}

		m.scroll = m.ensureVisible(m.cursor)
	}

	return m, nil
}

func (m Model) ensureVisible(cursor int) int {
	if m.height <= 0 {
		return 0
	}

	cols := m.cols()
	row := cursor / cols

	headerHeight := 3
	helpHeight := 3
	ch := m.cardOuterHeight()
	available := m.height - headerHeight - helpHeight
	visibleRows := available / ch
	if visibleRows < 1 {
		visibleRows = 1
	}

	scroll := m.scroll
	if row < scroll {
		scroll = row
	}
	if row >= scroll+visibleRows {
		scroll = row - visibleRows + 1
	}
	return scroll
}

func (m Model) View() string {
	header := headerStyle.Render("⊞ tyle")

	cols := m.cols()
	grid := renderGrid(m.layouts, m.cursor, cols)

	gridLines := strings.Split(grid, "\n")

	headerHeight := 3
	helpHeight := 3
	available := m.height - headerHeight - helpHeight
	if available < 1 {
		available = 1
	}

	ch := m.cardOuterHeight()
	startLine := m.scroll * ch
	if startLine > len(gridLines) {
		startLine = len(gridLines)
	}
	endLine := startLine + available
	if endLine > len(gridLines) {
		endLine = len(gridLines)
	}

	visibleGrid := strings.Join(gridLines[startLine:endLine], "\n")

	help := helpStyle.Render(
		helpKeyStyle.Render("←→↑↓") + " navigate  " +
			helpKeyStyle.Render("enter") + " select  " +
			helpKeyStyle.Render("esc") + " cancel",
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, visibleGrid, help)
}

func (m Model) Selected() *layout.Layout {
	return m.selected
}

func (m Model) Cancelled() bool {
	return m.cancelled
}
