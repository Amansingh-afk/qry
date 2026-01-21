package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Cursor struct{}

func (c *Cursor) Name() string {
	return "cursor"
}

func (c *Cursor) Available() bool {
	_, err := exec.LookPath("cursor")
	return err == nil
}

func (c *Cursor) Query(ctx context.Context, prompt string, workDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "cursor", "--prompt", prompt)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cursor: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

var _ Backend = (*Cursor)(nil)
