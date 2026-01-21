package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	success = lipgloss.Color("#5AF78E")
	muted   = lipgloss.Color("#636363")

	sqlKeyword = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6AC1")).
			Bold(true)

	metaStyle = lipgloss.NewStyle().
			Foreground(muted)

	successDot = lipgloss.NewStyle().
			Foreground(success).
			Render("●")
)

type JSONOutput struct {
	SQL     string `json:"sql"`
	Backend string `json:"backend"`
	Elapsed int64  `json:"elapsed_ms"`
	Safe    bool   `json:"safe"`
}

func JSON(w io.Writer, sql, backend string, elapsed time.Duration) {
	out := JSONOutput{
		SQL:     sql,
		Backend: backend,
		Elapsed: elapsed.Milliseconds(),
		Safe:    true,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}

func Pretty(w io.Writer, sql, backend string, elapsed time.Duration) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, highlightSQL(sql))
	fmt.Fprintln(w)
	fmt.Fprintln(w, metaStyle.Render(fmt.Sprintf("%s %s · %dms", successDot, backend, elapsed.Milliseconds())))
}

var keywords = []string{
	"SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "IN", "IS", "NULL",
	"JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "ON", "AS",
	"ORDER", "BY", "ASC", "DESC", "LIMIT", "OFFSET",
	"GROUP", "HAVING", "DISTINCT", "COUNT", "SUM", "AVG", "MIN", "MAX",
	"INSERT", "INTO", "VALUES", "UPDATE", "SET", "DELETE",
	"CREATE", "TABLE", "ALTER", "DROP", "INDEX", "PRIMARY", "KEY",
	"FOREIGN", "REFERENCES", "CONSTRAINT", "DEFAULT", "NOT NULL",
	"TRUE", "FALSE", "BETWEEN", "LIKE", "ILIKE", "EXISTS",
	"UNION", "ALL", "CASE", "WHEN", "THEN", "ELSE", "END",
	"CAST", "COALESCE", "NULLIF", "INTERVAL", "CURRENT_DATE", "CURRENT_TIMESTAMP",
	"WITH", "RECURSIVE", "RETURNING",
}

func highlightSQL(sql string) string {
	result := sql

	for _, kw := range keywords {
		lower := strings.ToLower(kw)
		result = replaceKeyword(result, kw)
		result = replaceKeyword(result, lower)
	}

	return result
}

func replaceKeyword(sql, keyword string) string {
	highlighted := sqlKeyword.Render(strings.ToUpper(keyword))

	words := strings.Fields(sql)
	for i, word := range words {
		clean := strings.TrimRight(word, "(),;")
		suffix := word[len(clean):]

		if strings.EqualFold(clean, keyword) {
			words[i] = highlighted + suffix
		}
	}

	return strings.Join(words, " ")
}
