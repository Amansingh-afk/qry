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

// inputResult holds the result of reading a line from stdin
type inputResult struct {
	line string
	err  error
}

// readLine reads a line from stdin in a goroutine, allowing for interruption
func readLine(scanner *bufio.Scanner) <-chan inputResult {
	ch := make(chan inputResult, 1)
	go func() {
		if scanner.Scan() {
			ch <- inputResult{line: scanner.Text()}
		} else {
			ch <- inputResult{err: scanner.Err()}
		}
	}()
	return ch
}

func runChat(cmd *cobra.Command, args []string) {
	b, err := getBackend()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	model := getModel(b.Name())
	dialect := getDialect()

	ui.ChatStarting(b.Name(), model, workDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println()
		cancel()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// Get existing session or empty string for new session
	sessionID := getSession(b.Name())

	for {
		fmt.Print(ui.Prompt())

		// Read input with ability to interrupt
		select {
		case <-ctx.Done():
			return
		case inputRes := <-readLine(scanner):
			if inputRes.err != nil {
				return
			}

			input := strings.TrimSpace(inputRes.line)
			if input == "" {
				continue
			}

			if input == "exit" || input == "quit" {
				return
			}

			queryCtx, queryCancel := context.WithTimeout(ctx, getTimeout())

			opts := backend.Options{
				Model:     model,
				Dialect:   dialect,
				SessionID: sessionID, // Resume session if available
			}

			// Use full prompt for new sessions, minimal prompt for existing sessions
			var sqlPrompt string
			if sessionID == "" {
				sqlPrompt = prompt.BuildSQL(input, dialect)
			} else {
				sqlPrompt = prompt.BuildFollowUp(input)
			}

			// Show loading spinner with timer
			stopSpinner := ui.Spinner()
			queryRes, err := b.Query(queryCtx, sqlPrompt, workDir, opts)
			duration := stopSpinner()
			ui.QueryDone(duration)

			queryCancel()

			if err != nil {
				if ctx.Err() != nil {
					return
				}
				ui.Error(err.Error())
				continue
			}

			// Store session ID for next turn and persist
			if queryRes.SessionID != "" {
				sessionID = queryRes.SessionID
				saveSession(b.Name(), sessionID)
			}

			sql := prompt.ExtractSQL(queryRes.Response)

			if warning := guardrails.Check(sql); warning != "" {
				ui.Warning(warning)
			}

			output.Pretty(os.Stdout, sql, b.Name(), model)
			fmt.Println()
		}
	}
}
