package guardrails

import (
	"regexp"
	"strings"
)

var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bDROP\s+(TABLE|DATABASE|SCHEMA|INDEX)\b`),
	regexp.MustCompile(`(?i)\bTRUNCATE\s+TABLE\b`),
	regexp.MustCompile(`(?i)\bDELETE\s+FROM\s+\w+\s*;?\s*$`),
	regexp.MustCompile(`(?i)\bUPDATE\s+\w+\s+SET\s+.+\s*;?\s*$`),
	regexp.MustCompile(`(?i)\bALTER\s+TABLE\s+.+\s+DROP\b`),
}

func Check(sql string) string {
	normalized := strings.TrimSpace(sql)

	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(normalized) {
			return "⚠️  This query may be destructive. Review carefully before executing."
		}
	}

	upper := strings.ToUpper(normalized)
	if strings.Contains(upper, "DROP") ||
		strings.Contains(upper, "TRUNCATE") ||
		(strings.Contains(upper, "DELETE") && !strings.Contains(upper, "WHERE")) ||
		(strings.Contains(upper, "UPDATE") && !strings.Contains(upper, "WHERE")) {
		return "⚠️  This query may be destructive. Review carefully before executing."
	}

	return ""
}

func IsSafe(sql string) bool {
	return Check(sql) == ""
}
