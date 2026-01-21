package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/session"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Backend   string            `yaml:"backend"`
	Model     string            `yaml:"model,omitempty"`
	Dialect   string            `yaml:"dialect,omitempty"`
	DBVersion string            `yaml:"db_version,omitempty"`
	Timeout   string            `yaml:"timeout,omitempty"`
	Defaults  map[string]string `yaml:"defaults,omitempty"`
	Session   SessionConfig     `yaml:"session,omitempty"`
	Prompt    string            `yaml:"prompt,omitempty"`
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
- Use {{dialect}}{{version}} syntax

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

	ui.Header("Setup")

	// Step 1: Detect repository
	ui.Step("Detecting repository...")
	ui.Pause()
	repoName := detectRepoName(workDir)
	ui.StepDone(repoName)

	// Handle --force: clear existing session
	if forceInit {
		ui.Pause()
		ui.Step("Clearing session...")
		ui.Pause()
		if err := session.Delete(workDir); err != nil {
			ui.StepWarn("No existing session")
		} else {
			ui.StepDone("Session cleared")
		}
	}

	// Step 2: Check backends
	ui.Pause()
	ui.Step("Checking backends...")
	ui.Pause()

	available := backend.Available()

	if len(available) == 0 {
		ui.StepWarn("No backends found")
		ui.Print("")
		ui.Print("  Install one of:")
		for _, name := range backend.List() {
			b, _ := backend.Get(name)
			ui.Print("    %s  →  %s", name, b.InstallCmd())
		}
		ui.Print("")
		os.Exit(1)
	}

	for _, b := range available {
		ui.StepDone(b.Name())
	}

	selected := available[0].Name()

	// Step 3: Create config
	ui.Pause()
	configPath := ".qry.yaml"
	if _, err := os.Stat(configPath); err == nil {
		if !forceInit {
			ui.Step("Config exists")
			ui.StepItem(".qry.yaml already configured")
			ui.Done("Ready")
			ui.Hint("Use --force to reset session")
			ui.Hint("Run: qry")
			ui.Print("")
			return
		}
		// With --force, we cleared session but keep existing config
		ui.Step("Config exists")
		ui.StepItem("Keeping existing .qry.yaml")
		ui.Done("Session reset")
		ui.Hint("Codebase will be re-indexed on next query")
		ui.Hint("Run: qry")
		ui.Print("")
		return
	}

	ui.Step("Creating config...")
	ui.Pause()

	// Create .qry/ directory for session storage
	qryDir := session.DirPath(workDir)
	if err := os.MkdirAll(qryDir, 0755); err != nil {
		ui.Error("Failed to create .qry directory: %s", err)
		os.Exit(1)
	}

	config := Config{
		Backend:   selected,
		Dialect:   "postgresql",
		DBVersion: "",
		Timeout:   "2m",
		Defaults:  defaultModels,
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

	ui.StepItem("Backend:  %s", selected)
	ui.StepItem("Model:    %s", defaultModels[selected])
	ui.StepItem("Dialect:  %s", config.Dialect)
	ui.StepItem("Session:  %s TTL", config.Session.TTL)

	// Step 4: Update .gitignore
	ui.Pause()
	ui.Step("Updating .gitignore...")
	ui.Pause()
	addToGitignore(workDir)
	ui.StepDone("Added .qry/")

	// Done
	ui.Done("Ready")
	ui.Hint("Run: qry")
	ui.Print("")
}

// detectRepoName returns the repository or directory name
func detectRepoName(workDir string) string {
	// Try to get git remote name
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = workDir
	if out, err := cmd.Output(); err == nil {
		url := strings.TrimSpace(string(out))
		// Extract repo name from URL
		// git@github.com:user/repo.git or https://github.com/user/repo.git
		url = strings.TrimSuffix(url, ".git")
		if idx := strings.LastIndex(url, "/"); idx != -1 {
			return url[idx+1:]
		}
		if idx := strings.LastIndex(url, ":"); idx != -1 {
			return url[idx+1:]
		}
	}

	// Fall back to directory name
	return filepath.Base(workDir)
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
