package security

import "sync"

// Security provides the main interface for the security layer
type Security struct {
	config    *Config
	validator *Validator
}

var (
	instance *Security
	once     sync.Once
)

// Get returns the singleton security instance
// Loads config on first call
func Get() *Security {
	once.Do(func() {
		cfg := LoadConfig()
		instance = &Security{
			config:    cfg,
			validator: NewValidator(cfg),
		}
	})
	return instance
}

// Reset clears the singleton (useful for testing or config reload)
func Reset() {
	once = sync.Once{}
	instance = nil
}

// IsEnabled returns true if security is configured
func (s *Security) IsEnabled() bool {
	return s.config != nil && s.config.Enabled
}

// GetConfig returns the security config (may be nil)
func (s *Security) GetConfig() *Config {
	return s.config
}

// GetMode returns the security mode (strict/warn)
func (s *Security) GetMode() Mode {
	if s.config == nil {
		return ModeWarn
	}
	return s.config.Mode
}

// Validate checks SQL against security rules
func (s *Security) Validate(sql string) *Result {
	return s.validator.Validate(sql)
}

// IsBlocked returns true if the result should be blocked
func (s *Security) IsBlocked(result *Result) bool {
	return s.validator.IsBlocked(result)
}

// ShouldWarn returns true if warnings should be shown
func (s *Security) ShouldWarn(result *Result) bool {
	return s.validator.ShouldWarn(result)
}

// GetPromptAddition returns the security rules to add to the prompt
func (s *Security) GetPromptAddition() string {
	return BuildPromptAddition(s.config)
}

// GetPromptSummary returns a brief summary for logging
func (s *Security) GetPromptSummary() string {
	return BuildPromptSummary(s.config)
}

// --- Convenience functions for direct use ---

// Enabled returns true if security is configured
func Enabled() bool {
	return Get().IsEnabled()
}

// Validate checks SQL against security rules
func Validate(sql string) *Result {
	return Get().Validate(sql)
}

// PromptAddition returns the security rules to add to the prompt
func PromptAddition() string {
	return Get().GetPromptAddition()
}
