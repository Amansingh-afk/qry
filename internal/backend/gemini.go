package backend

import (
	"context"
	"encoding/json"
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

// geminiJSONResponse represents the JSON output from gemini CLI
type geminiJSONResponse struct {
	SessionID string `json:"session_id"`
	ChatID    string `json:"chat_id"`
	Result    string `json:"result"`
	Response  string `json:"response"`
}

func (g *Gemini) Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error) {
	args := []string{"-p", prompt, "--output-format", "json"}

	if opts.SessionID != "" {
		args = append(args, "--resume", opts.SessionID)
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, "gemini", args...)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("gemini: %w\n%s", err, string(out))
	}

	// Try to parse JSON response
	var resp geminiJSONResponse
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
