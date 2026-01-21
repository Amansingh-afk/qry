package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	cyan   = lipgloss.Color("#00FFFF")
	purple = lipgloss.Color("#BD93F9")
	pink   = lipgloss.Color("#FF79C6")
	green  = lipgloss.Color("#50FA7B")
	yellow = lipgloss.Color("#F1FA8C")
	red    = lipgloss.Color("#FF5555")
	gray   = lipgloss.Color("#6272A4")

	cyanStyle   = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	purpleStyle = lipgloss.NewStyle().Foreground(purple).Bold(true)
	pinkStyle   = lipgloss.NewStyle().Foreground(pink).Bold(true)
	dimStyle    = lipgloss.NewStyle().Foreground(gray)
	errStyle    = lipgloss.NewStyle().Foreground(red).Bold(true)
	warnStyle   = lipgloss.NewStyle().Foreground(yellow)
	okStyle     = lipgloss.NewStyle().Foreground(green)
	infoStyle   = lipgloss.NewStyle().Foreground(cyan)
)

const version = "0.2.0"

func Banner() string {
	q := cyanStyle.Render(`
   ██████╗ 
  ██╔═══██╗
  ██║   ██║
  ██║▄▄ ██║
  ╚██████╔╝
   ╚══▀▀═╝ `)

	r := pinkStyle.Render(`
  ██████╗ 
  ██╔══██╗
  ██████╔╝
  ██╔══██╗
  ██║  ██║
  ╚═╝  ╚═╝`)

	y := purpleStyle.Render(`
  ██╗   ██╗
  ╚██╗ ██╔╝
   ╚████╔╝ 
    ╚██╔╝  
     ██║   
     ╚═╝   `)

	qLines := splitLines(q)
	rLines := splitLines(r)
	yLines := splitLines(y)

	var banner string
	for i := 0; i < len(qLines); i++ {
		banner += qLines[i] + rLines[i] + yLines[i] + "\n"
	}

	tagline := dimStyle.Render("  Ask. Get SQL.\n")

	return "\n" + banner + tagline
}

func splitLines(s string) []string {
	var lines []string
	var line string
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(c)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func Version() string {
	return fmt.Sprintf("qry %s", version)
}

func Print(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Error(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, errStyle.Render("✗ "+fmt.Sprintf(format, args...)))
}

func Warning(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, warnStyle.Render("⚠ "+fmt.Sprintf(format, args...)))
}

func Success(format string, args ...interface{}) {
	fmt.Println(okStyle.Render("✓ " + fmt.Sprintf(format, args...)))
}

func Info(format string, args ...interface{}) {
	fmt.Println(infoStyle.Render("→ " + fmt.Sprintf(format, args...)))
}

func Thinking(backend string) {
	fmt.Print(dimStyle.Render(fmt.Sprintf("● %s ", backend)))
}

func ClearLine() {
	fmt.Print("\r\033[K")
}

func ServerStarting(port int, dir string) {
	fmt.Println()
	fmt.Println(infoStyle.Render(fmt.Sprintf("  QRY Server :%d", port)))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  %s", dir)))
	fmt.Println()
}

func ChatStarting(backend, model string) {
	fmt.Println()
	info := infoStyle.Render("QRY Chat")
	backendInfo := backend
	if model != "" {
		backendInfo += "/" + model
	} else {
		backendInfo += " (default)"
	}
	fmt.Printf("  %s %s\n", info, dimStyle.Render(backendInfo))
	fmt.Println(dimStyle.Render("  Type queries. 'exit' or Ctrl+C to quit."))
	fmt.Println()
}

func Prompt() string {
	return pinkStyle.Render("> ")
}
