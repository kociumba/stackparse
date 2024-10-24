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
//
// Mostly used internally, but exposed for some edge cases
func NewFormatter(config *Config) *Formatter {
	return &Formatter{
		config: config,
	}
}

// Format converts a StackTrace into a formatted string
//
// If you are using this by itself without calling stackparse.Parse(), you should call this after the parser.
func (f *Formatter) Format(traces []*StackTrace) string {
	if !f.config.Colorize {
		f.config.Theme.DisableStyles() // disable styles, leave formatting styles
	}

	var result []string

	for _, trace := range traces {
		result = append(result, f.formatTrace(trace))
	}

	return strings.Join(result, "\n\n")
}

func (f *Formatter) formatTrace(trace *StackTrace) string {
	// var result []string
	var result *tree.Tree

	// Format goroutine header
	header := fmt.Sprintf("Goroutine %s: %s",
		trace.GoroutineID, trace.GoroutineState)
	// result = append(result, f.config.Theme.Goroutine.Render(header))
	result = tree.New().Root(f.config.Theme.Goroutine.Render(header))

	// Calculate widths and counts
	functionCounts := f.countFunctions(trace.Entries)

	// Format entries
	for _, entry := range trace.Entries {
		formattedEntry := f.formatEntry(entry, functionCounts)
		// result = append(result, formattedEntry)
		result.Child(formattedEntry)
	}

	// return strings.Join(result, "\n")
	return result.String()
}

func (f *Formatter) countFunctions(entries []StackEntry) map[string]int {
	counts := make(map[string]int)
	for _, entry := range entries {
		counts[entry.FunctionName]++
	}
	return counts
}

func (f *Formatter) formatEntry(entry StackEntry, functionCounts map[string]int) string {
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
