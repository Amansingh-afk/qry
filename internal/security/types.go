package security

import "fmt"

// Mode determines how violations are handled
type Mode string

const (
	ModeStrict Mode = "strict" // Block SQL with violations
	ModeWarn   Mode = "warn"   // Warn but still return SQL
)

// ViolationType categorizes the type of security violation
type ViolationType string

const (
	ViolationTable   ViolationType = "table"
	ViolationColumn  ViolationType = "column"
	ViolationPattern ViolationType = "pattern"
)

// Config holds security settings from .qry.yaml
type Config struct {
	Enabled bool
	Mode    Mode
	Exclude ExcludeConfig
}

// ExcludeConfig defines what to exclude
type ExcludeConfig struct {
	Tables   []string // Exact table names
	Columns  []string // Exact column names
	Patterns []string // Wildcard patterns (*_secret, api_*)
}

// Violation represents a single security violation
type Violation struct {
	Type    ViolationType // table, column, pattern
	Name    string        // The matched name in SQL
	Rule    string        // The rule that matched
	Context string        // Additional context (e.g., "in FROM clause")
}

// Result holds the validation result
type Result struct {
	Valid      bool
	Violations []Violation
	SQL        string // Original SQL (for reference)
}

// Error returns a formatted error message for the violations
func (r *Result) Error() string {
	if r.Valid || len(r.Violations) == 0 {
		return ""
	}

	msg := "Security violation: query references excluded data\n"
	for _, v := range r.Violations {
		msg += "  - " + string(v.Type) + ": " + v.Name
		if v.Rule != "" && v.Rule != v.Name {
			msg += " (matched rule: " + v.Rule + ")"
		}
		if v.Context != "" {
			msg += " " + v.Context
		}
		msg += "\n"
	}
	return msg
}

// Summary returns a short summary of violations
func (r *Result) Summary() string {
	if r.Valid {
		return ""
	}
	if len(r.Violations) == 1 {
		v := r.Violations[0]
		return "blocked: references " + string(v.Type) + " '" + v.Name + "'"
	}
	return fmt.Sprintf("blocked: %d security violations", len(r.Violations))
}
