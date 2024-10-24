# stackparse

A Go library for parsing and formatting stack traces with powerful customization options. This library makes debugging crashes and deliberate stack dumps more manageable by providing beautifully formatted stack traces with customizable styling.

## Installation

```bash
go get github.com/kociumba/stackparse
```

## Basic Usage

The simplest way to use stackparse is to capture and parse a stack trace:

```go
package main

import (
    "os"
    "runtime"
    "github.com/kociumba/stackparse"
)

func main() {
    // Create a buffer for the stack trace
    buf := make([]byte, 1<<16)
    
    // Capture the stack trace
    runtime.Stack(buf, true)
    
    // Parse and format the stack trace
    parsed := stackparse.Parse(buf)
    // modify the buffer in place with:
    stackparse.ParseStatic(&buf)
    
    // Write to stderr (or any other output)
    os.Stderr.Write(parsed)
    // if using stackparse.ParseStatic(), simply write from the buffer
    os.Stderr.Write(buf)
}
```

> [!IMPORTANT]
> Right now stackparse does not parse the reason given when when calling `panic()`, but this is planned for the next release.

## Configuration Options

stackparse provides several configuration options to customize the output:

> [!NOTE]
> Default settings are: Colorize: true, Simple: true, Theme: stackparse.DefaultTheme()

### Colorization

Control whether the output includes ANSI color codes:

```go
// Disable coloring the output
parsed := stackparse.Parse(buf, stackparse.WithColor(false))
```

### Simple Mode

Toggle between simple and detailed output formats (shortens both the function names and file paths):

```go
// Disable simple mode for more detailed output, does not guarantee that the formatting will be correct in all cases
parsed := stackparse.Parse(buf, stackparse.WithSimple(false))
```

## Theming

stackparse uses [lipgloss](https://github.com/charmbracelet/lipgloss) for styling and provides powerful theming capabilities.

### Using the Default Theme

The default theme provides a color scheme based on the [catppuccin theme](https://catppuccin.com/).

```go
// Stackparse uses the default theme without having to pass it in
parsed := stackparse.Parse(buf)

// You can get the default theme like this
defaultTheme := stackparse.DefaultTheme()
```

### Custom Themes

You can create custom themes by modifying the default theme or creating a new one from scratch:

```go
// Create a custom theme based on the default
myTheme := stackparse.DefaultTheme()

// overwrite the default styles compleatly
myTheme.Goroutine = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#ff0000")) // Red goroutine labels

myTheme.Function = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#00ff00")). // Green function names
    Italic(true)

// or build on top of the default theme
myTheme.Repeat = myTheme.Repeat.Faint(true).Blink(true) 

// Apply the custom theme
parsed := stackparse.Parse(buf, stackparse.WithTheme(myTheme))
```

### Theme Components

The Theme struct provides the following customizable components:

- `Goroutine`: Styling for goroutine headers
- `Function`: Styling for function names
- `Args`: Styling for function arguments
- `File`: Styling for file paths
- `Line`: Styling for line numbers
- `CreatedBy`: Styling for "created by" sections
- `Repeat`: Styling for repeated stack frames

### Custom Style Disabling

You can control how styles are disabled when colors are turned off with:

```go
myTheme := stackparse.DefaultTheme()

// Custom disable function
myTheme.SetDisableStylesFunc(func(t *stackparse.Theme) {
    // Keep bold formatting but remove colors
    t.Goroutine = t.Goroutine.UnsetForeground()
    t.Function = t.Function.UnsetForeground()
    t.Repeat = t.Repeat.UnsetBlink().UnsetFeint()
    
    // Keep some styling intact
    // everything else remains unchanged
})

parsed := stackparse.Parse(buf, 
    stackparse.WithTheme(myTheme),
    stackparse.WithColor(false), // This will trigger the custom disable function
)
```

## Output Example

When using the default theme, your stack trace will be formatted like this:

![default styling](/assets/image.png)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request with anything: fixes, demos, improvements, etc.
