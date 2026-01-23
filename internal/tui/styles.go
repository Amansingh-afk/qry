package tui

import "github.com/charmbracelet/lipgloss"

// Colors - Dracula-inspired palette
var (
	cyan   = lipgloss.Color("#00FFFF")
	purple = lipgloss.Color("#BD93F9")
	pink   = lipgloss.Color("#FF79C6")
	green  = lipgloss.Color("#50FA7B")
	yellow = lipgloss.Color("#F1FA8C")
	red    = lipgloss.Color("#FF5555")
	gray   = lipgloss.Color("#6272A4")
	white  = lipgloss.Color("#F8F8F2")

	// SQL syntax highlighting colors
	sqlKeyword = lipgloss.Color("#FF79C6") // pink
	sqlString  = lipgloss.Color("#F1FA8C") // yellow
	sqlNumber  = lipgloss.Color("#BD93F9") // purple
	sqlFunc    = lipgloss.Color("#50FA7B") // green
)

// Box styles
var (
	// Main container border
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(0, 1)

	// QRY logo in header
	logoQ = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	logoR = lipgloss.NewStyle().Foreground(pink).Bold(true)
	logoY = lipgloss.NewStyle().Foreground(purple).Bold(true)

	// Separator line
	separatorStyle = lipgloss.NewStyle().Foreground(gray)

	// Prompt style
	promptStyle = lipgloss.NewStyle().Foreground(pink).Bold(true)
	inputStyle  = lipgloss.NewStyle().Foreground(white)

	// SQL output section
	sqlHeaderStyle = lipgloss.NewStyle().Foreground(gray).Italic(true)
	sqlLineStyle   = lipgloss.NewStyle().Foreground(gray)

	// Metadata section
	metaLabelStyle = lipgloss.NewStyle().Foreground(gray)
	metaValueStyle = lipgloss.NewStyle().Foreground(white)
	safetyOK       = lipgloss.NewStyle().Foreground(green).Bold(true)
	safetyWarn     = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	safetyDanger   = lipgloss.NewStyle().Foreground(red).Bold(true)

	// Footer shortcuts
	shortcutKeyStyle  = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	shortcutDescStyle = lipgloss.NewStyle().Foreground(gray)

	// Timing
	timerStyle = lipgloss.NewStyle().Foreground(gray)

	// Error
	errorStyle = lipgloss.NewStyle().Foreground(red)

	// Spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(cyan)

	// Dim text
	dimStyle = lipgloss.NewStyle().Foreground(gray)
)

// SQL Keywords for syntax highlighting
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
	"DISTINCT": true, "COUNT": true, "SUM": true, "AVG": true,
	"MIN": true, "MAX": true, "CASE": true, "WHEN": true, "THEN": true,
	"ELSE": true, "END": true, "UNION": true, "ALL": true,
	"WITH": true, "OVER": true, "PARTITION": true, "ROW_NUMBER": true,
	"RANK": true, "DENSE_RANK": true, "LAG": true, "LEAD": true,
	"COALESCE": true, "NULLIF": true, "CAST": true, "INTERVAL": true,
	"DATE_TRUNC": true, "NOW": true, "CURRENT_DATE": true, "CURRENT_TIMESTAMP": true,
}

// SQL functions for highlighting
var sqlFunctions = map[string]bool{
	"COUNT": true, "SUM": true, "AVG": true, "MIN": true, "MAX": true,
	"COALESCE": true, "NULLIF": true, "CAST": true, "CONCAT": true,
	"LOWER": true, "UPPER": true, "TRIM": true, "LENGTH": true,
	"SUBSTRING": true, "REPLACE": true, "DATE_TRUNC": true, "NOW": true,
	"EXTRACT": true, "TO_CHAR": true, "TO_DATE": true, "ROUND": true,
	"FLOOR": true, "CEIL": true, "ABS": true, "ROW_NUMBER": true,
	"RANK": true, "DENSE_RANK": true, "LAG": true, "LEAD": true,
}
