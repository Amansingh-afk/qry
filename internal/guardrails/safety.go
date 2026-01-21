package guardrails

import (
	"regexp"
	"strings"
)

var dangerous = []string{
	"DROP TABLE",
	"DROP DATABASE",
	"TRUNCATE",
	"DELETE FROM",
	"ALTER TABLE",
	"DROP INDEX",
}

var whereClause = regexp.MustCompile(`(?i)\bWHERE\b`)

func Check(sql string) string {
	upper := strings.ToUpper(sql)

	for _, d := range dangerous {
		if strings.Contains(upper, d) {
			return "Destructive operation: " + d
		}
	}

	if strings.Contains(upper, "UPDATE") && !whereClause.MatchString(sql) {
		return "UPDATE without WHERE clause"
	}

	return ""
}
