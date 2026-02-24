package layout

import (
	"fmt"
	"strings"
	"unicode"
)

func Slugify(name string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash && b.Len() > 0 {
			b.WriteByte('-')
			prevDash = true
		}
	}
	s := b.String()
	return strings.TrimRight(s, "-")
}

func GenerateLayout(name string, rowsPerCol []int) Layout {
	if len(rowsPerCol) == 0 {
		rowsPerCol = []int{1}
	}

	totalPanes := 0
	for _, r := range rowsPerCol {
		totalPanes += r
	}

	steps := generateSteps(rowsPerCol)
	preview := generatePreview(rowsPerCol)
	desc := generateDescription(rowsPerCol)

	return Layout{
		ID:          Slugify(name),
		Name:        name,
		Description: desc,
		Preview:     preview,
		Steps:       steps,
		PaneCount:   totalPanes,
		FinalFocus:  Previous,
	}
}

func generateSteps(rowsPerCol []int) []LayoutStep {
	var steps []LayoutStep

	numCols := len(rowsPerCol)

	for i := 0; i < numCols-1; i++ {
		steps = append(steps, LayoutStep{Action: ActionSplit, Direction: Right})
	}

	for i := numCols - 1; i >= 0; i-- {
		rows := rowsPerCol[i]
		for j := 0; j < rows-1; j++ {
			steps = append(steps, LayoutStep{Action: ActionSplit, Direction: Down})
		}
		if i > 0 {
			for j := 0; j < rows; j++ {
				steps = append(steps, LayoutStep{Action: ActionFocus, Direction: Previous})
			}
		}
	}

	for j := 0; j < rowsPerCol[0]-1; j++ {
		steps = append(steps, LayoutStep{Action: ActionFocus, Direction: Previous})
	}

	steps = append(steps, LayoutStep{Action: ActionEqualize})

	return steps
}

func generatePreview(rowsPerCol []int) []string {
	numCols := len(rowsPerCol)

	maxRows := 0
	for _, r := range rowsPerCol {
		if r > maxRows {
			maxRows = r
		}
	}

	targetWidth := 15
	if numCols > 3 {
		targetWidth = numCols*4 + 1
	}

	innerWidth := targetWidth - numCols - 1
	colWidths := make([]int, numCols)
	base := innerWidth / numCols
	remainder := innerWidth % numCols
	for i := 0; i < numCols; i++ {
		colWidths[i] = base
		if i < remainder {
			colWidths[i]++
		}
	}

	gridHeight := maxRows * 2
	if gridHeight > 6 {
		gridHeight = 6
	}
	if gridHeight < 2 {
		gridHeight = 2
	}

	rowHeights := distributeHeight(gridHeight, maxRows)

	paneLabel := 'A'
	labels := make([][]string, numCols)
	for c := 0; c < numCols; c++ {
		labels[c] = make([]string, rowsPerCol[c])
		for r := 0; r < rowsPerCol[c]; r++ {
			labels[c][r] = string(paneLabel)
			paneLabel++
		}
	}

	colRowBoundaries := make([][]int, numCols)
	for c := 0; c < numCols; c++ {
		colRowBoundaries[c] = computeRowBoundaries(rowHeights, maxRows, rowsPerCol[c])
	}

	totalLines := gridHeight + 1
	var lines []string

	for lineIdx := 0; lineIdx <= gridHeight; lineIdx++ {
		var b strings.Builder

		for c := 0; c < numCols; c++ {
			isHBorder := false
			for _, boundary := range colRowBoundaries[c] {
				if lineIdx == boundary {
					isHBorder = true
					break
				}
			}

			isTop := lineIdx == 0
			isBottom := lineIdx == gridHeight

			if c == 0 {
				if isTop {
					b.WriteString("┌")
				} else if isBottom {
					b.WriteString("└")
				} else if isHBorder {
					b.WriteString("├")
				} else {
					b.WriteString("│")
				}
			}

			if isTop || isBottom || isHBorder {
				b.WriteString(strings.Repeat("─", colWidths[c]))
			} else {
				paneRow, labelLine := findPaneInfo(lineIdx, colRowBoundaries[c], rowHeights, maxRows, rowsPerCol[c])
				label := ""
				if paneRow >= 0 && paneRow < len(labels[c]) && labelLine {
					label = labels[c][paneRow]
				}
				b.WriteString(centerPad(label, colWidths[c]))
			}

			if c < numCols-1 {
				leftHBorder := isHBorder || isTop || isBottom
				rightIsHBorder := false
				for _, boundary := range colRowBoundaries[c+1] {
					if lineIdx == boundary {
						rightIsHBorder = true
						break
					}
				}
				rightHBorder := rightIsHBorder || isTop || isBottom

				if isTop {
					b.WriteString("┬")
				} else if isBottom {
					b.WriteString("┴")
				} else if leftHBorder && rightHBorder {
					b.WriteString("┼")
				} else if leftHBorder {
					b.WriteString("┤")
				} else if rightHBorder {
					b.WriteString("├")
				} else {
					b.WriteString("│")
				}
			} else {
				if isTop {
					b.WriteString("┐")
				} else if isBottom {
					b.WriteString("┘")
				} else if isHBorder {
					b.WriteString("┤")
				} else {
					b.WriteString("│")
				}
			}
		}

		lines = append(lines, b.String())
	}

	_ = totalLines
	return lines
}

func distributeHeight(gridHeight int, maxRows int) []int {
	heights := make([]int, maxRows)
	perRow := gridHeight / maxRows
	remainder := gridHeight % maxRows
	for i := 0; i < maxRows; i++ {
		heights[i] = perRow
		if i >= maxRows-remainder {
			heights[i]++
		}
	}
	return heights
}

func computeRowBoundaries(rowHeights []int, maxRows int, colRows int) []int {
	boundaries := []int{0}

	if colRows == 1 {
		total := 0
		for _, h := range rowHeights {
			total += h
		}
		boundaries = append(boundaries, total)
		return boundaries
	}

	if colRows == maxRows {
		pos := 0
		for i := 0; i < maxRows; i++ {
			pos += rowHeights[i]
			boundaries = append(boundaries, pos)
		}
		return boundaries
	}

	totalHeight := 0
	for _, h := range rowHeights {
		totalHeight += h
	}

	subHeights := distributeHeight(totalHeight, colRows)
	pos := 0
	for i := 0; i < colRows; i++ {
		pos += subHeights[i]
		boundaries = append(boundaries, pos)
	}
	return boundaries
}

func findPaneInfo(lineIdx int, boundaries []int, rowHeights []int, maxRows int, colRows int) (int, bool) {
	for i := 0; i < len(boundaries)-1; i++ {
		top := boundaries[i]
		bottom := boundaries[i+1]
		if lineIdx > top && lineIdx < bottom {
			mid := (top + bottom) / 2
			if bottom-top > 2 {
				return i, lineIdx == mid
			}
			return i, lineIdx == top+1
		}
	}
	return -1, false
}

func centerPad(label string, width int) string {
	if len(label) >= width {
		return label[:width]
	}
	totalPad := width - len(label)
	left := totalPad / 2
	right := totalPad - left
	return strings.Repeat(" ", left) + label + strings.Repeat(" ", right)
}

func generateDescription(rowsPerCol []int) string {
	numCols := len(rowsPerCol)
	totalPanes := 0
	for _, r := range rowsPerCol {
		totalPanes += r
	}

	allSame := true
	for _, r := range rowsPerCol {
		if r != rowsPerCol[0] {
			allSame = false
			break
		}
	}

	if numCols == 1 && rowsPerCol[0] == 1 {
		return "Single pane"
	}

	if allSame && rowsPerCol[0] == 1 {
		return fmt.Sprintf("%d equal columns", numCols)
	}

	if numCols == 1 {
		return fmt.Sprintf("%d rows in a single column", rowsPerCol[0])
	}

	if allSame {
		return fmt.Sprintf("%dx%d grid with %d panes", numCols, rowsPerCol[0], totalPanes)
	}

	parts := make([]string, numCols)
	for i, r := range rowsPerCol {
		if r == 1 {
			parts[i] = "1 row"
		} else {
			parts[i] = fmt.Sprintf("%d rows", r)
		}
	}
	return fmt.Sprintf("%d columns (%s), %d panes total", numCols, strings.Join(parts, ", "), totalPanes)
}
