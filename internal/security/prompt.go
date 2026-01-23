package security

import (
	"fmt"
	"strings"
)

// BuildPromptAddition generates the security rules portion to add to the prompt
// Returns empty string if no security config
func BuildPromptAddition(cfg *Config) string {
	if cfg == nil || !cfg.HasExclusions() {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("\n\nSECURITY RULES (MUST FOLLOW):\n")
	sb.WriteString("You must NEVER access, query, or return data from the following:\n")

	// Tables
	if len(cfg.Exclude.Tables) > 0 {
		sb.WriteString("\nForbidden tables:\n")
		for _, t := range cfg.Exclude.Tables {
			sb.WriteString("  - ")
			sb.WriteString(t)
			sb.WriteString("\n")
		}
	}

	// Columns
	if len(cfg.Exclude.Columns) > 0 {
		sb.WriteString("\nForbidden columns:\n")
		for _, c := range cfg.Exclude.Columns {
			sb.WriteString("  - ")
			sb.WriteString(c)
			sb.WriteString("\n")
		}
	}

	// Patterns
	if len(cfg.Exclude.Patterns) > 0 {
		sb.WriteString("\nForbidden patterns (any table/column matching):\n")
		for _, p := range cfg.Exclude.Patterns {
			sb.WriteString("  - ")
			sb.WriteString(p)
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\nIf a query requires accessing forbidden data, respond with: ")
	sb.WriteString("\"Cannot generate this query: it would access restricted data.\"\n")

	return sb.String()
}

// BuildPromptSummary returns a brief summary for logging/debugging
func BuildPromptSummary(cfg *Config) string {
	if cfg == nil || !cfg.HasExclusions() {
		return "security: disabled"
	}

	parts := []string{}

	if len(cfg.Exclude.Tables) > 0 {
		parts = append(parts, fmt.Sprintf("%d tables", len(cfg.Exclude.Tables)))
	}
	if len(cfg.Exclude.Columns) > 0 {
		parts = append(parts, fmt.Sprintf("%d columns", len(cfg.Exclude.Columns)))
	}
	if len(cfg.Exclude.Patterns) > 0 {
		parts = append(parts, fmt.Sprintf("%d patterns", len(cfg.Exclude.Patterns)))
	}

	return "security: " + string(cfg.Mode) + " mode, excluding " + strings.Join(parts, ", ")
}
