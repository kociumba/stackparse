package stackparse

// Parse is the main entry point for parsing stack traces
func Parse(stack []byte, options ...Option) []byte {
	parser := NewParser(options...)
	trace := parser.Parse(stack)

	formatter := NewFormatter(parser.config)
	result := formatter.Format(trace)

	return []byte(result)
}
