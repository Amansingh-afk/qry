package security

// Validator checks SQL against security rules
type Validator struct {
	config  *Config
	matcher *Matcher
}

// NewValidator creates a new validator from config
func NewValidator(cfg *Config) *Validator {
	return &Validator{
		config:  cfg,
		matcher: NewMatcher(cfg),
	}
}

// Validate checks SQL for security violations
func (v *Validator) Validate(sql string) *Result {
	result := &Result{
		Valid: true,
		SQL:   sql,
	}

	// If no config or security disabled, always valid
	if v.config == nil || !v.config.Enabled {
		return result
	}

	// Analyze SQL to extract references
	refs := AnalyzeSQL(sql)

	// Check each reference against rules
	for _, ref := range refs {
		var matched bool
		var rule string

		switch ref.Type {
		case "table":
			matched, rule = v.matcher.MatchTable(ref.Name)
		case "column":
			matched, rule = v.matcher.MatchColumn(ref.Name)
		}

		if matched {
			result.Valid = false
			vType := ViolationTable
			if ref.Type == "column" {
				vType = ViolationColumn
			}

			result.Violations = append(result.Violations, Violation{
				Type:    vType,
				Name:    ref.Name,
				Rule:    rule,
				Context: ref.Context,
			})
		}
	}

	return result
}

// IsBlocked returns true if the result should be blocked (strict mode + violations)
func (v *Validator) IsBlocked(result *Result) bool {
	if result.Valid {
		return false
	}
	return v.config != nil && v.config.IsStrict()
}

// ShouldWarn returns true if warnings should be shown (warn mode + violations)
func (v *Validator) ShouldWarn(result *Result) bool {
	if result.Valid {
		return false
	}
	return v.config != nil && v.config.IsWarn()
}
