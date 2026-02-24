package tui

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"

	"github.com/primefaces/tyle/internal/layout"
)

type cardDimensions struct {
	width  int
	height int
}

func measureLayouts(layouts []layout.Layout) cardDimensions {
	maxW, maxH := 0, 0
	for _, l := range layouts {
		for _, line := range l.Preview {
			w := utf8.RuneCountInString(line)
			if w > maxW {
				maxW = w
			}
		}
		if len(l.Preview) > maxH {
			maxH = len(l.Preview)
		}
	}
	return cardDimensions{width: maxW + 2, height: maxH}
}

func renderGrid(layouts []layout.Layout, cursor int, cols int) string {
	dim := measureLayouts(layouts)
	var rows []string

	for i := 0; i < len(layouts); i += cols {
		end := i + cols
		if end > len(layouts) {
			end = len(layouts)
		}

		var cards []string
		for j := i; j < end; j++ {
			cards = append(cards, renderCard(layouts[j], j == cursor, dim))
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, cards...)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderCard(l layout.Layout, isSelected bool, dim cardDimensions) string {
	base := cardBase
	pvStyle := previewStyle

	if isSelected {
		base = selectedCardBase
		pvStyle = selectedPreviewStyle
	}

	style := base.Width(dim.width).Height(dim.height)
	preview := pvStyle.Render(strings.Join(l.Preview, "\n"))
	return style.Render(preview)
}
