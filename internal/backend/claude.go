package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Claude struct{}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) InstallCmd() string {
	return "npm i -g @anthropic-ai/claude-code"
}

func (c *Claude) Available() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// claudeJSONResponse represents the JSON output from claude CLI
type claudeJSONResponse struct {
	SessionID string `json:"session_id"`
	Result    string `json:"result"`
}

func (c *Claude) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	args := []string{"-p", prompt, "--output-format", "json"}

	if opts.SessionID != "" {
		args = append(args, "--resume", opts.SessionID)
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("claude: %w\n%s", err, string(out))
	}

	// Parse JSON response
	var resp claudeJSONResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		// Fallback: treat as plain text (backward compatibility)
		return Result{Response: strings.TrimSpace(string(out))}, nil
	}

	return Result{
		Response:  resp.Result,
		SessionID: resp.SessionID,
	}, nil
}
