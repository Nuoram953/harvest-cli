package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#874BFE")
	SecondaryColor = lipgloss.Color("#7D56F4")
	AccentColor    = lipgloss.Color("#F25D94")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(PrimaryColor).
			Padding(0, 1)
)
