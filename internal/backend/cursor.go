package backend

import (
	"context"
	"encoding/json"
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
	// Try cursor-agent first (newer CLI name)
	if _, err := exec.LookPath("cursor-agent"); err == nil {
		return true
	}
	// Fall back to cursor
	_, err := exec.LookPath("cursor")
	return err == nil
}

// cursorJSONResponse represents the JSON output from cursor CLI
type cursorJSONResponse struct {
	SessionID string `json:"session_id"`
	ChatID    string `json:"chat_id"`
	Result    string `json:"result"`
	Response  string `json:"response"`
}

func (c *Cursor) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	// Determine which CLI binary to use
	cliCmd := "cursor"
	if _, err := exec.LookPath("cursor-agent"); err == nil {
		cliCmd = "cursor-agent"
	}

	args := []string{"-p", prompt, "--output-format", "json"}

	if opts.SessionID != "" {
		args = append(args, "--resume", opts.SessionID)
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, cliCmd, args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("cursor: %w\n%s", err, string(out))
	}

	// Try to parse JSON response
	var resp cursorJSONResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		// Fallback: treat as plain text
		return Result{Response: strings.TrimSpace(string(out))}, nil
	}

	// Get session ID from response (try multiple field names)
	sessionID := resp.SessionID
	if sessionID == "" {
		sessionID = resp.ChatID
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
