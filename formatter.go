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
func (f *Formatter) Format(traces []*StackTrace) string {
	if !f.config.Colorize {
		f.disableColors()
	}

	var result []string

	for _, trace := range traces {
		result = append(result, f.formatTrace(trace))
	}

	return strings.Join(result, "\n\n")
}

func (f *Formatter) formatTrace(trace *StackTrace) string {
	var result []string

	// Format goroutine header
	header := fmt.Sprintf("Goroutine %s: %s",
		trace.GoroutineID, trace.GoroutineState)
	result = append(result, f.config.Theme.Goroutine.Render(header))

	// Calculate widths and counts
	functionCounts := f.countFunctions(trace.Entries)
	maxWidths := f.calculateMaxWidths(trace.Entries)

	// Format entries
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

	// Function name with potential truncation indicator
	displayName := entry.FunctionName
	if entry.FullName != "" && entry.FullName != entry.FunctionName {
		displayName = displayName + " ..."
	}

	if entry.IsCreatedBy {
		createdByInfo := fmt.Sprintf("Created by: %s (goroutine %s)",
			displayName, entry.CreatedByGoroutine)
		currentTree = tree.New().Root(
			f.config.Theme.CreatedBy.Render(createdByInfo),
		)
	} else {
		// Create base function node
		funcInfo := lipgloss.JoinHorizontal(
			lipgloss.Left,
			f.config.Theme.Function.Render(displayName),
			f.config.Theme.Args.Render(fmt.Sprintf("(%s)", entry.Args)),
		)

		// Add repeat count if needed
		if count := functionCounts[entry.FunctionName]; count > 1 {
			funcInfo += f.config.Theme.Repeat.Render(
				fmt.Sprintf(" (repeated %d times)", count),
			)
		}

		currentTree = tree.New().Root(funcInfo)
	}

	// Add file location as child node if available
	if entry.File != "" {
		fileInfo := lipgloss.JoinHorizontal(
			lipgloss.Left,
			f.config.Theme.File.Render(entry.File),
			f.config.Theme.Line.Render(fmt.Sprintf(":%s", entry.Line)),
		)
		if entry.Offset != "" {
			fileInfo += f.config.Theme.Line.Render(fmt.Sprintf(" +%s", entry.Offset))
		}
		currentTree.Child(fileInfo)
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
