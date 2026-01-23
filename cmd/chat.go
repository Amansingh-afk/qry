package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/prompt"
	"github.com/amansingh-afk/qry/internal/tui"
	"github.com/amansingh-afk/qry/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func runChat(cmd *cobra.Command, args []string) {
	b, err := getBackend()
	if err != nil {
		ui.Error("%s", err.Error())
		os.Exit(1)
	}

	model := getModel(b.Name())
	dialect := getDialect()

	// Get repo name from directory
	repo := filepath.Base(workDir)

	// Get existing session or empty string for new session
	sessionID := getSession(b.Name())

	// Create query function that the TUI will call
	queryFunc := func(ctx context.Context, query string) (tui.QueryResult, error) {
		opts := backend.Options{
			Model:     model,
			Dialect:   dialect,
			SessionID: sessionID,
		}

		// Use full prompt for new sessions, minimal prompt for existing sessions
		var sqlPrompt string
		if sessionID == "" {
			sqlPrompt = prompt.BuildSQL(query, dialect)
		} else {
			sqlPrompt = prompt.BuildFollowUp(query)
		}

		result, err := b.Query(ctx, sqlPrompt, workDir, opts)
		if err != nil {
			return tui.QueryResult{}, err
		}

		// Update session for next query
		if result.SessionID != "" {
			sessionID = result.SessionID
			saveSession(b.Name(), sessionID)
		}

		sql := prompt.ExtractSQL(result.Response)

		return tui.QueryResult{
			SQL:       sql,
			SessionID: result.SessionID,
		}, nil
	}

	// Create and run TUI
	version := strings.TrimPrefix(ui.Version(), "qry ")
	m := tui.NewModel(repo, b.Name(), model, version, workDir, queryFunc)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
