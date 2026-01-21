package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Gemini struct{}

func (g *Gemini) Name() string {
	return "gemini"
}

func (g *Gemini) Available() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

func (g *Gemini) Query(ctx context.Context, prompt string, workDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "gemini", "-p", prompt)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("gemini: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

var _ Backend = (*Gemini)(nil)
