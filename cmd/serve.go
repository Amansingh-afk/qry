package cmd

import (
	"github.com/amansingh-afk/qry/internal/server"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ServerStarting(port, workDir)
		if err := server.Start(port, workDir); err != nil {
			ui.Error(err.Error())
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 7133, "port")
}
