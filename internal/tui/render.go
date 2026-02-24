package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/primefaces/tyle/internal/layout"
)

func renderGrid(layouts []layout.Layout, cursor int) string {
	var rows []string

	for i := 0; i < len(layouts); i += columns {
		end := i + columns
		if end > len(layouts) {
			end = len(layouts)
		}

		var cards []string
		for j := i; j < end; j++ {
			cards = append(cards, renderCard(layouts[j], j == cursor))
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, cards...)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderCard(l layout.Layout, isSelected bool) string {
	style := cardStyle
	nameStyle := cardTitleStyle

	if isSelected {
		style = selectedCardStyle
		nameStyle = selectedCardTitleStyle
	}

	name := nameStyle.Render(l.Name)
	preview := previewStyle.Render(strings.Join(l.Preview, "\n"))
	info := paneCountStyle.Render(fmt.Sprintf("%d panes", l.PaneCount))

	content := lipgloss.JoinVertical(lipgloss.Left, name, "", preview, "", info)

	return style.Render(content)
}
