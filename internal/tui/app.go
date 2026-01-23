package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amansingh-afk/qry/internal/history"
	"github.com/amansingh-afk/qry/internal/security"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QueryFunc is the function signature for executing queries
type QueryFunc func(ctx context.Context, query string) (QueryResult, error)

// QueryResult holds the result of a query
type QueryResult struct {
	SQL       string
	SessionID string
	Duration  time.Duration
}

// HistoryItem represents a past query
type HistoryItem struct {
	Query    string
	SQL      string
	Duration time.Duration
	Tables   []string
	Safety   string
}

// Thinking phrases that rotate during loading
var thinkingPhrases = []string{
	"parsing your intent...",
	"consulting the schema...",
	"summoning SQL...",
	"brewing your query...",
	"connecting the dots...",
	"reading the tables...",
	"crafting magic...",
	"almost there...",
}

// Model is the Bubble Tea model for the TUI
type Model struct {
	// Config
	repo    string
	backend string
	model   string
	version string
	workDir string
	width   int
	height  int

	// State
	textInput      textinput.Model
	spinner        spinner.Model
	loading        bool
	loadingStart   time.Time
	err            error
	currentSQL     string
	currentQuery   string // Store the query text for history
	currentTime    time.Duration
	tables         []string
	safety         string
	expanded       bool
	history        []HistoryItem
	historyIdx     int
	showHistory    bool
	showHelp       bool
	copied         bool
	historyCleared bool

	// Security
	securityResult  *security.Result
	securityBlocked bool

	// Query execution
	queryFunc QueryFunc
	queryCtx  context.Context
	cancelFn  context.CancelFunc
}

// queryResultMsg is sent when a query completes
type queryResultMsg struct {
	result QueryResult
	err    error
}

// copyResetMsg resets the copy indicator
type copyResetMsg struct{}

// historyClearedResetMsg resets the history cleared indicator
type historyClearedResetMsg struct{}

// NewModel creates a new TUI model
func NewModel(repo, backend, model, version, workDir string, queryFunc QueryFunc) Model {
	ti := textinput.New()
	ti.Placeholder = "Ask a question..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 200 // Large width, container will clip
	ti.PromptStyle = promptStyle
	ti.TextStyle = inputStyle
	ti.Prompt = "❯ "

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	// Load history from disk
	var historyItems []HistoryItem
	if entries, err := history.Load(workDir); err == nil {
		for _, e := range entries {
			historyItems = append(historyItems, HistoryItem{
				Query:    e.Query,
				SQL:      e.SQL,
				Duration: time.Duration(e.Duration * float64(time.Second)),
			})
		}
	}

	return Model{
		repo:       repo,
		backend:    backend,
		model:      model,
		version:    version,
		workDir:    workDir,
		textInput:  ti,
		spinner:    s,
		queryFunc:  queryFunc,
		history:    historyItems,
		historyIdx: -1,
		width:      120,
		height:     24,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		key := msg.String()

		// Handle ctrl+c always
		if key == "ctrl+c" {
			if m.loading && m.cancelFn != nil {
				m.cancelFn()
				m.loading = false
				return m, nil
			}
			return m, tea.Quit
		}

		// Handle esc - close history panel
		if key == "esc" {
			m.showHistory = false
			m.historyIdx = -1
			return m, nil
		}

		// Handle enter for query submission or vim commands
		if key == "enter" {
			if m.loading {
				return m, nil
			}
			query := strings.TrimSpace(m.textInput.Value())
			if query == "" {
				return m, nil
			}

			// Handle vim-style colon commands
			if strings.HasPrefix(query, ":") {
				cmd := strings.ToLower(strings.TrimPrefix(query, ":"))
				m.textInput.SetValue("")

				switch cmd {
				case "q", "quit", "exit":
					return m, tea.Quit

				case "c", "copy":
					if m.currentSQL != "" {
						if err := clipboard.WriteAll(m.currentSQL); err == nil {
							m.copied = true
							return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
								return copyResetMsg{}
							})
						}
					}
					return m, nil

				case "h", "history":
					m.showHistory = !m.showHistory
					return m, nil

				case "e", "expand":
					m.expanded = !m.expanded
					return m, nil

				case "clear":
					m.currentSQL = ""
					m.err = nil
					m.tables = nil
					m.safety = ""
					m.showHistory = false
					return m, nil

				case "clear-history":
					m.history = []HistoryItem{}
					m.historyIdx = -1
					m.showHistory = false
					m.historyCleared = true
					_ = history.Clear(m.workDir)
					return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
						return historyClearedResetMsg{}
					})

				case "help", "?":
					m.showHelp = !m.showHelp
					return m, nil

				default:
					// Unknown command, show error briefly
					m.err = fmt.Errorf("unknown command: %s", cmd)
					return m, nil
				}
			}

			// Regular query
			if query == "exit" || query == "quit" {
				return m, tea.Quit
			}

			// Start query
			m.loading = true
			m.loadingStart = time.Now()
			m.currentQuery = query // Store for history
			m.err = nil
			m.currentSQL = ""
			m.copied = false
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			m.queryCtx = ctx
			m.cancelFn = cancel
			return m, tea.Batch(m.spinner.Tick, m.executeQuery(query))
		}

		// Handle up/down for history navigation (always works)
		if key == "up" && len(m.history) > 0 && !m.loading {
			if m.historyIdx < len(m.history)-1 {
				m.historyIdx++
				m.textInput.SetValue(m.history[len(m.history)-1-m.historyIdx].Query)
				m.textInput.CursorEnd()
			}
			return m, nil
		}

		if key == "down" && !m.loading {
			if m.historyIdx > 0 {
				m.historyIdx--
				m.textInput.SetValue(m.history[len(m.history)-1-m.historyIdx].Query)
				m.textInput.CursorEnd()
			} else if m.historyIdx == 0 {
				m.historyIdx = -1
				m.textInput.SetValue("")
			}
			return m, nil
		}

	case queryResultMsg:
		m.loading = false
		if m.cancelFn != nil {
			m.cancelFn()
		}
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		// Security validation
		m.securityResult = security.Validate(msg.result.SQL)
		sec := security.Get()
		m.securityBlocked = sec.IsBlocked(m.securityResult)

		// If blocked by security, don't show SQL
		if m.securityBlocked {
			m.currentSQL = ""
			m.err = fmt.Errorf("%s", m.securityResult.Error())
			m.textInput.SetValue("")
			return m, nil
		}

		m.currentSQL = msg.result.SQL
		m.currentTime = msg.result.Duration
		m.tables = extractTables(msg.result.SQL)
		m.safety = checkSafety(msg.result.SQL)

		// Add to in-memory history
		item := HistoryItem{
			Query:    m.currentQuery,
			SQL:      msg.result.SQL,
			Duration: msg.result.Duration,
			Tables:   m.tables,
			Safety:   m.safety,
		}
		m.history = append(m.history, item)
		m.historyIdx = -1
		m.textInput.SetValue("")

		// Persist to disk
		_ = history.Add(m.workDir, m.currentQuery, msg.result.SQL, msg.result.Duration)

	case copyResetMsg:
		m.copied = false

	case historyClearedResetMsg:
		m.historyCleared = false

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update text input
	if !m.loading {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// executeQuery runs the query in the background
func (m *Model) executeQuery(query string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		result, err := m.queryFunc(m.queryCtx, query)
		result.Duration = time.Since(start)
		return queryResultMsg{result: result, err: err}
	}
}

// View renders the UI
func (m Model) View() string {
	var b strings.Builder

	// Calculate width for content
	contentWidth := min(m.width-4, 80)

	// Header
	header := m.renderHeader(contentWidth)
	b.WriteString(header)
	b.WriteString("\n")

	// Separator
	b.WriteString(separatorStyle.Render(strings.Repeat("─", contentWidth)))
	b.WriteString("\n")

	// Input section
	if m.loading {
		elapsed := time.Since(m.loadingStart)
		// Rotate phrases every 3 seconds
		phraseIdx := int(elapsed.Seconds()/3) % len(thinkingPhrases)
		phrase := thinkingPhrases[phraseIdx]
		timeStr := fmt.Sprintf("%.1fs", elapsed.Seconds())
		b.WriteString(fmt.Sprintf(" %s %s %s\n", m.spinner.View(), dimStyle.Render(phrase), timerStyle.Render(timeStr)))
	} else {
		b.WriteString(" ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")
	}

	// SQL output
	if m.currentSQL != "" || m.err != nil {
		b.WriteString("\n")
		b.WriteString(m.renderSQL(contentWidth))
	}

	// History view
	if m.showHistory && len(m.history) > 0 {
		b.WriteString("\n")
		b.WriteString(m.renderHistory(contentWidth))
	}

	// Help view
	if m.showHelp {
		b.WriteString("\n")
		b.WriteString(m.renderHelp())
	}

	// Error
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf(" Error: %s", m.err.Error())))
		b.WriteString("\n")
	}

	// Metadata (only when there's SQL)
	if m.currentSQL != "" {
		b.WriteString("\n")
		b.WriteString(m.renderMetadata(contentWidth))
	}

	// Footer (always shown)
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", contentWidth)))
	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	// Wrap in container
	content := b.String()
	container := containerStyle.Width(contentWidth + 2).Render(content)

	return "\n" + container + "\n"
}

func (m Model) renderHeader(width int) string {
	// QRY logo - bigger and bolder
	logo := logoQ.Render("Q") + logoR.Render("R") + logoY.Render("Y")

	// Version
	ver := dimStyle.Render("v" + m.version)

	// Info
	info := fmt.Sprintf("%s | %s", m.repo, m.backend)
	if m.model != "" {
		info = fmt.Sprintf("%s | %s/%s", m.repo, m.backend, m.model)
	}

	// Layout: QRY v0.3.1                    repo | claude/haiku
	leftPart := " " + logo + " " + ver
	rightPart := dimStyle.Render(info)

	padding := width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart) - 1
	if padding < 1 {
		padding = 1
	}

	return leftPart + strings.Repeat(" ", padding) + rightPart
}

func (m Model) renderSQL(width int) string {
	var b strings.Builder

	b.WriteString(sqlHeaderStyle.Render(" Generated SQL:"))
	b.WriteString("\n")
	b.WriteString(sqlLineStyle.Render(" " + strings.Repeat("─", width-2)))
	b.WriteString("\n")

	// Syntax highlight and format SQL
	sql := m.currentSQL
	if !m.expanded && len(sql) > 500 {
		sql = sql[:500] + "..."
	}

	highlighted := highlightSQL(sql)
	lines := strings.Split(highlighted, "\n")
	for _, line := range lines {
		b.WriteString(" " + line + "\n")
	}

	return b.String()
}

func (m Model) renderMetadata(width int) string {
	var parts []string

	// Timer
	if m.currentTime > 0 {
		parts = append(parts, timerStyle.Render(fmt.Sprintf("⏱ %.1fs", m.currentTime.Seconds())))
	}

	// Tables
	if len(m.tables) > 0 {
		tables := strings.Join(m.tables, ", ")
		parts = append(parts, metaLabelStyle.Render("Tables: ")+metaValueStyle.Render(tables))
	}

	// Safety
	var safetyStr string
	switch m.safety {
	case "OK":
		safetyStr = safetyOK.Render("✓ READ-ONLY")
	case "WARN":
		safetyStr = safetyWarn.Render("⚠ MODIFIES DATA")
	case "DANGER":
		safetyStr = safetyDanger.Render("✗ DESTRUCTIVE")
	}
	parts = append(parts, safetyStr)

	// Security warning (warn mode)
	if m.securityResult != nil && !m.securityResult.Valid && !m.securityBlocked {
		parts = append(parts, safetyWarn.Render("⚠ SECURITY WARNING"))
	}

	// Copied indicator
	if m.copied {
		parts = append(parts, safetyOK.Render("✓ Copied!"))
	}

	// History cleared indicator
	if m.historyCleared {
		parts = append(parts, safetyOK.Render("✓ History cleared"))
	}

	return " " + strings.Join(parts, "  |  ")
}

func (m Model) renderFooter() string {
	var shortcuts []struct {
		key  string
		desc string
	}

	// During loading, show cancel option
	if m.loading {
		shortcuts = append(shortcuts, struct {
			key  string
			desc string
		}{"^C", "cancel"})

		var parts []string
		for _, s := range shortcuts {
			parts = append(parts, shortcutKeyStyle.Render(s.key)+" "+shortcutDescStyle.Render(s.desc))
		}
		return " " + strings.Join(parts, "  ")
	}

	// Normal mode shortcuts
	shortcuts = append(shortcuts, struct {
		key  string
		desc string
	}{":h", "history"})

	shortcuts = append(shortcuts, struct {
		key  string
		desc string
	}{"↑↓", "prev"})

	// Show copy only when there's SQL
	if m.currentSQL != "" {
		shortcuts = append(shortcuts, struct {
			key  string
			desc string
		}{":c", "copy"})
	}

	// Show expand only for long SQL
	if m.currentSQL != "" && len(m.currentSQL) > 500 {
		shortcuts = append(shortcuts, struct {
			key  string
			desc string
		}{":e", "expand"})
	}

	// Help and quit
	shortcuts = append(shortcuts, struct {
		key  string
		desc string
	}{":?", "help"})

	shortcuts = append(shortcuts, struct {
		key  string
		desc string
	}{":q", "quit"})

	var parts []string
	for _, s := range shortcuts {
		parts = append(parts, shortcutKeyStyle.Render(s.key)+" "+shortcutDescStyle.Render(s.desc))
	}

	return " " + strings.Join(parts, "  ")
}

func (m Model) renderHistory(width int) string {
	var b strings.Builder

	b.WriteString(sqlHeaderStyle.Render(" History:"))
	b.WriteString("\n")

	// Show last 5 items
	start := 0
	if len(m.history) > 5 {
		start = len(m.history) - 5
	}

	for i := len(m.history) - 1; i >= start; i-- {
		item := m.history[i]
		query := item.Query
		if len(query) > 40 {
			query = query[:40] + "..."
		}
		b.WriteString(fmt.Sprintf(" %s %s\n", dimStyle.Render("·"), metaValueStyle.Render(query)))
	}

	return b.String()
}

func (m Model) renderHelp() string {
	var b strings.Builder

	b.WriteString(sqlHeaderStyle.Render(" Commands:"))
	b.WriteString("\n")

	commands := []struct {
		cmd  string
		desc string
	}{
		{":c, :copy", "Copy SQL to clipboard"},
		{":h, :history", "Toggle history panel"},
		{":e, :expand", "Expand/collapse long SQL"},
		{":clear", "Clear current result"},
		{":clear-history", "Wipe all saved history"},
		{":help, :?", "Show this help"},
		{":q, :quit", "Exit"},
		{"↑ / ↓", "Navigate query history"},
		{"Ctrl+C", "Cancel query / exit"},
	}

	for _, c := range commands {
		b.WriteString(fmt.Sprintf(" %s  %s\n",
			shortcutKeyStyle.Render(fmt.Sprintf("%-15s", c.cmd)),
			dimStyle.Render(c.desc)))
	}

	return b.String()
}

// highlightSQL applies syntax highlighting to SQL
func highlightSQL(sql string) string {
	words := tokenizeSQL(sql)
	var result strings.Builder

	kwStyle := lipgloss.NewStyle().Foreground(sqlKeyword).Bold(true)
	fnStyle := lipgloss.NewStyle().Foreground(sqlFunc)
	strStyle := lipgloss.NewStyle().Foreground(sqlString)
	numStyle := lipgloss.NewStyle().Foreground(sqlNumber)
	defaultStyle := lipgloss.NewStyle().Foreground(white)

	for _, word := range words {
		upper := strings.ToUpper(word)

		switch {
		case sqlKeywords[upper]:
			result.WriteString(kwStyle.Render(word))
		case sqlFunctions[upper]:
			result.WriteString(fnStyle.Render(word))
		case isStringLiteral(word):
			result.WriteString(strStyle.Render(word))
		case isNumber(word):
			result.WriteString(numStyle.Render(word))
		default:
			result.WriteString(defaultStyle.Render(word))
		}
	}

	return result.String()
}

// tokenizeSQL splits SQL into tokens while preserving whitespace
func tokenizeSQL(sql string) []string {
	var tokens []string
	var current strings.Builder

	inString := false
	stringChar := byte(0)

	for i := 0; i < len(sql); i++ {
		c := sql[i]

		if inString {
			current.WriteByte(c)
			if c == stringChar {
				inString = false
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		switch {
		case c == '\'' || c == '"':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			inString = true
			stringChar = c
			current.WriteByte(c)

		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(c))

		case c == '(' || c == ')' || c == ',' || c == ';' || c == '.' || c == '=' || c == '<' || c == '>' || c == '*':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(c))

		default:
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func isStringLiteral(s string) bool {
	return len(s) >= 2 && (s[0] == '\'' || s[0] == '"')
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

// extractTables extracts table names from SQL
func extractTables(sql string) []string {
	tables := make(map[string]bool)
	upper := strings.ToUpper(sql)
	words := strings.Fields(sql)

	for i, word := range words {
		upperWord := strings.ToUpper(word)
		if upperWord == "FROM" || upperWord == "JOIN" || upperWord == "INTO" || upperWord == "UPDATE" {
			if i+1 < len(words) {
				table := strings.Trim(words[i+1], "(),;")
				if table != "" && !sqlKeywords[strings.ToUpper(table)] {
					tables[table] = true
				}
			}
		}
	}

	// Also check for table aliases after JOIN
	_ = upper // suppress unused warning

	var result []string
	for t := range tables {
		result = append(result, t)
	}
	return result
}

// checkSafety determines if the SQL is safe
func checkSafety(sql string) string {
	upper := strings.ToUpper(sql)

	// Destructive operations
	if strings.Contains(upper, "DROP ") ||
		strings.Contains(upper, "TRUNCATE ") ||
		strings.Contains(upper, "DELETE ") && !strings.Contains(upper, "WHERE") {
		return "DANGER"
	}

	// Modifying operations
	if strings.Contains(upper, "INSERT ") ||
		strings.Contains(upper, "UPDATE ") ||
		strings.Contains(upper, "DELETE ") ||
		strings.Contains(upper, "ALTER ") ||
		strings.Contains(upper, "CREATE ") {
		return "WARN"
	}

	return "OK"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
