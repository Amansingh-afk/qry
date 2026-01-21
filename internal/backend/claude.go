package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Claude struct{}

func (c *Claude) Name() string {
	return "claude"
}

func (c *Claude) Available() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

func (c *Claude) Query(ctx context.Context, prompt string, workDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", "-p", prompt)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

var _ Backend = (*Claude)(nil)
