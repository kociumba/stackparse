package stackparse

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/term"
)

// Parser handles the parsing of stack traces
type Parser struct {
	config   *Config
	patterns *Patterns
}

// Patterns holds all regex patterns used for parsing
type Patterns struct {
	Goroutine *regexp.Regexp
	Function  *regexp.Regexp
	Location  *regexp.Regexp
	CreatedBy *regexp.Regexp
	LongFunc  *regexp.Regexp
}

// NewParser creates a new Parser instance with the given options
//
// Mostly used internally but exposed for some edge cases
func NewParser(options ...Option) *Parser {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}

	return &Parser{
		config: config,
		patterns: &Patterns{
			// universal all function regex ^(\S+)\((.*)\)$
			Goroutine: regexp.MustCompile(`goroutine (\d+) \[(.*?)\]:`),
			Function:  regexp.MustCompile(`^([^\s/]+)\((.*)\)$`),
			Location:  regexp.MustCompile(`^\s*(.+\.go):(\d+)(?:\s+\+([0-9a-fA-Fx]+))?$`),
			CreatedBy: regexp.MustCompile(`created by (.+) in goroutine (\d+)`),
			LongFunc:  regexp.MustCompile(`^(\S*/\S+)\((.*)\)$`),
		},
	}
}

// Parse converts a byte slice containing a stack trace into a StackTrace
//
// I you are usting this by itself and not by simply calling stackparse.Parse(), always call this before the formatter.
func (p *Parser) Parse(stack []byte) []*StackTrace {
	lines := strings.Split(string(stack), "\n")
	return p.parseLines(lines)
}

// Assigns the types of lines to the StackTrace
func (p *Parser) parseLines(lines []string) []*StackTrace {
	var traces []*StackTrace
	var currentTrace *StackTrace
	var currentEntry *StackEntry

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Handle goroutine header - starts a new trace
		if match := p.patterns.Goroutine.FindStringSubmatch(line); match != nil {
			// Save current entry if exists
			if currentEntry != nil && currentTrace != nil {
				currentTrace.Entries = append(currentTrace.Entries, *currentEntry)
				currentEntry = nil
			}

			// Create new trace
			currentTrace = &StackTrace{
				GoroutineID:    match[1],
				GoroutineState: match[2],
			}
			traces = append(traces, currentTrace)
			continue
		}

		// Skip if no current trace
		if currentTrace == nil {
			continue
		}

		// Handle "created by" lines
		if match := p.patterns.CreatedBy.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				currentTrace.Entries = append(currentTrace.Entries, *currentEntry)
			}

			currentEntry = &StackEntry{
				FunctionName:       match[1],
				IsCreatedBy:        true,
				CreatedByGoroutine: match[2],
			}

			// Look ahead for location
			if i+1 < len(lines) && p.patterns.Location.MatchString(lines[i+1]) {
				locMatch := p.patterns.Location.FindStringSubmatch(lines[i+1])
				filePath := locMatch[1]
				if p.config.Simple {
					filePath = p.simplifyPath(filePath)
				}
				currentEntry.File = filePath
				currentEntry.Line = locMatch[2]
				if len(locMatch) > 3 && locMatch[3] != "" {
					currentEntry.Offset = locMatch[3]
				}
				i++
			}
			continue
		}

		// Try to match function calls - first try simple function pattern
		match := p.patterns.Function.FindStringSubmatch(line)
		if match == nil {
			// If simple pattern doesn't match, try long function pattern
			match = p.patterns.LongFunc.FindStringSubmatch(line)
		}

		if match != nil {
			if currentEntry != nil {
				currentTrace.Entries = append(currentTrace.Entries, *currentEntry)
			}

			funcName := match[1]
			args := match[2]

			fd := int(os.Stdout.Fd())
			termWidth, _, err := term.GetSize(fd)
			if err != nil {
				termWidth = 60 * 8
			}

			// Handle long function names
			if len(funcName) > termWidth/8 && p.config.Simple { // Only simplify if config.Simple is true
				parts := strings.Split(funcName, "/")
				if len(parts) > 3 {
					// Keep the last three parts
					funcName = ".../" + strings.Join(parts[len(parts)-3:], "/")
				}
			}

			currentEntry = &StackEntry{
				FunctionName: funcName,
				Args:         args,
				FullName:     match[1], // Store full name for reference
			}

			// Look ahead for location
			if i+1 < len(lines) && p.patterns.Location.MatchString(lines[i+1]) {
				locMatch := p.patterns.Location.FindStringSubmatch(lines[i+1])
				filePath := locMatch[1]
				if p.config.Simple {
					filePath = p.simplifyPath(filePath)
				}
				currentEntry.File = filePath
				currentEntry.Line = locMatch[2]
				if len(locMatch) > 3 && locMatch[3] != "" {
					currentEntry.Offset = locMatch[3]
				}
				i++
			}
			continue
		}

		// Handle file locations if we somehow missed them earlier
		if match := p.patterns.Location.FindStringSubmatch(line); match != nil && currentEntry != nil {
			filePath := match[1]
			if p.config.Simple {
				filePath = p.simplifyPath(filePath)
			}
			currentEntry.File = filePath
			currentEntry.Line = match[2]
			if len(match) > 3 && match[3] != "" {
				currentEntry.Offset = match[3]
			}
		}
	}

	// Add final entry if exists
	if currentEntry != nil && currentTrace != nil {
		currentTrace.Entries = append(currentTrace.Entries, *currentEntry)
	}

	return traces
}

// Simple utility wrapper
func (p *Parser) simplifyPath(path string) string {
	if p.config.Simple {
		return filepath.Base(path)
	}
	return path
}
