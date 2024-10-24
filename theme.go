package stackparse

import "github.com/charmbracelet/lipgloss"

// StyleDisabler represents types that can disable their styles
type StyleDisabler interface {
	DisableStyles()
}

// Theme defines the styling for different components
//
// You can provide you own theme to stackparse, and all of the output will be styled accordingly
// Setting a custom DisableStylesFunc is not required but if you want to preserve specific styling, you can do so.
// Usage would look something like this.
//
// Method 1: Using SetDisableStylesFunc:
//	theme := stackparse.DefaultTheme()
//	theme.SetDisableStylesFunc(func(t *stackparse.Theme) {
//  	// Custom implementation
//  	t.Base = t.Base.UnsetForeground()
//  	t.Function = t.Function.UnsetBold()
//  	// ... other custom unset operations
//	})
//
// Method 2: Creating a custom theme with custom disable function:
//	customTheme := &stackparse.Theme{
//  	// ... theme settings ...
//	}
//	customTheme.SetDisableStylesFunc(func(t *stackparse.Theme) {
//  	// Custom implementation
//	})
type Theme struct {
	Base      lipgloss.Style
	Goroutine lipgloss.Style
	Function  lipgloss.Style
	Args      lipgloss.Style
	File      lipgloss.Style
	Line      lipgloss.Style
	CreatedBy lipgloss.Style
	Repeat    lipgloss.Style

	// Optional custom disable function
	disableStylesFunc func(*Theme)
}

// DisableStyles disables all color and text styles (e.g. bold, underline, etc.) of a theme
//
// When providing your own theme, you can use this default or create your own DisableStyles to unset only specific styling.
//
// Keep in mind you need to reassign the styles after unsetting parts of them, for example:
//	func (t *Theme) DisableStyles() {
//		t.Base = t.Base.UnsetForeground()
//	}
//
func (t *Theme) DisableStyles() {
	if t.disableStylesFunc != nil {
		t.disableStylesFunc(t)
		return
	}

	// Default implementation that disables everything
	for _, field := range []*lipgloss.Style{
		&t.Base,
		&t.Goroutine,
		&t.Function,
		&t.Args,
		&t.File,
		&t.Line,
		&t.CreatedBy,
		&t.Repeat,
	} {
		*field = field.
			UnsetForeground().
			UnsetBackground().
			UnsetBold().
			UnsetFaint().
			UnsetItalic().
			UnsetBlink()
	}
}

// SetDisableStylesFunc allows setting a custom DisableStyles implementation
//
// For more info read the documentation of
//	Theme{}
func (t *Theme) SetDisableStylesFunc(f func(*Theme)) {
	t.disableStylesFunc = f
}

// DefaultTheme returns the default styling theme
func DefaultTheme() *Theme {
	return &Theme{
		Base: lipgloss.NewStyle().PaddingLeft(2),
		Goroutine: lipgloss.NewStyle().
			Bold(true).
			// Foreground(lipgloss.Color("#00ADD8")).
			Foreground(lipgloss.Color("#ed8796")).
			MarginTop(1),
		// MarginBottom(1),
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
