package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Codex struct{}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) InstallCmd() string {
	return "npm i -g @openai/codex"
}

func (c *Codex) Available() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

func (c *Codex) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	args := []string{"-p", prompt}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "codex", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("codex: %w\n%s", err, string(out))
	}
	return Result{Response: strings.TrimSpace(string(out))}, nil
}
