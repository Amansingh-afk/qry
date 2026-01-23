package security

import "github.com/spf13/viper"

// LoadConfig loads security configuration from viper
// Returns nil if security is not configured
func LoadConfig() *Config {
	// Check if security section exists
	if !viper.IsSet("security") {
		return nil
	}

	mode := viper.GetString("security.mode")
	if mode == "" {
		mode = string(ModeWarn) // Default to warn mode
	}

	tables := viper.GetStringSlice("security.exclude.tables")
	columns := viper.GetStringSlice("security.exclude.columns")
	patterns := viper.GetStringSlice("security.exclude.patterns")

	// If nothing to exclude, security is effectively disabled
	if len(tables) == 0 && len(columns) == 0 && len(patterns) == 0 {
		return nil
	}

	return &Config{
		Enabled: true,
		Mode:    Mode(mode),
		Exclude: ExcludeConfig{
			Tables:   tables,
			Columns:  columns,
			Patterns: patterns,
		},
	}
}

// IsStrict returns true if mode is strict (block on violation)
func (c *Config) IsStrict() bool {
	return c != nil && c.Mode == ModeStrict
}

// IsWarn returns true if mode is warn (warn but allow)
func (c *Config) IsWarn() bool {
	return c != nil && c.Mode == ModeWarn
}

// HasExclusions returns true if there are any exclusion rules
func (c *Config) HasExclusions() bool {
	if c == nil {
		return false
	}
	return len(c.Exclude.Tables) > 0 ||
		len(c.Exclude.Columns) > 0 ||
		len(c.Exclude.Patterns) > 0
}
