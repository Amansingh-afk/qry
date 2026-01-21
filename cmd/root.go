package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/session"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	backendFlag string
	modelFlag   string
	dialectFlag string
	timeoutFlag time.Duration
	jsonFlag    bool
	dryRunFlag  bool
	workDir     string
)

var rootCmd = &cobra.Command{
	Use:   "qry",
	Short: "Ask. Get SQL.",
	Long:  ui.Banner(),
	Run:   runChat, // Default: interactive chat
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)

	rootCmd.PersistentFlags().StringVarP(&backendFlag, "backend", "b", "", "backend (claude, gemini, codex, cursor)")
	rootCmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "", "model to use")
	rootCmd.PersistentFlags().StringVarP(&dialectFlag, "dialect", "d", "", "SQL dialect (postgresql, mysql, sqlite)")
	rootCmd.PersistentFlags().DurationVarP(&timeoutFlag, "timeout", "t", 0, "timeout")

	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
}

func loadConfig() {
	workDir, _ = os.Getwd()

	viper.SetConfigName(".qry")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir)

	// Hardcoded fallbacks (only used if no config)
	viper.SetDefault("backend", "claude")
	viper.SetDefault("dialect", "postgresql")
	viper.SetDefault("timeout", "2m")
	viper.SetDefault("session.ttl", "7d")
	viper.SetDefault("defaults.claude", "haiku")
	viper.SetDefault("defaults.gemini", "gemini-2.0-flash")
	viper.SetDefault("defaults.codex", "gpt-4o-mini")
	viper.SetDefault("defaults.cursor", "gpt-4o-mini")

	_ = viper.ReadInConfig()
}

func getBackend() (backend.Backend, error) {
	name := backendFlag
	if name == "" {
		name = viper.GetString("backend")
	}

	b, err := backend.Get(name)
	if err != nil {
		return nil, err
	}

	if !b.Available() {
		return nil, fmt.Errorf("%s not installed\n\n  Install: %s", b.Name(), b.InstallCmd())
	}

	return b, nil
}

func getModel(backendName string) string {
	// 1. CLI flag takes priority
	if modelFlag != "" {
		return modelFlag
	}

	// 2. Config file "model" field
	if model := viper.GetString("model"); model != "" {
		return model
	}

	// 3. Backend-specific default from config
	if model := viper.GetString("defaults." + backendName); model != "" {
		return model
	}

	// 4. Hardcoded fallback (shouldn't reach here if config exists)
	return ""
}

func getDialect() string {
	if dialectFlag != "" {
		return dialectFlag
	}
	return viper.GetString("dialect")
}

func getTimeout() time.Duration {
	if timeoutFlag > 0 {
		return timeoutFlag
	}
	if t := viper.GetDuration("timeout"); t > 0 {
		return t
	}
	return 2 * time.Minute
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.Version())
	},
}

// getSessionTTL parses the session TTL from config
// Supports formats like "7d", "24h", "168h"
func getSessionTTL() time.Duration {
	ttlStr := viper.GetString("session.ttl")
	if ttlStr == "" {
		return 7 * 24 * time.Hour // Default 7 days
	}

	// Handle "Xd" format (days)
	if len(ttlStr) > 1 && ttlStr[len(ttlStr)-1] == 'd' {
		var days int
		if _, err := fmt.Sscanf(ttlStr, "%dd", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}

	// Try standard duration format
	if d, err := time.ParseDuration(ttlStr); err == nil {
		return d
	}

	return 7 * 24 * time.Hour // Fallback
}

// getSession returns the current session ID if valid, empty string otherwise
func getSession(backendName string) string {
	s, err := session.GetOrCreate(workDir, backendName, getSessionTTL())
	if err != nil || s == nil {
		return ""
	}
	return s.SessionID
}

// saveSession persists the session ID for future use
func saveSession(backendName, sessionID string) {
	if sessionID == "" {
		return
	}
	_ = session.Update(workDir, backendName, sessionID)
}
