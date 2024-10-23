package stackparse

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// Formatter handles the formatting of parsed stack traces
type Formatter struct {
	config *Config
}

// NewFormatter creates a new Formatter instance
func NewFormatter(config *Config) *Formatter {
	return &Formatter{config: config}
}

// Format converts a StackTrace into a formatted string
func (f *Formatter) Format(trace *StackTrace) string {
	if !f.config.Colorize {
		f.disableColors()
	}

	var result []string

	// Add goroutine header
	if trace.GoroutineID != "" {
		header := fmt.Sprintf("Goroutine %s: %s",
			trace.GoroutineID, trace.GoroutineState)
		result = append(result, f.config.Theme.Goroutine.Render(header))
	}

	// Format entries
	functionCounts := f.countFunctions(trace.Entries)
	maxWidths := f.calculateMaxWidths(trace.Entries)

	for _, entry := range trace.Entries {
		formattedEntry := f.formatEntry(entry, functionCounts, maxWidths)
		result = append(result, formattedEntry)
	}

	return strings.Join(result, "\n")
}

// MaxWidths holds the maximum widths for different components
type MaxWidths struct {
	Function int
	File     int
	Line     int
}

func (f *Formatter) calculateMaxWidths(entries []StackEntry) MaxWidths {
	var maxWidths MaxWidths
	for _, entry := range entries {
		maxWidths.Function = max(maxWidths.Function, lipgloss.Width(entry.FunctionName))
		maxWidths.File = max(maxWidths.File, lipgloss.Width(entry.File))
		maxWidths.Line = max(maxWidths.Line, lipgloss.Width(entry.Line))
	}
	return maxWidths
}

func (f *Formatter) countFunctions(entries []StackEntry) map[string]int {
	counts := make(map[string]int)
	for _, entry := range entries {
		counts[entry.FunctionName]++
	}
	return counts
}

func (f *Formatter) formatEntry(entry StackEntry, functionCounts map[string]int, maxWidths MaxWidths) string {
	var currentTree *tree.Tree

	if functionCounts[entry.FunctionName] <= 1 {
		if entry.IsCreatedBy {
			currentTree = tree.New().Root(
				f.config.Theme.CreatedBy.Render("Created by: ") +
					entry.FunctionName,
			)
		} else {
			currentTree = tree.New().Root(
				f.config.Theme.Function.Render("Function: ") +
					lipgloss.NewStyle().Width(maxWidths.Function).Render(entry.FunctionName) +
					f.config.Theme.Args.Render(fmt.Sprintf("(%s)", entry.Args)),
			)
		}

		// Add file location as a child
		fileInfo := lipgloss.JoinHorizontal(
			lipgloss.Left,
			f.config.Theme.File.Render("At: "),
			lipgloss.NewStyle().Width(maxWidths.File).Render(entry.File),
			f.config.Theme.Line.Render(fmt.Sprintf(" Line: %s", entry.Line)),
		)
		currentTree.Child(fileInfo)
	} else {
		repCount := f.config.Theme.Repeat.Render(
			fmt.Sprintf(" (repeated %d times)",
				functionCounts[entry.FunctionName]),
		)
		currentTree = tree.New().Root(
			f.config.Theme.Function.Render("Function: ") +
				lipgloss.NewStyle().Width(maxWidths.Function).Render(entry.FunctionName) +
				f.config.Theme.Args.Render(fmt.Sprintf("(%s)", entry.Args)) +
				repCount,
		)
	}

	return currentTree.String()
}

func (f *Formatter) disableColors() {
	f.config.Theme.Goroutine.UnsetForeground()
	f.config.Theme.Function.UnsetForeground()
	f.config.Theme.Args.UnsetForeground()
	f.config.Theme.File.UnsetForeground()
	f.config.Theme.Line.UnsetForeground()
	f.config.Theme.CreatedBy.UnsetForeground()
	f.config.Theme.Repeat.UnsetForeground()
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
