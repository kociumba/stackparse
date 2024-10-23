package stackparse

import "github.com/charmbracelet/lipgloss"

// Theme defines the styling for different components
type Theme struct {
	Base      lipgloss.Style
	Goroutine lipgloss.Style
	Function  lipgloss.Style
	Args      lipgloss.Style
	File      lipgloss.Style
	Line      lipgloss.Style
	CreatedBy lipgloss.Style
	Repeat    lipgloss.Style
}

// DefaultTheme returns the default styling theme
func DefaultTheme() *Theme {
	return &Theme{
		Base: lipgloss.NewStyle().PaddingLeft(2),
		Goroutine: lipgloss.NewStyle().
			Bold(true).
			// Foreground(lipgloss.Color("#00ADD8")).
			Foreground(lipgloss.Color("#ed8796")).
			MarginTop(1).
			MarginBottom(1),
		Function: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0c6c6")),
		Args: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7dc4e4")).
			PaddingLeft(2),
		File: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f5a97f")),
		Line: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eed49f")),
		// PaddingLeft(2),
		CreatedBy: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ee99a0")),
		Repeat: lipgloss.NewStyle().
			Italic(true).
			Faint(true),
	}
}
