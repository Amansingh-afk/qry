package ui

import (
	"fmt"
	"os"
	"time"

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

// version is set via ldflags at build time
// -X github.com/amansingh-afk/qry/internal/ui.version=v1.0.0
var version = "dev"

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

// Spinner shows an animated loading spinner with elapsed time and returns a stop function
// The stop function returns the elapsed duration
func Spinner() func() time.Duration {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan struct{})
	start := time.Now()

	go func() {
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				elapsed := time.Since(start).Seconds()
				fmt.Print(dimStyle.Render(fmt.Sprintf("\r  %s thinking... %.1fs", frames[i%len(frames)], elapsed)))
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	return func() time.Duration {
		close(done)
		elapsed := time.Since(start)
		fmt.Print("\r\033[K") // Clear the line
		return elapsed
	}
}

func ClearLine() {
	fmt.Print("\r\033[K")
}

// QueryDone shows completion time
func QueryDone(duration time.Duration) {
	fmt.Printf("  %s\n", dimStyle.Render(fmt.Sprintf("%.1fs", duration.Seconds())))
}

func ServerStarting(port int, dir string) {
	fmt.Println()
	fmt.Println(infoStyle.Render(fmt.Sprintf("  QRY Server :%d", port)))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  %s", dir)))
	fmt.Println()
}

func ChatStarting(backend, model, workDir string) {
	fmt.Println()

	// Mini ASCII QRY
	q := cyanStyle.Render("█▀█")
	r := pinkStyle.Render("█▀█")
	y := purpleStyle.Render("█ █")

	q2 := cyanStyle.Render("▀▀█")
	r2 := pinkStyle.Render("█▀▄")
	y2 := purpleStyle.Render("▀█▀")

	fmt.Printf("  %s %s %s\n", q, r, y)
	fmt.Printf("  %s %s %s  %s\n", q2, r2, y2, dimStyle.Render("v"+version))

	backendInfo := backend
	if model != "" {
		backendInfo += "/" + model
	}

	fmt.Println()
	fmt.Printf("  %s  %s\n", dimStyle.Render(backendInfo), dimStyle.Render(workDir))
	fmt.Println()
}

func Prompt() string {
	return pinkStyle.Render("> ")
}

// Step shows a step in progress with arrow
func Step(format string, args ...interface{}) {
	fmt.Printf("  %s %s\n", cyanStyle.Render("→"), fmt.Sprintf(format, args...))
}

// StepDone shows a completed step with checkmark
func StepDone(format string, args ...interface{}) {
	fmt.Printf("    %s\n", okStyle.Render(fmt.Sprintf(format, args...)))
}

// StepItem shows an indented item under a step
func StepItem(format string, args ...interface{}) {
	fmt.Printf("    %s\n", dimStyle.Render(fmt.Sprintf(format, args...)))
}

// StepWarn shows a warning item under a step
func StepWarn(format string, args ...interface{}) {
	fmt.Printf("    %s\n", warnStyle.Render(fmt.Sprintf(format, args...)))
}

// Header shows a styled header
func Header(text string) {
	q := cyanStyle.Render("Q")
	r := pinkStyle.Render("R")
	y := purpleStyle.Render("Y")
	fmt.Printf("\n  %s%s%s %s\n\n", q, r, y, dimStyle.Render(text))
}

// Done shows final success message
func Done(format string, args ...interface{}) {
	fmt.Printf("\n  %s %s\n\n", okStyle.Render("✓"), fmt.Sprintf(format, args...))
}

// Hint shows a hint/suggestion
func Hint(format string, args ...interface{}) {
	fmt.Printf("  %s\n", dimStyle.Render(fmt.Sprintf(format, args...)))
}

// Pause adds a small delay for premium feel
func Pause() {
	time.Sleep(150 * time.Millisecond)
}
