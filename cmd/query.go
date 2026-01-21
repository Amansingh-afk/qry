package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/guardrails"
	"github.com/amansingh-afk/qry/internal/output"
	"github.com/amansingh-afk/qry/internal/prompt"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "q [query]",
	Aliases: []string{"query"},
	Short:   "Generate SQL from natural language",
	Example: `  qry q "get active users"
  qry q "count orders" --json
  qry q "find users" -b gemini -m pro
  qry q "get recent" -d postgresql`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runQuery(args[0])
	},
}

func init() {
	queryCmd.Flags().BoolVar(&jsonFlag, "json", false, "output JSON")
	queryCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "show prompt without running")
}

func runQuery(query string) {
	b, err := getBackend()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	dialect := getDialect()

	// Get existing session or empty string for new session
	sessionID := getSession(b.Name())

	// Use full prompt for new sessions, minimal prompt for existing sessions
	var sqlPrompt string
	if sessionID == "" {
		sqlPrompt = prompt.BuildSQL(query, dialect)
	} else {
		sqlPrompt = prompt.BuildFollowUp(query)
	}

	if dryRunFlag {
		ui.Info("Prompt:")
		fmt.Println(sqlPrompt)
		return
	}

	model := getModel(b.Name())

	opts := backend.Options{
		Model:     model,
		Dialect:   dialect,
		SessionID: sessionID,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, getTimeout())
	defer cancelTimeout()

	ui.Thinking(b.Name())

	result, err := b.Query(ctx, sqlPrompt, workDir, opts)

	ui.ClearLine()

	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	// Save session for future queries
	saveSession(b.Name(), result.SessionID)

	sql := prompt.ExtractSQL(result.Response)

	if warning := guardrails.Check(sql); warning != "" {
		ui.Warning(warning)
	}

	if jsonFlag {
		output.JSON(os.Stdout, sql, b.Name(), model, dialect)
	} else {
		output.Pretty(os.Stdout, sql, b.Name(), model)
	}
}
