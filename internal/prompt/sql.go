package prompt

import (
	"regexp"
	"strings"
)

func BuildSQL(query string, dialect string) string {
	dialectHint := ""
	if dialect != "" {
		dialectHint = "\n- Use " + dialect + " syntax"
	}

	return `You are a SQL expert. Based on the codebase context (schemas, migrations, models), generate ONLY the SQL query.

Rules:
- Output ONLY the SQL, no explanation
- Use actual table/column names from the codebase` + dialectHint + `

Request: ` + query
}

var sqlBlock = regexp.MustCompile("(?s)```sql\\s*(.+?)\\s*```")
var anyBlock = regexp.MustCompile("(?s)```\\s*(.+?)\\s*```")

func ExtractSQL(response string) string {
	if m := sqlBlock.FindStringSubmatch(response); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	if m := anyBlock.FindStringSubmatch(response); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	return strings.TrimSpace(response)
}
