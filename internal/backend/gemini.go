package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Gemini struct{}

func (g *Gemini) Name() string { return "gemini" }

func (g *Gemini) InstallCmd() string {
	return "npm i -g @google/gemini-cli"
}

func (g *Gemini) Available() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

func (g *Gemini) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	args := []string{"-p", prompt}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "gemini", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("gemini: %w\n%s", err, string(out))
	}
	return Result{Response: strings.TrimSpace(string(out))}, nil
}
