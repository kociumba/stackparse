// utility to parse go stack traces and make them more readable
package stackparse

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type stackEntry struct {
	functionName string
	args         string
	file         string
	line         string
	isCreatedBy  bool
}

type Option func(*config)

// WithColorize returns an Option that enables colorizing of the stack trace output using ANSI escape codes.
func Color(opt bool) Option {
	return func(cfg *config) {
		cfg.colorize = opt
	}
}

// do not color the output, usefull if piping to another program or a file
func NoColor() Option {
	return Color(false)
}

// WithSimple returns an Option that enables simplified output.
func Simple(opt bool) Option {
	return func(cfg *config) {
		cfg.simple = opt
	}
}

// print full filepaths and extra info
func NoSimple() Option {
	return Simple(false)
}

type config struct {
	colorize bool // should the output be colorized with ansi escape codes
	simple   bool // should the output be simplified by omitting certain details
}

var (
	baseStyle      lipgloss.Style
	goroutineStyle lipgloss.Style
	functionStyle  lipgloss.Style
	argsStyle      lipgloss.Style
	fileStyle      lipgloss.Style
	lineStyle      lipgloss.Style
	createdByStyle lipgloss.Style
	repeatStyle    lipgloss.Style
)

// Style definitions
func initStyles() {
	// Base styles
	baseStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	// Goroutine style
	goroutineStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ADD8")). // Go blue
		MarginTop(1).
		MarginBottom(1)

	// Function style
	functionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#98C379")) // Soft green
		// PaddingLeft(4)

	// Args style
	argsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61AFEF")). // Light blue
		PaddingLeft(2)

	// File style
	fileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C678DD")) // Purple
		// PaddingLeft(4)

	// Line number style
	lineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5C07B")). // Gold
		PaddingLeft(2)

	// Created by style
	createdByStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E06C75")) // Soft red
		// PaddingLeft(4)

	// Repeat count style
	repeatStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E06C75")). // Soft red
		Italic(true)
}

// Regular expressions for different parts of the stack trace
var (
	// Matches goroutine header line
	goroutineRegex = regexp.MustCompile(`goroutine (\d+) \[([\w\.]+)\]:`)

	// Matches function calls with arguments
	functionRegex = regexp.MustCompile(`^(\S+)\((.*)\)$`)

	// Matches file location lines
	locationRegex = regexp.MustCompile(`^\s*(.+\.go):(\d+)(.*)$`)

	// Matches "created by" lines
	createdByRegex = regexp.MustCompile(`created by (.+) in goroutine (\d+)`)
)

// Parse is the main entry point for parsing stack traces
func Parse(stack []byte, options ...Option) []byte {
	cfg := config{
		colorize: true,
		simple:   true,
	}

	for _, opt := range options {
		opt(&cfg)
	}

	initStyles()

	lines := strings.Split(string(stack), "\n")
	return []byte(parseStackTrace(lines, cfg))
}

// func formatStackLine(label, content string, style lipgloss.Style) string {
// 	return style.Render(label) + content
// }

func parseStackTrace(lines []string, cfg config) string {
	var result []string
	var entries []stackEntry
	var currentEntry *stackEntry
	functionCounts := make(map[string]int)

	// reset all styles that have color
	if !cfg.colorize {
		goroutineStyle = lipgloss.NewStyle()
		functionStyle = lipgloss.NewStyle()
		argsStyle = lipgloss.NewStyle()
		fileStyle = lipgloss.NewStyle()
		lineStyle = lipgloss.NewStyle()
		createdByStyle = lipgloss.NewStyle()
		repeatStyle = lipgloss.NewStyle()
	}

	// First pass: collect all entries
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Handle goroutine header
		if match := goroutineRegex.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
				currentEntry = nil
			}
			header := fmt.Sprintf("Goroutine %s: %s", match[1], match[2])
			result = append(result, goroutineStyle.Render(header))
			continue
		}

		// Handle function calls
		if match := functionRegex.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
			}
			currentEntry = &stackEntry{
				functionName: match[1],
				args:         match[2],
			}
			functionCounts[currentEntry.functionName]++
			continue
		}

		// Handle file locations
		if match := locationRegex.FindStringSubmatch(line); match != nil && currentEntry != nil {
			filePath := match[1]
			if cfg.simple {
				filePath = simplifyPath(filePath)
			}
			currentEntry.file = filePath
			currentEntry.line = match[2]
			continue
		}

		// Handle "created by" lines
		if match := createdByRegex.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
			}
			currentEntry = &stackEntry{
				functionName: match[1],
				isCreatedBy:  true,
			}
			if i+1 < len(lines) && locationRegex.MatchString(lines[i+1]) {
				locMatch := locationRegex.FindStringSubmatch(lines[i+1])
				filePath := locMatch[1]
				if cfg.simple {
					filePath = simplifyPath(filePath)
				}
				currentEntry.file = filePath
				currentEntry.line = locMatch[2]
				i++
			}
			continue
		}
	}

	if currentEntry != nil {
		entries = append(entries, *currentEntry)
	}

	// Calculate maximum widths
	var maxFunctionWidth, maxFileWidth, maxLineWidth int
	for _, entry := range entries {
		maxFunctionWidth = max(maxFunctionWidth, lipgloss.Width(entry.functionName))
		maxFileWidth = max(maxFileWidth, lipgloss.Width(entry.file))
		maxLineWidth = max(maxLineWidth, lipgloss.Width(entry.line))
	}

	// Create style for aligned content
	// alignedStyle := lipgloss.NewStyle().Width(maxFunctionWidth + maxFileWidth + maxLineWidth + 20)

	// Second pass: format entries with tree structure
	var currentTree *tree.Tree
	for _, entry := range entries {
		if functionCounts[entry.functionName] <= 1 {
			// Create a tree for each function call
			if entry.isCreatedBy {
				currentTree = tree.New().Root(
					createdByStyle.Render("Created by: ") +
						entry.functionName,
				)
			} else {
				currentTree = tree.New().Root(
					functionStyle.Render("Function: ") +
						lipgloss.NewStyle().Width(maxFunctionWidth).Render(entry.functionName) +
						argsStyle.Render(fmt.Sprintf("(%s)", entry.args)),
				)
			}

			// Add file location as a child of the function
			fileInfo := lipgloss.JoinHorizontal(
				lipgloss.Left,
				fileStyle.Render("At: "),
				lipgloss.NewStyle().Width(maxFileWidth).Render(entry.file),
				lineStyle.Render(fmt.Sprintf(" Line: %s", entry.line)),
			)
			currentTree.Child(fileInfo)

			// Render the tree and add it to results
			result = append(result, currentTree.String())
		} else {
			// For repeated functions, add a counter
			repCount := repeatStyle.Render(fmt.Sprintf(" (repeated %d times)", functionCounts[entry.functionName]))
			currentTree = tree.New().Root(
				functionStyle.Render("Function: ") +
					lipgloss.NewStyle().Width(maxFunctionWidth).Render(entry.functionName) +
					argsStyle.Render(fmt.Sprintf("(%s)", entry.args)) +
					repCount,
			)
			result = append(result, currentTree.String())
		}
	}

	initStyles()

	return strings.Join(result, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func simplifyPath(path string) string {
	return filepath.Base(path)
}
