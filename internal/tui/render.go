package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/primefaces/tyle/internal/layout"
)

func renderGrid(layouts []layout.Layout, cursor int, cols int) string {
	var rows []string

	for i := 0; i < len(layouts); i += cols {
		end := i + cols
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
	pvStyle := previewStyle

	if isSelected {
		style = selectedCardStyle
		pvStyle = selectedPreviewStyle
	}

	preview := pvStyle.Render(strings.Join(l.Preview, "\n"))
	return style.Render(preview)
}
