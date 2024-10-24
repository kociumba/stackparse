package stackparse

// Parse is the main entry point for parsing stack traces
//
// You can use this to parse stack traces from your code like this:
//
//	buf := make([]byte, 1<<16)
//	runtime.Stack(buf, true)
//	parsed := stackparse.Parse(buf)
//	// use the parsed stack however you want
//	// the return is also []byte, so you can do things like:
//	os.Stderr.Write(parsed)
func Parse(stack []byte, options ...Option) []byte {
	parser := NewParser(options...)
	trace := parser.Parse(stack)

	formatter := NewFormatter(parser.config)
	result := formatter.Format(trace)

	return []byte(result)
}

// Use this to parse the stack trace in place overwriting the original buffer
//
//	buf := make([]byte, 1<<16)
//	runtime.Stack(buf, true)
//	stackparse.Parse(&buf)
//	os.Stderr.Write(buf)
func ParseStatic(stack *[]byte, options ...Option) {
	*stack = Parse(*stack, options...)
}
