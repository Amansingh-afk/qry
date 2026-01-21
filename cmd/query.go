package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/guardrails"
	"github.com/amansingh-afk/qry/internal/output"
	"github.com/amansingh-afk/qry/internal/prompt"
	"github.com/amansingh-afk/qry/internal/ui"
)

func runQuery(query string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := backend.Get(getBackend())
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	if !b.Available() {
		ui.Error("%s is not installed. Install it first.", b.Name())
		os.Exit(1)
	}

	ui.Thinking(b.Name())

	start := time.Now()
	sqlPrompt := prompt.BuildSQL(query)
	result, err := b.Query(ctx, sqlPrompt, getWorkDir())
	elapsed := time.Since(start)

	ui.ClearLine()

	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	sql := prompt.ExtractSQL(result)

	if warning := guardrails.Check(sql); warning != "" {
		ui.Warning(warning)
	}

	if jsonOut {
		output.JSON(os.Stdout, sql, b.Name(), elapsed)
	} else {
		output.Pretty(os.Stdout, sql, b.Name(), elapsed)
	}
}
