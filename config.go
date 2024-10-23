package stackparse

// Config holds all configuration options for the parser
type Config struct {
	Colorize bool
	Simple   bool
	Theme    *Theme
}

// Option is a function type for setting configuration options
type Option func(*Config)

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Colorize: true,
		Simple:   true,
		Theme:    DefaultTheme(),
	}
}

// WithColor enables or disables colorized output
func WithColor(enabled bool) Option {
	return func(c *Config) {
		c.Colorize = enabled
	}
}

// WithSimple enables or disables simplified output
func WithSimple(enabled bool) Option {
	return func(c *Config) {
		c.Simple = enabled
	}
}

// WithTheme sets a custom theme
func WithTheme(theme *Theme) Option {
	return func(c *Config) {
		c.Theme = theme
	}
}
