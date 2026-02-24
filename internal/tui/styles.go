package tui

import "github.com/charmbracelet/lipgloss"

var (
	cyan    = lipgloss.Color("86")
	gray    = lipgloss.Color("241")
	white   = lipgloss.Color("255")

	cardBase = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(0, 1)

	selectedCardBase = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(cyan).
				Padding(0, 1)

	previewStyle = lipgloss.NewStyle().
			Foreground(gray)

	selectedPreviewStyle = lipgloss.NewStyle().
				Foreground(cyan)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cyan).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(gray).
			Padding(1, 2)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(white).
			Bold(true)
)
