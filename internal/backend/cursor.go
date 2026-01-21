package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Cursor struct{}

func (c *Cursor) Name() string { return "cursor" }

func (c *Cursor) InstallCmd() string {
	return "curl -fsSL https://cursor.com/install | sh"
}

func (c *Cursor) Available() bool {
	cmd := exec.Command("cursor", "--version")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "cursor")
}

func (c *Cursor) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	args := []string{"--prompt", prompt}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "cursor", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("cursor: %w\n%s", err, string(out))
	}
	return Result{Response: strings.TrimSpace(string(out))}, nil
}
