package prompt

import (
	"regexp"
	"strings"
)

const sqlSystemPrompt = `You are a SQL expert. Generate SQL based on the user's natural language query.

Rules:
- Output ONLY the SQL query, no explanations
- Use standard SQL syntax
- Wrap the SQL in a code block with sql language tag
- If the query is ambiguous, make reasonable assumptions
- Prefer readable, well-formatted SQL`

func BuildSQL(query string) string {
	return sqlSystemPrompt + "\n\nUser query: " + query
}

var sqlBlockRegex = regexp.MustCompile("(?s)```sql\\s*(.+?)```")
var genericBlockRegex = regexp.MustCompile("(?s)```\\s*(.+?)```")

func ExtractSQL(response string) string {
	if matches := sqlBlockRegex.FindStringSubmatch(response); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	if matches := genericBlockRegex.FindStringSubmatch(response); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	lines := strings.Split(response, "\n")
	var sqlLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)
		if strings.HasPrefix(upper, "SELECT") ||
			strings.HasPrefix(upper, "INSERT") ||
			strings.HasPrefix(upper, "UPDATE") ||
			strings.HasPrefix(upper, "DELETE") ||
			strings.HasPrefix(upper, "CREATE") ||
			strings.HasPrefix(upper, "ALTER") ||
			strings.HasPrefix(upper, "DROP") ||
			strings.HasPrefix(upper, "WITH") ||
			strings.HasPrefix(upper, "FROM") ||
			strings.HasPrefix(upper, "WHERE") ||
			strings.HasPrefix(upper, "JOIN") ||
			strings.HasPrefix(upper, "ORDER") ||
			strings.HasPrefix(upper, "GROUP") ||
			strings.HasPrefix(upper, "HAVING") ||
			strings.HasPrefix(upper, "LIMIT") ||
			strings.HasPrefix(upper, "(") ||
			strings.HasPrefix(upper, ")") ||
			strings.HasPrefix(upper, "AND") ||
			strings.HasPrefix(upper, "OR") ||
			len(sqlLines) > 0 {
			sqlLines = append(sqlLines, trimmed)
		}
	}

	if len(sqlLines) > 0 {
		return strings.Join(sqlLines, "\n")
	}

	return strings.TrimSpace(response)
}
