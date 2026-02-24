package tui

import "github.com/charmbracelet/lipgloss"

var (
	cyan    = lipgloss.Color("86")
	gray    = lipgloss.Color("241")
	white   = lipgloss.Color("255")
	dimGray = lipgloss.Color("238")

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(0, 1).
			Width(22)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(cyan).
				Padding(0, 1).
				Width(22)

	cardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	selectedCardTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(cyan)

	previewStyle = lipgloss.NewStyle().
			Foreground(gray)

	paneCountStyle = lipgloss.NewStyle().
			Foreground(dimGray).
			Italic(true)

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
