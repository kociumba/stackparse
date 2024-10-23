package stackparse

import (
	"path/filepath"
	"regexp"
	"strings"
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
}

// NewParser creates a new Parser instance with the given options
func NewParser(options ...Option) *Parser {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}

	return &Parser{
		config: config,
		patterns: &Patterns{
			Goroutine: regexp.MustCompile(`goroutine (\d+) \[([\w\.]+)\]:`),
			Function:  regexp.MustCompile(`^(\S+)\((.*)\)$`),
			Location:  regexp.MustCompile(`^\s*(.+\.go):(\d+)(.*)$`),
			CreatedBy: regexp.MustCompile(`created by (.+) in goroutine (\d+)`),
		},
	}
}

// Parse converts a byte slice containing a stack trace into a StackTrace
func (p *Parser) Parse(stack []byte) *StackTrace {
	lines := strings.Split(string(stack), "\n")
	return p.parseLines(lines)
}

func (p *Parser) parseLines(lines []string) *StackTrace {
	trace := &StackTrace{}
	var currentEntry *StackEntry

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Handle goroutine header
		if match := p.patterns.Goroutine.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				trace.Entries = append(trace.Entries, *currentEntry)
				currentEntry = nil
			}
			trace.GoroutineID = match[1]
			trace.GoroutineState = match[2]
			continue
		}

		// Handle function calls
		if match := p.patterns.Function.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				trace.Entries = append(trace.Entries, *currentEntry)
			}
			currentEntry = &StackEntry{
				FunctionName: match[1],
				Args:         match[2],
			}
			continue
		}

		// Handle file locations
		if match := p.patterns.Location.FindStringSubmatch(line); match != nil && currentEntry != nil {
			filePath := match[1]
			if p.config.Simple {
				filePath = p.simplifyPath(filePath)
			}
			currentEntry.File = filePath
			currentEntry.Line = match[2]
			continue
		}

		// Handle "created by" lines
		if match := p.patterns.CreatedBy.FindStringSubmatch(line); match != nil {
			if currentEntry != nil {
				trace.Entries = append(trace.Entries, *currentEntry)
			}
			currentEntry = &StackEntry{
				FunctionName: match[1],
				IsCreatedBy:  true,
			}
			if i+1 < len(lines) && p.patterns.Location.MatchString(lines[i+1]) {
				locMatch := p.patterns.Location.FindStringSubmatch(lines[i+1])
				filePath := locMatch[1]
				if p.config.Simple {
					filePath = p.simplifyPath(filePath)
				}
				currentEntry.File = filePath
				currentEntry.Line = locMatch[2]
				i++
			}
			continue
		}
	}

	if currentEntry != nil {
		trace.Entries = append(trace.Entries, *currentEntry)
	}

	return trace
}

func (p *Parser) simplifyPath(path string) string {
	return filepath.Base(path)
}
