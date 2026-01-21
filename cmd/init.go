package cmd

import (
	"os"

	"github.com/amansingh-afk/qry/internal/backend"
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
}

var defaultModels = map[string]string{
	"claude": "haiku",
	"gemini": "gemini-2.0-flash",
	"codex":  "gpt-4o-mini",
	"cursor": "gpt-4o-mini",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize QRY in current directory",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
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

	configPath := ".qry.yaml"
	if _, err := os.Stat(configPath); err == nil {
		ui.Print("")
		ui.Info(".qry.yaml already exists")
		return
	}

	config := Config{
		Backend:  selected,
		Dialect:  "postgresql",
		Timeout:  "2m",
		Defaults: defaultModels,
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
	ui.Print("")
	ui.Print("  Try: qry q \"get active users\"")
	ui.Print("")
}
