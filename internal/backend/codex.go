package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Codex struct{}

func (c *Codex) Name() string {
	return "codex"
}

func (c *Codex) Available() bool {
	_, err := exec.LookPath("codex")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("npx")
	return err == nil
}

func (c *Codex) Query(ctx context.Context, prompt string, workDir string) (string, error) {
	var cmd *exec.Cmd

	if _, err := exec.LookPath("codex"); err == nil {
		cmd = exec.CommandContext(ctx, "codex", "-p", prompt)
	} else {
		cmd = exec.CommandContext(ctx, "npx", "-y", "codex", "-p", prompt)
	}

	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("codex: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

var _ Backend = (*Codex)(nil)
