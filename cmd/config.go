package cmd

import (
	"fmt"
	"os"

	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a config value in .qry.yaml",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		viper.Set(key, value)

		configPath := ".qry.yaml"
		if err := viper.WriteConfigAs(configPath); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
		ui.Success("Set %s = %s", key, value)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString(args[0]))
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%s: %v\n", k, v)
		}
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}
