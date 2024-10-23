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
			Foreground(lipgloss.Color("#00ADD8")).
			MarginTop(1).
			MarginBottom(1),
		Function: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")),
		Args: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61AFEF")).
			PaddingLeft(2),
		File: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C678DD")),
		Line: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5C07B")).
			PaddingLeft(2),
		CreatedBy: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")),
		Repeat: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")).
			Italic(true),
	}
}
