package backend

import (
	"context"
	"encoding/json"
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

// codexJSONResponse represents the JSON output from codex CLI
type codexJSONResponse struct {
	SessionID string `json:"session_id"`
	RolloutID string `json:"rollout_id"`
	Result    string `json:"result"`
	Response  string `json:"response"`
}

func (c *Codex) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	var args []string

	// Codex uses different syntax for resume: `codex exec resume <id> "prompt"`
	if opts.SessionID != "" {
		args = []string{"exec", "resume", opts.SessionID, prompt, "--output-format", "json"}
	} else {
		args = []string{"exec", "-p", prompt, "--output-format", "json"}
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "codex", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("codex: %w\n%s", err, string(out))
	}

	// Try to parse JSON response
	var resp codexJSONResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		// Fallback: treat as plain text
		return Result{Response: strings.TrimSpace(string(out))}, nil
	}

	// Get session ID from response (try multiple field names)
	sessionID := resp.SessionID
	if sessionID == "" {
		sessionID = resp.RolloutID
	}

	// Get response content
	response := resp.Result
	if response == "" {
		response = resp.Response
	}
	if response == "" {
		response = strings.TrimSpace(string(out))
	}

	return Result{
		Response:  response,
		SessionID: sessionID,
	}, nil
}
