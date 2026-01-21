package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	primary   = lipgloss.Color("#00D9FF")
	secondary = lipgloss.Color("#FF6AC1")
	success   = lipgloss.Color("#5AF78E")
	warning   = lipgloss.Color("#F3F99D")
	danger    = lipgloss.Color("#FF5C57")
	muted     = lipgloss.Color("#636363")
	white     = lipgloss.Color("#F1F1F0")

	logoStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	accentStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true)

	taglineStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warning)

	successStyle = lipgloss.NewStyle().
			Foreground(success)

	infoStyle = lipgloss.NewStyle().
			Foreground(primary)

	dimStyle = lipgloss.NewStyle().
			Foreground(muted)

	sqlStyle = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(0, 1).
			Bold(true)
)

const version = "0.1.0"

func Banner() string {
	q := logoStyle.Render("Q")
	r := accentStyle.Render("R")
	y := logoStyle.Render("Y")

	logo := fmt.Sprintf(`
    %s%s%s
`, q, r, y)

	tagline := taglineStyle.Render("    Ask. Get SQL.\n")

	return logo + tagline
}

func BannerFull() string {
	line1 := logoStyle.Render("  ╔═══╗") + accentStyle.Render("╔═══╗") + logoStyle.Render("╗   ╗")
	line2 := logoStyle.Render("  ║   ║") + accentStyle.Render("║   ║") + logoStyle.Render("║   ║")
	line3 := logoStyle.Render("  ║   ║") + accentStyle.Render("╠═══╝") + logoStyle.Render(" ╚═╦═╝")
	line4 := logoStyle.Render("  ║ ╔╗║") + accentStyle.Render("║  ╚╗") + logoStyle.Render("   ║")
	line5 := logoStyle.Render("  ╚═╝╚╝") + accentStyle.Render("╝   ╚") + logoStyle.Render("   ╝")

	logo := fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n", line1, line2, line3, line4, line5)
	tagline := taglineStyle.Render("\n    Ask. Get SQL.\n")

	return logo + tagline
}

func Version() string {
	return fmt.Sprintf("%s %s", logoStyle.Render("qry"), dimStyle.Render(version))
}

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, errorStyle.Render("✗ "+msg))
}

func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, warningStyle.Render("⚠ "+msg))
}

func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(successStyle.Render("✓ " + msg))
}

func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(infoStyle.Render("→ " + msg))
}

func Thinking(backend string) {
	fmt.Print(dimStyle.Render(fmt.Sprintf("● querying %s ", backend)))
}

func ClearLine() {
	fmt.Print("\r\033[K")
}

func ServerStarting(port int) {
	fmt.Println()
	fmt.Println(headerStyle.Render(" QRY SERVER "))
	fmt.Println()
	fmt.Println(infoStyle.Render(fmt.Sprintf("  ➜ http://localhost:%d", port)))
	fmt.Println(dimStyle.Render("  Press Ctrl+C to stop"))
	fmt.Println()
}

func SQL(sql string) string {
	return sqlStyle.Render(sql)
}

func Muted(s string) string {
	return dimStyle.Render(s)
}
