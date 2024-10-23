package stackparse

// StackEntry represents a single entry in the stack trace
type StackEntry struct {
	FunctionName       string
	FullName           string
	Args               string
	File               string
	Line               string
	Offset             string
	IsCreatedBy        bool
	CreatedByGoroutine string
}

// StackTrace represents a complete stack trace
type StackTrace struct {
	Entries        []StackEntry
	GoroutineID    string
	GoroutineState string
}
