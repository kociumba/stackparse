package stackparse

// StackEntry represents a single entry in the stack trace
type StackEntry struct {
	FunctionName string
	Args         string
	File         string
	Line         string
	IsCreatedBy  bool
}

// StackTrace represents a complete stack trace
type StackTrace struct {
	Entries        []StackEntry
	GoroutineID    string
	GoroutineState string
}
