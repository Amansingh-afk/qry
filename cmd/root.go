package cmd

import (
	"fmt"
	"os"

	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	backendArg string
	jsonOut    bool
	workDir    string
)

var rootCmd = &cobra.Command{
	Use:   "qry [query]",
	Short: "Ask. Get SQL.",
	Long:  ui.Banner(),
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		runQuery(args[0])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&backendArg, "backend", "b", "", "LLM backend (cursor, claude, gemini, codex)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	workDir, _ = os.Getwd()

	viper.AddConfigPath(workDir)
	viper.SetConfigName(".qry")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("QRY")
	viper.AutomaticEnv()

	viper.SetDefault("backend", "claude")
	viper.SetDefault("dialect", "postgres")

	_ = viper.ReadInConfig()
}

func getBackend() string {
	if backendArg != "" {
		return backendArg
	}
	return viper.GetString("backend")
}

func getWorkDir() string {
	return workDir
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.Version())
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize QRY in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := ".qry.yaml"

		if _, err := os.Stat(configPath); err == nil {
			ui.Info(".qry.yaml already exists")
			return
		}

		content := `# QRY Configuration
backend: claude
dialect: postgres
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			ui.Error("Failed to create config: %s", err)
			os.Exit(1)
		}

		ui.Success("Created .qry.yaml")
	},
}
