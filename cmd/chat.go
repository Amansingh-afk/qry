package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/guardrails"
	"github.com/amansingh-afk/qry/internal/output"
	"github.com/amansingh-afk/qry/internal/prompt"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive SQL chat session",
	Long: `Start an interactive chat session for SQL generation.
Type your questions, get SQL. Ctrl+C to exit.`,
	Example: `  qry chat
  qry chat -m sonnet
  qry chat -d postgresql`,
	Run: runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

func runChat(cmd *cobra.Command, args []string) {
	b, err := getBackend()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	model := getModel(b.Name())
	dialect := getDialect()

	ui.ChatStarting(b.Name(), model)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	scanner := bufio.NewScanner(os.Stdin)

	var sessionID string // Track session for multi-turn conversation

	for {
		fmt.Print(ui.Prompt())

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			break
		}

		queryCtx, queryCancel := context.WithTimeout(ctx, getTimeout())

		opts := backend.Options{
			Model:     model,
			Dialect:   dialect,
			SessionID: sessionID, // Resume session if available
		}

		sqlPrompt := prompt.BuildSQL(input, dialect)
		result, err := b.Query(queryCtx, sqlPrompt, workDir, opts)

		queryCancel()

		if err != nil {
			if ctx.Err() != nil {
				fmt.Println()
				break
			}
			ui.Error(err.Error())
			continue
		}

		// Store session ID for next turn
		if result.SessionID != "" {
			sessionID = result.SessionID
		}

		sql := prompt.ExtractSQL(result.Response)

		if warning := guardrails.Check(sql); warning != "" {
			ui.Warning(warning)
		}

		output.Pretty(os.Stdout, sql, b.Name(), model)
		fmt.Println()
	}
}
