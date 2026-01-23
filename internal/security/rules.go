package security

import (
	"regexp"
	"strings"
)

// Matcher checks if a name matches exclusion rules
type Matcher struct {
	exactTables  map[string]bool
	exactColumns map[string]bool
	patterns     []*compiledPattern
}

type compiledPattern struct {
	original string
	regex    *regexp.Regexp
}

// NewMatcher creates a new matcher from exclusion config
func NewMatcher(cfg *Config) *Matcher {
	if cfg == nil {
		return &Matcher{
			exactTables:  make(map[string]bool),
			exactColumns: make(map[string]bool),
		}
	}

	m := &Matcher{
		exactTables:  make(map[string]bool),
		exactColumns: make(map[string]bool),
	}

	// Index exact matches for O(1) lookup
	for _, t := range cfg.Exclude.Tables {
		m.exactTables[strings.ToLower(t)] = true
	}
	for _, c := range cfg.Exclude.Columns {
		m.exactColumns[strings.ToLower(c)] = true
	}

	// Compile patterns
	for _, p := range cfg.Exclude.Patterns {
		if compiled := compilePattern(p); compiled != nil {
			m.patterns = append(m.patterns, compiled)
		}
	}

	return m
}

// compilePattern converts a wildcard pattern to regex
// Supports: * (any characters), ? (single character)
func compilePattern(pattern string) *compiledPattern {
	if pattern == "" {
		return nil
	}

	// Escape regex special chars except * and ?
	escaped := regexp.QuoteMeta(pattern)

	// Convert wildcards to regex
	escaped = strings.ReplaceAll(escaped, `\*`, `.*`)
	escaped = strings.ReplaceAll(escaped, `\?`, `.`)

	// Anchor the pattern
	regexStr := "^" + escaped + "$"

	re, err := regexp.Compile("(?i)" + regexStr) // Case insensitive
	if err != nil {
		return nil
	}

	return &compiledPattern{
		original: pattern,
		regex:    re,
	}
}

// MatchTable checks if a table name is excluded
func (m *Matcher) MatchTable(table string) (matched bool, rule string) {
	lower := strings.ToLower(table)

	// Check exact match
	if m.exactTables[lower] {
		return true, table
	}

	// Check patterns
	for _, p := range m.patterns {
		if p.regex.MatchString(table) {
			return true, p.original
		}
	}

	return false, ""
}

// MatchColumn checks if a column name is excluded
func (m *Matcher) MatchColumn(column string) (matched bool, rule string) {
	lower := strings.ToLower(column)

	// Check exact match
	if m.exactColumns[lower] {
		return true, column
	}

	// Check patterns
	for _, p := range m.patterns {
		if p.regex.MatchString(column) {
			return true, p.original
		}
	}

	return false, ""
}

// MatchAny checks if a name matches any exclusion rule (table or column)
func (m *Matcher) MatchAny(name string) (matched bool, rule string, vType ViolationType) {
	if matched, rule := m.MatchTable(name); matched {
		return true, rule, ViolationTable
	}
	if matched, rule := m.MatchColumn(name); matched {
		return true, rule, ViolationColumn
	}
	return false, "", ""
}
