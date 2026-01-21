package cmd

import (
	"os"
	"strings"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/session"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Backend  string            `yaml:"backend"`
	Model    string            `yaml:"model,omitempty"`
	Dialect  string            `yaml:"dialect,omitempty"`
	Timeout  string            `yaml:"timeout,omitempty"`
	Defaults map[string]string `yaml:"defaults,omitempty"`
	Session  SessionConfig     `yaml:"session,omitempty"`
	Prompt   string            `yaml:"prompt,omitempty"`
}

type SessionConfig struct {
	TTL string `yaml:"ttl,omitempty"`
}

var defaultModels = map[string]string{
	"claude": "haiku",
	"gemini": "gemini-2.0-flash",
	"codex":  "gpt-4o-mini",
	"cursor": "gpt-4o-mini",
}

const defaultPrompt = `You are a SQL expert. Based on the codebase context (schemas, migrations, models), generate ONLY the SQL query.

Rules:
- Output ONLY the SQL, no explanation
- Use actual table/column names from the codebase
- Use {{dialect}} syntax

Request: {{query}}`

var forceInit bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize QRY in current directory",
	Long: `Initialize QRY configuration in the current directory.
Creates .qry.yaml config file and .qry/ directory for session storage.`,
	Example: `  qry init
  qry init --force  # Reset session and re-index codebase`,
	Run: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&forceInit, "force", false, "force re-initialization (clears session)")
}

func runInit(cmd *cobra.Command, args []string) {
	workDir, _ := os.Getwd()

	// Handle --force: clear existing session
	if forceInit {
		if err := session.Delete(workDir); err != nil {
			ui.Error("Failed to clear session: %s", err)
		} else {
			ui.Success("Session cleared")
		}
	}

	ui.Info("Checking backends...")

	available := backend.Available()

	if len(available) == 0 {
		ui.Error("No backends found")
		ui.Print("")
		ui.Print("  Install one of:")
		for _, name := range backend.List() {
			b, _ := backend.Get(name)
			ui.Print("    %s  →  %s", name, b.InstallCmd())
		}
		ui.Print("")
		os.Exit(1)
	}

	ui.Print("")
	for _, b := range available {
		ui.Success("%s", b.Name())
	}

	selected := available[0].Name()

	// Create .qry/ directory for session storage
	qryDir := session.DirPath(workDir)
	if err := os.MkdirAll(qryDir, 0755); err != nil {
		ui.Error("Failed to create .qry directory: %s", err)
		os.Exit(1)
	}

	// Add .qry/ to .gitignore if not already present
	addToGitignore(workDir)

	configPath := ".qry.yaml"
	if _, err := os.Stat(configPath); err == nil {
		if !forceInit {
			ui.Print("")
			ui.Info(".qry.yaml already exists (use --force to reset session)")
			return
		}
		// With --force, we cleared session but keep existing config
		ui.Print("")
		ui.Info(".qry.yaml exists, session reset")
		ui.Print("")
		ui.Print("  Codebase will be re-indexed on next query")
		ui.Print("")
		return
	}

	config := Config{
		Backend:  selected,
		Dialect:  "postgresql",
		Timeout:  "2m",
		Defaults: defaultModels,
		Session: SessionConfig{
			TTL: "7d",
		},
		Prompt: defaultPrompt,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		ui.Error("Failed to create config: %s", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		ui.Error("Failed to write config: %s", err)
		os.Exit(1)
	}

	ui.Print("")
	ui.Success("Created .qry.yaml")
	ui.Print("")
	ui.Print("  Backend: %s", selected)
	ui.Print("  Model:   %s (from defaults)", defaultModels[selected])
	ui.Print("  Session: %s TTL", config.Session.TTL)
	ui.Print("")
	ui.Print("  Try: qry q \"get active users\"")
	ui.Print("")
}

// addToGitignore adds .qry/ to .gitignore if not already present
func addToGitignore(workDir string) {
	gitignorePath := workDir + "/.gitignore"
	entry := ".qry/"

	// Read existing .gitignore
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return // Can't read, skip
	}

	// Check if already present
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			return // Already present
		}
	}

	// Append .qry/ to .gitignore
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// Add newline if file doesn't end with one
	if len(content) > 0 && content[len(content)-1] != '\n' {
		_, _ = f.WriteString("\n")
	}

	_, _ = f.WriteString(entry + "\n")
}
