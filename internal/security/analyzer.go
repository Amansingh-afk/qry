package security

import (
	"regexp"
	"strings"
)

// SQLRef represents a reference found in SQL
type SQLRef struct {
	Name    string
	Type    string // "table" or "column"
	Context string // e.g., "FROM clause", "JOIN", "SELECT"
}

// AnalyzeSQL extracts table and column references from SQL
func AnalyzeSQL(sql string) []SQLRef {
	var refs []SQLRef

	// Normalize SQL
	normalized := strings.ToUpper(sql)

	// Extract table references
	refs = append(refs, extractTables(sql, normalized)...)

	// Extract column references
	refs = append(refs, extractColumns(sql, normalized)...)

	return refs
}

// extractTables finds table names after FROM, JOIN, INTO, UPDATE
func extractTables(sql, normalized string) []SQLRef {
	var refs []SQLRef

	// Patterns for table extraction
	patterns := []struct {
		regex   *regexp.Regexp
		context string
	}{
		{regexp.MustCompile(`(?i)\bFROM\s+([a-zA-Z_][a-zA-Z0-9_]*)`), "FROM clause"},
		{regexp.MustCompile(`(?i)\bJOIN\s+([a-zA-Z_][a-zA-Z0-9_]*)`), "JOIN clause"},
		{regexp.MustCompile(`(?i)\bINTO\s+([a-zA-Z_][a-zA-Z0-9_]*)`), "INTO clause"},
		{regexp.MustCompile(`(?i)\bUPDATE\s+([a-zA-Z_][a-zA-Z0-9_]*)`), "UPDATE clause"},
		{regexp.MustCompile(`(?i)\bTRUNCATE\s+(?:TABLE\s+)?([a-zA-Z_][a-zA-Z0-9_]*)`), "TRUNCATE"},
		{regexp.MustCompile(`(?i)\bDROP\s+TABLE\s+(?:IF\s+EXISTS\s+)?([a-zA-Z_][a-zA-Z0-9_]*)`), "DROP TABLE"},
		{regexp.MustCompile(`(?i)\bDELETE\s+FROM\s+([a-zA-Z_][a-zA-Z0-9_]*)`), "DELETE FROM"},
	}

	seen := make(map[string]bool)

	for _, p := range patterns {
		matches := p.regex.FindAllStringSubmatch(sql, -1)
		for _, m := range matches {
			if len(m) > 1 {
				table := m[1]
				key := strings.ToLower(table)
				if !seen[key] && !isSQLKeyword(table) {
					seen[key] = true
					refs = append(refs, SQLRef{
						Name:    table,
						Type:    "table",
						Context: p.context,
					})
				}
			}
		}
	}

	return refs
}

// extractColumns finds column names from SELECT, WHERE, etc.
func extractColumns(sql, normalized string) []SQLRef {
	var refs []SQLRef
	seen := make(map[string]bool)

	// Extract from SELECT clause (before FROM)
	selectPattern := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM`)
	if m := selectPattern.FindStringSubmatch(sql); len(m) > 1 {
		cols := parseColumnList(m[1])
		for _, col := range cols {
			key := strings.ToLower(col)
			if !seen[key] {
				seen[key] = true
				refs = append(refs, SQLRef{
					Name:    col,
					Type:    "column",
					Context: "SELECT clause",
				})
			}
		}
	}

	// Extract from WHERE clause
	wherePattern := regexp.MustCompile(`(?i)\bWHERE\s+(.+?)(?:\s+ORDER\s+BY|\s+GROUP\s+BY|\s+LIMIT|\s+HAVING|;|$)`)
	if m := wherePattern.FindStringSubmatch(sql); len(m) > 1 {
		cols := extractColumnRefs(m[1])
		for _, col := range cols {
			key := strings.ToLower(col)
			if !seen[key] {
				seen[key] = true
				refs = append(refs, SQLRef{
					Name:    col,
					Type:    "column",
					Context: "WHERE clause",
				})
			}
		}
	}

	// Extract from ORDER BY
	orderPattern := regexp.MustCompile(`(?i)\bORDER\s+BY\s+(.+?)(?:\s+LIMIT|;|$)`)
	if m := orderPattern.FindStringSubmatch(sql); len(m) > 1 {
		cols := parseColumnList(m[1])
		for _, col := range cols {
			key := strings.ToLower(col)
			if !seen[key] {
				seen[key] = true
				refs = append(refs, SQLRef{
					Name:    col,
					Type:    "column",
					Context: "ORDER BY clause",
				})
			}
		}
	}

	// Extract from GROUP BY
	groupPattern := regexp.MustCompile(`(?i)\bGROUP\s+BY\s+(.+?)(?:\s+ORDER\s+BY|\s+HAVING|\s+LIMIT|;|$)`)
	if m := groupPattern.FindStringSubmatch(sql); len(m) > 1 {
		cols := parseColumnList(m[1])
		for _, col := range cols {
			key := strings.ToLower(col)
			if !seen[key] {
				seen[key] = true
				refs = append(refs, SQLRef{
					Name:    col,
					Type:    "column",
					Context: "GROUP BY clause",
				})
			}
		}
	}

	return refs
}

// parseColumnList parses a comma-separated column list
func parseColumnList(s string) []string {
	var cols []string

	// Handle common patterns
	// Remove function calls, keep only column names
	colPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\b`)
	matches := colPattern.FindAllStringSubmatch(s, -1)

	for _, m := range matches {
		if len(m) > 0 {
			col := m[1]
			if !isSQLKeyword(col) && !isSQLFunction(col) {
				cols = append(cols, col)
			}
		}
	}

	return cols
}

// extractColumnRefs extracts column references from expressions
func extractColumnRefs(s string) []string {
	return parseColumnList(s)
}

// Common SQL keywords to ignore
var sqlKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true, "AND": true, "OR": true,
	"JOIN": true, "LEFT": true, "RIGHT": true, "INNER": true, "OUTER": true,
	"ON": true, "AS": true, "IN": true, "NOT": true, "NULL": true,
	"IS": true, "LIKE": true, "BETWEEN": true, "EXISTS": true,
	"GROUP": true, "BY": true, "ORDER": true, "HAVING": true,
	"LIMIT": true, "OFFSET": true, "ASC": true, "DESC": true,
	"INSERT": true, "INTO": true, "VALUES": true, "UPDATE": true,
	"SET": true, "DELETE": true, "CREATE": true, "ALTER": true,
	"DROP": true, "TABLE": true, "INDEX": true, "VIEW": true,
	"DISTINCT": true, "ALL": true, "CASE": true, "WHEN": true,
	"THEN": true, "ELSE": true, "END": true, "UNION": true,
	"WITH": true, "OVER": true, "PARTITION": true, "TRUE": true,
	"FALSE": true, "INTERVAL": true, "CURRENT_DATE": true,
	"CURRENT_TIMESTAMP": true, "IF": true, "TRUNCATE": true,
}

func isSQLKeyword(s string) bool {
	return sqlKeywords[strings.ToUpper(s)]
}

// Common SQL functions to ignore
var sqlFunctions = map[string]bool{
	"COUNT": true, "SUM": true, "AVG": true, "MIN": true, "MAX": true,
	"COALESCE": true, "NULLIF": true, "CAST": true, "CONCAT": true,
	"LOWER": true, "UPPER": true, "TRIM": true, "LENGTH": true,
	"SUBSTRING": true, "REPLACE": true, "DATE_TRUNC": true, "NOW": true,
	"EXTRACT": true, "TO_CHAR": true, "TO_DATE": true, "ROUND": true,
	"FLOOR": true, "CEIL": true, "ABS": true, "ROW_NUMBER": true,
	"RANK": true, "DENSE_RANK": true, "LAG": true, "LEAD": true,
	"DATE": true, "YEAR": true, "MONTH": true, "DAY": true,
}

func isSQLFunction(s string) bool {
	return sqlFunctions[strings.ToUpper(s)]
}
