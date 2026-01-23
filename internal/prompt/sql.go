package prompt

import (
	"regexp"
	"strings"

	"github.com/amansingh-afk/qry/internal/security"
	"github.com/spf13/viper"
)

const defaultPromptTemplate = `You are a SQL expert. Based on the codebase context (schemas, migrations, models), generate ONLY the SQL query.

Rules:
- Output ONLY the SQL, no explanation
- Use actual table/column names from the codebase
- Use {{dialect}}{{version}} syntax

Request: {{query}}`

// BuildSQL builds the full prompt for the first query in a session.
// This includes the role, rules, security rules, and the query.
func BuildSQL(query string, dialect string) string {
	// Get prompt template from config or use default
	template := viper.GetString("prompt")
	if template == "" {
		template = defaultPromptTemplate
	}

	// Get db version from config
	dbVersion := viper.GetString("db_version")
	versionStr := ""
	if dbVersion != "" {
		versionStr = " " + dbVersion
	}

	// Replace template variables
	result := template
	result = strings.ReplaceAll(result, "{{query}}", query)
	result = strings.ReplaceAll(result, "{{dialect}}", dialect)
	result = strings.ReplaceAll(result, "{{version}}", versionStr)

	// Add security rules if configured
	securityRules := security.PromptAddition()
	if securityRules != "" {
		result += securityRules
	}

	return result
}

// BuildFollowUp builds a minimal prompt for subsequent queries in an existing session.
// The LLM already knows its role from the first query.
func BuildFollowUp(query string) string {
	return query
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
